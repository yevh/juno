package trie

import (
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

// https://github.com/starknet-io/starknet-p2p-specs/blob/main/p2p/proto/snapshot.proto#L6
type ProofNode struct {
	Binary *Binary
	Edge   *Edge
}

// Note: does not work for leaves
func (pn *ProofNode) Hash(hash hashFunc) *felt.Felt {
	switch {
	case pn.Binary != nil:
		return hash(pn.Binary.LeftHash, pn.Binary.RightHash)
	case pn.Edge != nil:
		length := make([]byte, len(pn.Edge.Path.bitset))
		length[len(pn.Edge.Path.bitset)-1] = pn.Edge.Path.len
		pathFelt := pn.Edge.Path.Felt()
		lengthFelt := new(felt.Felt).SetBytes(length)
		return new(felt.Felt).Add(hash(pn.Edge.Child, &pathFelt), lengthFelt)
	default:
		return nil
	}
}

func (pn *ProofNode) PrettyPrint() {
	if pn.Binary != nil {
		fmt.Printf("  Binary:\n")
		fmt.Printf("    LeftHash: %v\n", pn.Binary.LeftHash)
		fmt.Printf("    RightHash: %v\n", pn.Binary.RightHash)
	}
	if pn.Edge != nil {
		fmt.Printf("  Edge:\n")
		fmt.Printf("    Child: %v\n", pn.Edge.Child)
		fmt.Printf("    Path: %v\n", pn.Edge.Path)
		fmt.Printf("    Value: %v\n", pn.Edge.Value)
	}
}

type Binary struct {
	LeftHash  *felt.Felt
	RightHash *felt.Felt
}

type Edge struct {
	Child *felt.Felt
	Path  *Key
	Value *felt.Felt
}

func isEdge(parentKey *Key, sNode storageNode) bool {
	sNodeLen := sNode.key.len
	if parentKey == nil { // Root
		return sNodeLen != 0
	}
	return sNodeLen-parentKey.len > 1
}

// Note: we need to account for the fact that Junos Trie has nodes that are Binary AND Edge,
// whereas the protocol requires nodes that are Binary XOR Edge
func transformNode(tri *Trie, parentKey *Key, sNode storageNode) (*Edge, *Binary, error) {
	isEdgeBool := isEdge(parentKey, sNode)

	var edge *Edge
	if isEdgeBool {
		edgePath := path(sNode.key, parentKey)
		edge = &Edge{
			Path:  &edgePath,
			Child: sNode.node.Value,
		}
	}
	if sNode.key.len == tri.height { // Leaf
		return edge, nil, nil
	}
	lNode, err := tri.GetNodeFromKey(sNode.node.Left)
	if err != nil {
		return nil, nil, err
	}
	rNode, err := tri.GetNodeFromKey(sNode.node.Right)
	if err != nil {
		return nil, nil, err
	}

	rightHash := rNode.Value
	if isEdge(sNode.key, storageNode{node: rNode, key: sNode.node.Right}) {
		edgePath := path(sNode.node.Right, sNode.key)
		rEdge := ProofNode{Edge: &Edge{
			Path:  &edgePath,
			Child: rNode.Value,
		}}
		rightHash = rEdge.Hash(tri.hash)
	}
	leftHash := lNode.Value
	if isEdge(sNode.key, storageNode{node: lNode, key: sNode.node.Left}) {
		edgePath := path(sNode.node.Left, sNode.key)
		lEdge := ProofNode{Edge: &Edge{
			Path:  &edgePath,
			Child: lNode.Value,
		}}
		leftHash = lEdge.Hash(tri.hash)
	}
	binary := &Binary{
		LeftHash:  leftHash,
		RightHash: rightHash,
	}

	return edge, binary, nil
}

// https://github.com/eqlabs/pathfinder/blob/main/crates/merkle-tree/src/tree.rs#L514
func GetProof(leaf *felt.Felt, tri *Trie) ([]ProofNode, error) {
	leafKey := tri.feltToKey(leaf)
	nodesToLeaf, err := tri.nodesFromRoot(&leafKey)
	if err != nil {
		return nil, err
	}
	proofNodes := []ProofNode{}

	var parentKey *Key

	for i := 0; i < len(nodesToLeaf); i++ {
		sNode := nodesToLeaf[i]
		sNodeEdge, sNodeBinary, err := transformNode(tri, parentKey, sNode)
		if err != nil {
			return nil, err
		}
		isLeaf := sNode.key.len == tri.height

		if sNodeEdge != nil && !isLeaf { // Internal Edge
			proofNodes = append(proofNodes, []ProofNode{{Edge: sNodeEdge}, {Binary: sNodeBinary}}...)
		} else if sNodeEdge == nil && !isLeaf { // Internal Binary
			proofNodes = append(proofNodes, []ProofNode{{Binary: sNodeBinary}}...)
		} else if sNodeEdge != nil && isLeaf { // Leaf Edge
			proofNodes = append(proofNodes, []ProofNode{{Edge: sNodeEdge}}...)
		} else if sNodeEdge == nil && sNodeBinary == nil { // sNode is a binary leaf
			break
		}
		parentKey = nodesToLeaf[i].key
	}
	return proofNodes, nil
}

func GetProofs(startKey, endKey *felt.Felt, tri *Trie) ([][]ProofNode, error) {
	oneFelt := new(felt.Felt).SetUint64(1)
	iterKey := startKey
	leafRange := new(felt.Felt).Sub(endKey, startKey).Uint64()
	proofs := make([][]ProofNode, leafRange)
	for i := range leafRange {
		proof, err := GetProof(iterKey, tri)
		if err != nil {
			return nil, err
		}
		proofs[i] = proof
		iterKey.Add(iterKey, oneFelt)
	}
	return proofs, nil
}

// verifyProof checks if `leafPath` leads from `root` to `leafHash` along the `proofNodes`
// https://github.com/eqlabs/pathfinder/blob/main/crates/merkle-tree/src/tree.rs#L2006
func VerifyProof(root *felt.Felt, key *Key, value *felt.Felt, proofs []ProofNode, hash hashFunc) bool {
	if key.Len() != 251 { //nolint:gomnd
		return false
	}

	expectedHash := root
	remainingPath := key

	for _, proofNode := range proofs {
		if !proofNode.Hash(hash).Equal(expectedHash) {
			return false
		}
		switch {
		case proofNode.Binary != nil:
			if remainingPath.Test(remainingPath.Len() - 1) {
				expectedHash = proofNode.Binary.RightHash
			} else {
				expectedHash = proofNode.Binary.LeftHash
			}
			remainingPath.RemoveLastBit()
		case proofNode.Edge != nil:
			if !proofNode.Edge.Path.Equal(remainingPath.SubKey(proofNode.Edge.Path.Len())) {
				return false
			}
			expectedHash = proofNode.Edge.Child
			remainingPath.Truncate(proofNode.Edge.Path.Len())
		}
	}
	return expectedHash.Equal(value)
}

// verifyRangeProof verifies the range proof for the given range of keys.
// ref: https://github.com/ethereum/go-ethereum/blob/v1.14.3/trie/proof.go#L484
func verifyRangeProof(root *felt.Felt, firstKey *Key, firstProof []ProofNode, lastKey *Key, lastProof []ProofNode,
	innerKeys []*Key, innerValues []*felt.Felt, innerProofs [][]ProofNode, hash hashFunc) (bool, error) {
	// Step 0: checks
	if len(innerKeys) != len(innerValues) {
		return false, fmt.Errorf("inconsistent proof data, keys: %d, values: %d", len(innerKeys), len(innerValues))
	}
	// Ensure the received batch is monotonic increasing
	for i := 0; i < len(innerKeys)-1; i++ {
		if innerKeys[i].cmp(innerKeys[i+1]) >= 0 {
			return false, errors.New("range is not monotonically increasing")
		}
	}
	// Ensure the received batch contains no deletions
	for _, value := range innerValues {
		if value.Equal(&felt.Zero) {
			return false, errors.New("range contains deletion")
		}
	}

	// Step 1: Verify the first edge proof (firstKey)
	if !VerifyProof(root, firstKey, value, firstProof, hash) { // Todo: value??
		return false, errors.New("Invalid first edge proof")
	}

	// Step 2: Verify the last edge proof (lastKey)
	if !VerifyProof(root, lastKey, value, lastProof, hash) { // Todo: value??
		return false, errors.New("Invalid last edge proof")
	}

	// Step 3: Verify each key and value in the given range
	for i, key := range innerKeys {

		// Verify the key using its proof
		if !VerifyProof(root, key, innerValues[i], innerProofs[i], hash) {
			return false, errors.New(fmt.Sprintf("Invalid proof for key %x", key))
		}

		// Verify the value associated with the key
		if !verifyValueFromProof(innerProofs[i], innerValues[i], hash) {
			return false, errors.New(fmt.Sprintf("Incorrect value for key %x", key))
		}
	}

	// Step 4: Recompute the root hash from the verified paths
	recomputedRoot := recomputeRootHash(innerKeys, innerValues, innerProofs, hash)

	// Verify that the recomputed root hash matches the provided root hash
	if !recomputedRoot.Eq(root) {
		return false, errors.New("Root hash mismatch")
	}

	return true, errors.New("Proof is valid")
}
