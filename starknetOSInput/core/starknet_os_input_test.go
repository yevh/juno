package osinput

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/juno/blockchain"
	"github.com/NethermindEth/juno/clients/feeder"
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/db/pebble"
	"github.com/NethermindEth/juno/mocks"
	adaptfeeder "github.com/NethermindEth/juno/starknetdata/feeder"
	"github.com/NethermindEth/juno/utils"
	"github.com/NethermindEth/juno/vm"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Note this works for Sepolia and not mainnet (mainnet had bugs in the early days)
func getNewClasses(client adaptfeeder.Feeder, blockNumber uint64) (map[felt.Felt]core.Class, error) {
	su, err := client.StateUpdate(context.Background(), blockNumber)
	if err != nil {
		return nil, err
	}
	var classesToFetch []*felt.Felt
	classesToFetch = append(classesToFetch, su.StateDiff.DeclaredV0Classes...)
	for classHash := range su.StateDiff.DeclaredV1Classes {
		classesToFetch = append(classesToFetch, &classHash)
	}
	classes := make(map[felt.Felt]core.Class)
	for _, classHash := range classesToFetch {
		fmt.Println("classHash", classHash)
		class, err := client.Class(context.Background(), classHash)
		if err != nil {
			return nil, err
		}
		classes[*classHash] = class
	}
	return classes, nil
}

func TestGenerateStarknetOSInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	network := utils.Sepolia

	// exampleConfig := LoadExampleStarknetOSConfig()
	// Todo: get test data for first mainet block (may need to feed this into run_os.py somehow)
	expectedOSInptsEmpty := StarknetOsInput{}

	// Two new (deprecated) classes, three new contracts, all three contracts have non-zero storage, but only one has non-zero nonce.
	t.Run("empty state to sepolia block 0", func(t *testing.T) {

		testDB := pebble.NewMemTest(t)
		chain := blockchain.New(testDB, &network)
		client := feeder.NewTestClient(t, &network)
		gw := adaptfeeder.New(client)
		mockVM := mocks.NewMockVM(mockCtrl)

		// old empty state
		txn, err := testDB.NewTransaction(false)
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, txn.Discard())
		})
		oldState := core.NewState(txn)

		// get and apply block0 to get new state
		block0, err := gw.BlockByNumber(context.Background(), uint64(0))
		require.NoError(t, err)
		su0, err := gw.StateUpdate(context.Background(), uint64(0))
		require.NoError(t, err)
		newClasses, err := getNewClasses(*gw, 0)
		require.NoError(t, err)

		require.NoError(t, chain.Store(block0, &core.BlockCommitments{}, su0, newClasses))

		newState, closer, err := chain.StateAtBlockNumber(0)
		require.NoError(t, err)

		var classes []core.Class
		for _, class := range newClasses {
			classes = append(classes, class)
		}

		vmParas := VMParameters{
			Txns:            block0.Transactions,
			DeclaredClasses: classes,
			PaidFeesOnL1:    nil,
			State:           oldState,
			Network:         &network,
			SkipChargeFee:   false,
			SkipValidate:    false,
			ErrOnRevert:     false,
			UseBlobData:     false,
			BlockInfo:       &vm.BlockInfo{Header: block0.Header},
		}

		mockVM.EXPECT().Execute(vmParas.Txns, vmParas.DeclaredClasses, vmParas.PaidFeesOnL1,
			vmParas.BlockInfo, vmParas.State, vmParas.Network, vmParas.SkipChargeFee, vmParas.SkipValidate,
			vmParas.ErrOnRevert, vmParas.UseBlobData).Return(nil, nil, nil, nil)

		osinput, err := GenerateStarknetOSInput(oldState, newState, *block0, mockVM, vmParas)
		require.NoError(t, err)

		qwe, err := json.MarshalIndent(osinput, "", "")
		require.NoError(t, err)
		os.WriteFile("osinput.json", qwe, 0644)

		require.Equal(t, expectedOSInptsEmpty.ContractStateCommitmentInfo, osinput.ContractStateCommitmentInfo)
		require.Equal(t, expectedOSInptsEmpty.ContractClassCommitmentInfo, osinput.ContractClassCommitmentInfo)
		require.Equal(t, expectedOSInptsEmpty.DeprecatedCompiledClasses, osinput.DeprecatedCompiledClasses)
		require.Equal(t, expectedOSInptsEmpty.CompiledClasses, osinput.CompiledClasses)
		require.Equal(t, expectedOSInptsEmpty.CompiledClassVisitedPcs, osinput.CompiledClassVisitedPcs)
		require.Equal(t, expectedOSInptsEmpty.Contracts, osinput.Contracts)
		require.Equal(t, expectedOSInptsEmpty.ClassHashToCompiledClassHash, osinput.ClassHashToCompiledClassHash)
		require.Equal(t, expectedOSInptsEmpty.GeneralConfig, osinput.GeneralConfig)
		require.Equal(t, expectedOSInptsEmpty.Transactions, osinput.Transactions)
		require.Equal(t, expectedOSInptsEmpty.BlockHash.String(), osinput.BlockHash.String())
		require.NoError(t, closer())
	})

	// t.Run(" mainnet block 0 to 1", func(t *testing.T) {
	// 	testDB := pebble.NewMemTest(t)
	// 	chain := blockchain.New(testDB, &network)
	// 	client := feeder.NewTestClient(t, &network)
	// 	gw := adaptfeeder.New(client)
	// 	mockVM := mocks.NewMockVM(mockCtrl)

	// 	// get and apply block0 to get new state
	// 	block0, err := gw.BlockByNumber(context.Background(), uint64(0))
	// 	require.NoError(t, err)
	// 	su0, err := gw.StateUpdate(context.Background(), uint64(0))
	// 	require.NoError(t, err)
	// 	classHash0 := utils.HexToFelt(t, "0x10455c752b86932ce552f2b0fe81a880746649b9aee7e0d842bf3f52378f9f8")
	// 	class0, err := gw.Class(context.Background(), utils.HexToFelt(t, "0x10455c752b86932ce552f2b0fe81a880746649b9aee7e0d842bf3f52378f9f8"))
	// 	require.NoError(t, err)
	// 	require.NoError(t, chain.Store(block0, &core.BlockCommitments{}, su0, map[felt.Felt]core.Class{*classHash0: class0}))

	// 	// get and apply block1 to get new state
	// 	block1, err := gw.BlockByNumber(context.Background(), uint64(1))
	// 	require.NoError(t, err)
	// 	su1, err := gw.StateUpdate(context.Background(), uint64(1))
	// 	require.NoError(t, err)
	// 	require.NoError(t, err)
	// 	require.NoError(t, chain.Store(block1, &core.BlockCommitments{}, su1, nil))

	// 	oldState, oldCloser, err := chain.StateAtBlockNumber(0)
	// 	require.NoError(t, err)
	// 	newState, newCloser, err := chain.StateAtBlockNumber(1)
	// 	require.NoError(t, err)

	// 	vmParas := VMParameters{
	// 		Txns:            block0.Transactions,
	// 		DeclaredClasses: []core.Class{class0},
	// 		PaidFeesOnL1:    nil,
	// 		State:           oldState,
	// 		Network:         &network,
	// 		SkipChargeFee:   false,
	// 		SkipValidate:    false,
	// 		ErrOnRevert:     false,
	// 		UseBlobData:     false,
	// 		BlockInfo:       &vm.BlockInfo{Header: block0.Header},
	// 	}

	// 	mockVM.EXPECT().Execute(vmParas.Txns, vmParas.DeclaredClasses, vmParas.PaidFeesOnL1,
	// 		vmParas.BlockInfo, vmParas.State, vmParas.Network, vmParas.SkipChargeFee, vmParas.SkipValidate,
	// 		vmParas.ErrOnRevert, vmParas.UseBlobData).Return(nil, nil, nil, nil)

	// 	osinput, err := GenerateStarknetOSInput(oldState, newState, *block0, mockVM, vmParas)
	// 	require.NoError(t, err)

	// 	qwe, err := json.MarshalIndent(osinput, "", "")
	// 	require.NoError(t, err)
	// 	os.WriteFile("osinput.json", qwe, 0644)

	// 	require.Equal(t, expectedOSInptsEmpty.ContractStateCommitmentInfo, osinput.ContractStateCommitmentInfo)
	// 	require.Equal(t, expectedOSInptsEmpty.ContractClassCommitmentInfo, osinput.ContractClassCommitmentInfo)
	// 	require.Equal(t, expectedOSInptsEmpty.DeprecatedCompiledClasses, osinput.DeprecatedCompiledClasses)
	// 	require.Equal(t, expectedOSInptsEmpty.CompiledClasses, osinput.CompiledClasses)
	// 	require.Equal(t, expectedOSInptsEmpty.CompiledClassVisitedPcs, osinput.CompiledClassVisitedPcs)
	// 	require.Equal(t, expectedOSInptsEmpty.Contracts, osinput.Contracts)
	// 	require.Equal(t, expectedOSInptsEmpty.ClassHashToCompiledClassHash, osinput.ClassHashToCompiledClassHash)
	// 	require.Equal(t, expectedOSInptsEmpty.GeneralConfig, osinput.GeneralConfig)
	// 	require.Equal(t, expectedOSInptsEmpty.Transactions, osinput.Transactions)
	// 	require.Equal(t, expectedOSInptsEmpty.BlockHash.String(), osinput.BlockHash.String())
	// 	require.NoError(t, oldCloser())
	// 	require.NoError(t, newCloser())
	// })

	// expectedOSInptsEmpty := StarknetOsInput{
	// 	// "0x0" has no state (no nonce, no classhash)
	// 	// "0x1" has storage, and class hash "0x0", and nonce "0x0". The "old_root" for block 0 is "0x0" suggesting it has no state?
	// 	ContractStateCommitmentInfo: CommitmentInfo{
	// 		PreviousRoot:    *new(felt.Felt).SetUint64(0),
	// 		UpdatedRoot:     *new(felt.Felt).SetUint64(0),
	// 		TreeHeight:      251,
	// 		CommitmentFacts: map[felt.Felt][]felt.Felt{},
	// 	},
	// 	// "0x0" has no storage / classhash, and therefore no corresponding class
	// 	// "0x1" has storage, and class hash "0x0", but no class exists with class hash "0x0"
	// 	ContractClassCommitmentInfo: CommitmentInfo{
	// 		PreviousRoot:    *new(felt.Felt).SetUint64(0),
	// 		UpdatedRoot:     *new(felt.Felt).SetUint64(0),
	// 		TreeHeight:      251,
	// 		CommitmentFacts: map[felt.Felt][]felt.Felt{},
	// 	},
	// 	DeprecatedCompiledClasses: nil,
	// 	CompiledClasses:           nil,
	// 	CompiledClassVisitedPcs:   nil,
	// 	// run_os.py returns zeros, even for non-deployed contracts
	// 	Contracts: map[felt.Felt]ContractState{
	// 		*new(felt.Felt).SetUint64(0): {
	// 			ContractHash: *new(felt.Felt).SetUint64(0),
	// 			StorageCommitmentTree: PatriciaTree{
	// 				Root:   *new(felt.Felt).SetUint64(0),
	// 				Height: 251,
	// 			},
	// 			Nonce: *new(felt.Felt).SetUint64(0),
	// 		},
	// 		*new(felt.Felt).SetUint64(1): {
	// 			ContractHash: *new(felt.Felt).SetUint64(0),
	// 			StorageCommitmentTree: PatriciaTree{
	// 				Root:   *new(felt.Felt).SetUint64(0),
	// 				Height: 251,
	// 			},
	// 			Nonce: *new(felt.Felt).SetUint64(0),
	// 		},
	// 	},
	// 	ClassHashToCompiledClassHash: map[felt.Felt]felt.Felt{},
	// 	GeneralConfig:                exampleConfig,
	// 	Transactions:                 nil,
	// 	BlockHash:                    *utils.HexToFelt(t, "2535437458273622887584459710067137978693525181086955024571735059458497227738"),
	// }

	// Declare and deploy dummy_token.json
	// expectedOSInptsDummyToken := StarknetOsInput{
	// 	ContractStateCommitmentInfo: CommitmentInfo{
	// 		PreviousRoot: *new(felt.Felt).SetUint64(0),
	// 		UpdatedRoot:  *utils.HexToFelt(t, "587553090332532752877781043098123845316873655073341242723423693336333123978"),
	// 		TreeHeight:   251,
	// 		CommitmentFacts: map[felt.Felt][]felt.Felt{
	// 			*utils.HexToFelt(t, "0x14c8b135d7babe1581dd8f67002c5482be3e7a52bc0235f875ee6dfc582018a"): {
	// 				*utils.HexToFelt(t, "0xfb"),
	// 				*utils.HexToFelt(t, "0x3400a86fdc294a70fac1cf84f81a2127419359096b846be9814786d4fc056b8"),
	// 				*utils.HexToFelt(t, "0x7a7555584f4d26fd18050fb1ab401491b77f8b664c14e2fe21cbb6d3df0dfe5"),
	// 			},
	// 		},
	// 	},

	// 	ContractClassCommitmentInfo: CommitmentInfo{
	// 		PreviousRoot:    *new(felt.Felt).SetUint64(0),
	// 		UpdatedRoot:     *new(felt.Felt).SetUint64(0),
	// 		TreeHeight:      251,
	// 		CommitmentFacts: map[felt.Felt][]felt.Felt{},
	// 	},
	// 	DeprecatedCompiledClasses: map[felt.Felt]core.Cairo0Class{
	// 		*utils.HexToFelt(t, "0x7cea4d7710723fa9e33472b6ceb71587a0ce4997ef486638dd0156bdb6c2daa"): {}, // Todo
	// 	},
	// 	CompiledClasses:         nil,
	// 	CompiledClassVisitedPcs: nil,
	// 	// run_os.py returns zeros, even for non-deployed contracts
	// 	Contracts: map[felt.Felt]ContractState{
	// 		*new(felt.Felt).SetUint64(0): {
	// 			ContractHash: *new(felt.Felt).SetUint64(0),
	// 			StorageCommitmentTree: PatriciaTree{
	// 				Root:   *new(felt.Felt).SetUint64(0),
	// 				Height: 251,
	// 			},
	// 			Nonce: *new(felt.Felt).SetUint64(0),
	// 		},
	// 		*new(felt.Felt).SetUint64(1): {
	// 			ContractHash: *new(felt.Felt).SetUint64(0),
	// 			StorageCommitmentTree: PatriciaTree{
	// 				Root:   *new(felt.Felt).SetUint64(0),
	// 				Height: 251,
	// 			},
	// 			Nonce: *new(felt.Felt).SetUint64(0),
	// 		},
	// 		*utils.HexToFelt(t, "1470089414715992704702781317133162679047468004062084455026636858461958198968"): {
	// 			ContractHash: *new(felt.Felt).SetUint64(0),
	// 			StorageCommitmentTree: PatriciaTree{
	// 				Root:   *new(felt.Felt).SetUint64(0),
	// 				Height: 251,
	// 			},
	// 			Nonce: *new(felt.Felt).SetUint64(0),
	// 		},
	// 	},
	// 	ClassHashToCompiledClassHash: map[felt.Felt]felt.Felt{},
	// 	GeneralConfig:                exampleConfig,
	// 	Transactions: []core.Transaction{
	// 		&core.DeployAccountTransaction{
	// 			MaxFee: utils.HexToFelt(t, "0x10000000000000000000000000"),
	// 			Nonce:  &felt.Zero,
	// 			DeployTransaction: core.DeployTransaction{
	// 				TransactionHash: utils.HexToFelt(t, "0xcd76933991f9453baa217e0c0f9090b0a48c6922c74ede5d5e2faa36e4e68"),
	// 				Version:         new(core.TransactionVersion).SetUint64(1),
	// 				// SenderAddress:       utils.HexToFelt(t, "0x3400a86fdc294a70fac1cf84f81a2127419359096b846be9814786d4fc056b8"), // Todo:Not used??
	// 				// Type: "DEPLOY_ACCOUNT", // Todo: Switch from core.Transaction to rpc.Transaction?
	// 				ContractAddressSalt: &felt.Zero,
	// 				ClassHash:           utils.HexToFelt(t, "0x7cea4d7710723fa9e33472b6ceb71587a0ce4997ef486638dd0156bdb6c2daa"),
	// 				ConstructorCallData: []*felt.Felt{},
	// 			},
	// 		},
	// 	},
	// 	BlockHash: *utils.HexToFelt(t, "2535437458273622887584459710067137978693525181086955024571735059458497227738"),
	// }

	// // Todo
	// t.Run("0x0 and 0x2 contracts + declare + deploy dummy_token - todo", func(t *testing.T) {
	// 	su := &core.StateUpdate{
	// 		OldRoot: &felt.Zero,
	// 		NewRoot: &felt.Zero,
	// 		// BlockHash: &felt.Zero, // Not used
	// 		StateDiff: &core.StateDiff{
	// 			StorageDiffs:      nil,
	// 			Nonces:            nil,
	// 			DeployedContracts: nil,
	// 			DeclaredV0Classes: nil,
	// 			DeclaredV1Classes: nil,
	// 		},
	// 	}
	// 	err := newState.Update(0, su, nil)
	// 	require.NoError(t, err)

	// 	vmParas := VMParameters{
	// 		Txns:            nil,
	// 		DeclaredClasses: nil,
	// 		PaidFeesOnL1:    nil,
	// 		State:           oldState,
	// 		Network:         &network,
	// 		SkipChargeFee:   false,
	// 		SkipValidate:    false,
	// 		ErrOnRevert:     false,
	// 		UseBlobData:     false,
	// 	}
	// 	block := core.Block{
	// 		Header: &core.Header{
	// 			Hash: utils.HexToFelt(t, "0x59b01ba262c999f2617412ffbba780f80b0103d928cbce1aecbaa50de90abda"),
	// 		},
	// 	}
	// 	mockVM.EXPECT().Execute(vmParas.Txns, vmParas.DeclaredClasses, vmParas.PaidFeesOnL1,
	// 		vmParas.BlockInfo, oldState, vmParas.Network, vmParas.SkipChargeFee, vmParas.SkipValidate, vmParas.ErrOnRevert, vmParas.UseBlobData).Return(nil, nil, nil, nil)

	// 	osinput, err := GenerateStarknetOSInput(&block, oldState, newState, mockVM, vmParas)
	// 	require.NoError(t, err)

	// 	require.Equal(t, expectedOSInptsDummyToken.ContractStateCommitmentInfo, osinput.ContractStateCommitmentInfo)
	// 	require.Equal(t, expectedOSInptsDummyToken.ContractClassCommitmentInfo, osinput.ContractClassCommitmentInfo)
	// 	require.Equal(t, expectedOSInptsDummyToken.DeprecatedCompiledClasses, osinput.DeprecatedCompiledClasses)
	// 	require.Equal(t, expectedOSInptsDummyToken.CompiledClasses, osinput.CompiledClasses)
	// 	require.Equal(t, expectedOSInptsDummyToken.CompiledClassVisitedPcs, osinput.CompiledClassVisitedPcs)
	// 	require.Equal(t, expectedOSInptsDummyToken.Contracts, osinput.Contracts)
	// 	require.Equal(t, expectedOSInptsDummyToken.ClassHashToCompiledClassHash, osinput.ClassHashToCompiledClassHash)
	// 	require.Equal(t, expectedOSInptsDummyToken.GeneralConfig, osinput.GeneralConfig)
	// 	require.Equal(t, expectedOSInptsDummyToken.Transactions, osinput.Transactions)
	// 	require.Equal(t, expectedOSInptsDummyToken.BlockHash.String(), osinput.BlockHash.String())
	// })

}

// func loadInitClasses() ([]core.Cairo0Class, []core.Cairo1Class, error) {
// 	loadJSON := func(filePath string, target interface{}) {
// 		data, err := os.ReadFile(filePath)
// 		if err != nil {
// 			log.Fatalf("unable to read json file: %v", err)
// 		}
// 		if err := json.Unmarshal(data, target); err != nil {
// 			panic(fmt.Errorf("unable to unmarshal json compiled class: %v", err))
// 		}
// 	}

// 	testContractClass := new(core.Cairo0Class)
// 	loadJSON("testdata/test_contract.json", testContractClass)

// 	dummyAccountClass := new(core.Cairo0Class)
// 	loadJSON("testdata/dummy_accout.json", dummyAccountClass)

// 	dummyTokenClass := new(core.Cairo0Class)
// 	loadJSON("testdata/dummy_token.json", dummyTokenClass)

// 	return []core.Cairo0Class{*testContractClass, *dummyAccountClass, *dummyTokenClass}, nil, nil
// }

// func stringToFelt(s string) *felt.Felt {
// 	f, err := new(felt.Felt).SetString(s)
// 	if err != nil {
// 		panic(err)

// 	}
// 	return f
// }

// func getInitTxns(depClasses []core.Cairo0Class) ([]core.Transaction, error) {

// 	deployDummyTokenTransaction := core.DeployAccountTransaction{
// 		DeployTransaction: core.DeployTransaction{
// 			TransactionHash:     stringToFelt("22688876470218804543887986415455328819098091743988100398197353790124740200"),
// 			Version:             new(core.TransactionVersion).SetUint64(1),
// 			ContractAddressSalt: new(felt.Felt).SetUint64(0),
// 			ClassHash:           stringToFelt("3531298130119845387864440863187980726515137569165069484670944625223023734186"),
// 			ConstructorCalldata: []byte{},
// 		},
// 		MaxFee:               new(felt.Felt).SetBigInt(big.NewInt(1267650600228229401496703205376)),
// 		TransactionSignature: []*felt.Felt{},
// 		Nonce:                new(felt.Felt).SetUint64(0),
// 	}

// 	deployDummyAccountTransaction := core.DeployAccountTransaction{
// 		DeployTransaction: core.DeployTransaction{
// 			TransactionHash:     new(felt.Felt).SetBigInt(big.NewInt(96374521715508826444566467091393680183010464453336720810014746622481735761)),
// 			Version:             new(core.TransactionVersion).SetUint64(1),
// 			SenderAddress:       new(felt.Felt).SetBigInt(big.NewInt(2618767603815038378512366346550627731109766804643583016834052353912473402832)),
// 			ContractAddressSalt: new(felt.Felt).SetUint64(0),
// 			ClassHash:           big.NewInt(646245114977324210659279014519951538684823368221946044944492064370769527799),
// 			ConstructorCalldata: []byte{},
// 		},
// 		MaxFee:               new(felt.Felt).SetBigInt(big.NewInt(1267650600228229401496703205376)),
// 		TransactionSignature: []*felt.Felt{},
// 		Nonce:                new(felt.Felt).SetUint64(0),
// 	}

// 	fundInvokeTransaction := core.InvokeTransaction{
// 		TransactionHash:      new(felt.Felt).SetBigInt(big.NewInt(2852915394592604060963909836770256627436576776991723431020681987492769528640)),
// 		Version:              new(core.TransactionVersion).SetUint64(1),
// 		MaxFee:               new(felt.Felt).SetBigInt(big.NewInt(1267650600228229401496703205376)),
// 		TransactionSignature: []*felt.Felt{},
// 		Nonce:                new(felt.Felt).SetUint64(1),
// 		SenderAddress:        new(felt.Felt).SetBigInt(big.NewInt(1470089414715992704702781317133162679047468004062084455026636858461958198968)),
// 		EntryPointSelector:   new(felt.Felt).SetBigInt(big.NewInt(617075754465154585683856897856256838130216341506379215893724690153393808813)),
// 		EntryPointType:       core.EntryPointType.EXTERNAL,
// 		CallData: []*felt.Felt{
// 			new(felt.Felt).SetBigInt(big.NewInt(1470089414715992704702781317133162679047468004062084455026636858461958198968)),
// 			new(felt.Felt).SetBigInt(big.NewInt(232670485425082704932579856502088130646006032362877466777181098476241604910)),
// 			new(felt.Felt).SetUint64(3),
// 			new(felt.Felt).SetBigInt(big.NewInt(2618767603815038378512366346550627731109766804643583016834052353912473402832)),
// 			new(felt.Felt).SetBigInt(big.NewInt(1329227995784915872903807060280344576)),
// 			new(felt.Felt).SetUint64(0),
// 		},
// 	}

// 	return []core.Transaction{&deployDummyAccountTransaction, &deployDummyTokenTransaction, &fundInvokeTransaction}, nil
// }
