package ipfs

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/moresearch/swechain/testutil/sample"
	ipfssimulation "github.com/moresearch/swechain/x/ipfs/simulation"
	"github.com/moresearch/swechain/x/ipfs/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ipfsGenesis := types.GenesisState{
		Params: types.DefaultParams(), CodingTrajList: []types.CodingTraj{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ipfsGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateCodingTraj          = "op_weight_msg_ipfs"
		defaultWeightMsgCreateCodingTraj int = 100
	)

	var weightMsgCreateCodingTraj int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCodingTraj, &weightMsgCreateCodingTraj, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCodingTraj = defaultWeightMsgCreateCodingTraj
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCodingTraj,
		ipfssimulation.SimulateMsgCreateCodingTraj(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateCodingTraj          = "op_weight_msg_ipfs"
		defaultWeightMsgUpdateCodingTraj int = 100
	)

	var weightMsgUpdateCodingTraj int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCodingTraj, &weightMsgUpdateCodingTraj, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCodingTraj = defaultWeightMsgUpdateCodingTraj
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCodingTraj,
		ipfssimulation.SimulateMsgUpdateCodingTraj(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteCodingTraj          = "op_weight_msg_ipfs"
		defaultWeightMsgDeleteCodingTraj int = 100
	)

	var weightMsgDeleteCodingTraj int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCodingTraj, &weightMsgDeleteCodingTraj, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCodingTraj = defaultWeightMsgDeleteCodingTraj
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCodingTraj,
		ipfssimulation.SimulateMsgDeleteCodingTraj(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
