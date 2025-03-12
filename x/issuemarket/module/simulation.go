package issuemarket

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/moresearch/swechain/testutil/sample"
	issuemarketsimulation "github.com/moresearch/swechain/x/issuemarket/simulation"
	"github.com/moresearch/swechain/x/issuemarket/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	issuemarketGenesis := types.GenesisState{
		Params: types.DefaultParams(), AuctionList: []types.Auction{{Id: 0, Creator: sample.AccAddress()}, {Id: 1, Creator: sample.AccAddress()}}, AuctionCount: 2, BidList: []types.Bid{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&issuemarketGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateAuction          = "op_weight_msg_issuemarket"
		defaultWeightMsgCreateAuction int = 100
	)

	var weightMsgCreateAuction int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAuction, &weightMsgCreateAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAuction = defaultWeightMsgCreateAuction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAuction,
		issuemarketsimulation.SimulateMsgCreateAuction(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateAuction          = "op_weight_msg_issuemarket"
		defaultWeightMsgUpdateAuction int = 100
	)

	var weightMsgUpdateAuction int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAuction, &weightMsgUpdateAuction, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAuction = defaultWeightMsgUpdateAuction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAuction,
		issuemarketsimulation.SimulateMsgUpdateAuction(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteAuction          = "op_weight_msg_issuemarket"
		defaultWeightMsgDeleteAuction int = 100
	)

	var weightMsgDeleteAuction int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAuction, &weightMsgDeleteAuction, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAuction = defaultWeightMsgDeleteAuction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAuction,
		issuemarketsimulation.SimulateMsgDeleteAuction(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateBid          = "op_weight_msg_issuemarket"
		defaultWeightMsgCreateBid int = 100
	)

	var weightMsgCreateBid int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateBid, &weightMsgCreateBid, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBid = defaultWeightMsgCreateBid
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateBid,
		issuemarketsimulation.SimulateMsgCreateBid(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateBid          = "op_weight_msg_issuemarket"
		defaultWeightMsgUpdateBid int = 100
	)

	var weightMsgUpdateBid int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateBid, &weightMsgUpdateBid, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateBid = defaultWeightMsgUpdateBid
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateBid,
		issuemarketsimulation.SimulateMsgUpdateBid(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteBid          = "op_weight_msg_issuemarket"
		defaultWeightMsgDeleteBid int = 100
	)

	var weightMsgDeleteBid int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteBid, &weightMsgDeleteBid, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteBid = defaultWeightMsgDeleteBid
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteBid,
		issuemarketsimulation.SimulateMsgDeleteBid(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
