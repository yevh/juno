package osinput

import (
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
)

// StarknetOsInput defines the input required to execute the OS.
// Note: this type will need adapted to the specific Runner being used (eg cairo-lang, snos, etc)
// starknet ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/os_input.py#L29
// snos ref: https://github.com/keep-starknet-strange/snos/blob/main/src/io/input.rs#L18
type StarknetOsInput struct {
	// Todo: Set of nodes that are changed in the contract/class trie from updating state?
	ContractStateCommitmentInfo CommitmentInfo `json:"contract_state_commitment_info"`
	ContractClassCommitmentInfo CommitmentInfo `json:"contract_class_commitment_info"`

	// New classes to be declared in the block
	DeprecatedCompiledClasses map[felt.Felt]core.Cairo0Class   `json:"deprecated_compiled_classes"`
	CompiledClasses           map[felt.Felt]core.CompiledClass `json:"compiled_classes"`

	// Todo ??
	CompiledClassVisitedPcs map[felt.Felt][]felt.Felt `json:"compiled_class_visited_pcs"`

	// Contract data associated with every contract that the batch of transactions require
	Contracts map[felt.Felt]ContractState `json:"contracts"`

	// Mapping from contract-tries class_hash to class-tries compiled_class_hash
	ClassHashToCompiledClassHash map[felt.Felt]felt.Felt `json:"class_hash_to_compiled_class_hash"`

	// Fixed Starknet Config
	GeneralConfig StarknetGeneralConfig `json:"general_config"`

	// New transactions defined in the block
	Transactions []core.Transaction `json:"transactions"`

	// Todo: initial, or final blockhash?
	BlockHash felt.Felt `json:"block_hash"`
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/storage/starknet_storage.py#L29
// Context: When updating a trie (class trie, contract trie), we keep track
// of the nodes, and thier (length,path,value) values, that change.
type CommitmentInfo struct {
	PreviousRoot    felt.Felt                 `json:"previous_root"`
	UpdatedRoot     felt.Felt                 `json:"updated_root"`
	TreeHeight      uint                      `json:"tree_height"`
	CommitmentFacts map[felt.Felt][]felt.Felt `json:"commitment_facts"`
}

type StarknetGeneralConfig struct {
	StarknetOsConfig        StarknetOsConfig   `json:"starknet_os_config"`
	GasPriceBounds          GasPriceBounds     `json:"gas_price_bounds"`
	InvokeTxMaxNSteps       int                `json:"invoke_tx_max_n_steps"`
	ValidateMaxNSteps       int                `json:"validate_max_n_steps"`
	DefaultEthPriceInFri    felt.Felt          `json:"default_eth_price_in_fri"`
	ConstantGasPrice        bool               `json:"constant_gas_price"`
	SequencerAddress        *felt.Felt         `json:"sequencer_address"`
	CairoResourceFeeWeights map[string]float64 `json:"cairo_resource_fee_weights"`
	EnforceL1HandlerFee     bool               `json:"enforce_l1_handler_fee"`
	UseKzgDa                bool               `json:"use_kzg_da"`
}

type StarknetOsConfig struct {
	ChainID                   felt.Felt `json:"chain_id"`
	DeprecatedFeeTokenAddress felt.Felt `json:"deprecated_fee_token_address"`
	FeeTokenAddress           felt.Felt `json:"fee_token_address"`
}

type ContractState struct {
	ContractHash          felt.Felt    `json:"contract_hash"`
	StorageCommitmentTree PatriciaTree `json:"storage_commitment_tree"`
	Nonce                 felt.Felt    `json:"nonce"`
}

type PatriciaTree struct {
	Root   felt.Felt `json:"root"`
	Height uint      `json:"height"`
}

type GasPriceBounds struct {
	MinWeiL1GasPrice     felt.Felt `json:"min_wei_l1_gas_price"`
	MinFriL1GasPrice     felt.Felt `json:"min_fri_l1_gas_price"`
	MaxFriL1GasPrice     felt.Felt `json:"max_fri_l1_gas_price"`
	MinWeiL1DataGasPrice felt.Felt `json:"min_wei_l1_data_gas_price"`
	MinFriL1DataGasPrice felt.Felt `json:"min_fri_l1_data_gas_price"`
	MaxFriL1DataGasPrice felt.Felt `json:"max_fri_l1_data_gas_price"`
}

func (s *StarknetOsInput) MarshalJSON() ([]byte, error) {
	type Alias StarknetOsInput
	aux := &struct {
		DeprecatedCompiledClasses    map[string]core.Cairo0Class   `json:"deprecated_compiled_classes"`
		CompiledClasses              map[string]core.CompiledClass `json:"compiled_classes"`
		CompiledClassVisitedPcs      map[string][]felt.Felt        `json:"compiled_class_visited_pcs"`
		Contracts                    map[string]ContractState      `json:"contracts"`
		ClassHashToCompiledClassHash map[string]felt.Felt          `json:"class_hash_to_compiled_class_hash"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	var ok bool
	aux.DeprecatedCompiledClasses, ok = convertMapKeysToString(s.DeprecatedCompiledClasses).(map[string]core.Cairo0Class)
	if !ok {
		return nil, errors.New("type assertion error in custom marshalling of StarknetOsInput")
	}
	aux.CompiledClasses, ok = convertMapKeysToString(s.CompiledClasses).(map[string]core.CompiledClass)
	if !ok {
		return nil, errors.New("type assertion error in custom marshalling of StarknetOsInput")
	}
	aux.CompiledClassVisitedPcs, ok = convertMapKeysToString(s.CompiledClassVisitedPcs).(map[string][]felt.Felt)
	if !ok {
		return nil, errors.New("type assertion error in custom marshalling of StarknetOsInput")
	}

	aux.Contracts, ok = convertMapKeysToString(s.Contracts).(map[string]ContractState)
	if !ok {
		return nil, errors.New("type assertion error in custom marshalling of StarknetOsInput")
	}
	aux.ClassHashToCompiledClassHash = convertMapKeysToString(s.ClassHashToCompiledClassHash).(map[string]felt.Felt)
	if !ok {
		return nil, errors.New("type assertion error in custom marshalling of StarknetOsInput")
	}
	return json.Marshal(aux)
}

func (ci CommitmentInfo) MarshalJSON() ([]byte, error) {

	newCommitmentFacts := make(map[string][]felt.Felt)
	for key, value := range ci.CommitmentFacts {
		newCommitmentFacts[key.String()] = value
	}

	type Alias CommitmentInfo
	return json.Marshal(&struct {
		PreviousRoot    string                 `json:"previous_root"`
		UpdatedRoot     string                 `json:"updated_root"`
		CommitmentFacts map[string][]felt.Felt `json:"commitment_facts"`
		*Alias
	}{
		PreviousRoot:    ci.PreviousRoot.String(),
		UpdatedRoot:     ci.UpdatedRoot.String(),
		CommitmentFacts: newCommitmentFacts,
		Alias:           (*Alias)(&ci),
	})
}

func (c ContractState) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ContractHash          string       `json:"contract_hash"`
		StorageCommitmentTree PatriciaTree `json:"storage_commitment_tree"`
		Nonce                 string       `json:"nonce"`
	}{
		ContractHash:          c.ContractHash.String(),
		StorageCommitmentTree: c.StorageCommitmentTree,
		Nonce:                 c.Nonce.String(),
	})
}

func (p PatriciaTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Root   string `json:"root"`
		Height uint   `json:"height"`
	}{
		Root:   p.Root.String(),
		Height: p.Height,
	})
}
func convertMapKeysToString(originalMap interface{}) interface{} {
	switch m := originalMap.(type) {
	case map[felt.Felt]core.Cairo0Class:
		newMap := make(map[string]core.Cairo0Class)
		for k, v := range m {
			newMap[k.String()] = v // Assuming Felt has a String() method
		}
		return newMap
	case map[felt.Felt]core.CompiledClass:
		newMap := make(map[string]core.CompiledClass)
		for k, v := range m {
			newMap[k.String()] = v
		}
		return newMap
	case map[felt.Felt][]felt.Felt:
		newMap := make(map[string][]felt.Felt)
		for k, v := range m {
			newMap[k.String()] = v
		}
		return newMap
	case map[felt.Felt]ContractState:
		newMap := make(map[string]ContractState)
		for k, v := range m {
			newMap[k.String()] = v
		}
		return newMap
	case map[felt.Felt]felt.Felt:
		newMap := make(map[string]felt.Felt)
		for k, v := range m {
			newMap[k.String()] = v
		}
		return newMap
	default:
		return nil
	}
}
