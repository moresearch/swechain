package issuemarket

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/moresearch/swechain/x/issuemarket/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListAuction",
					Use:       "list-auction",
					Short:     "List all auction",
				},
				{
					RpcMethod:      "GetAuction",
					Use:            "get-auction [id]",
					Short:          "Gets a auction by id",
					Alias:          []string{"show-auction"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ListBid",
					Use:       "list-bid",
					Short:     "List all bid",
				},
				{
					RpcMethod:      "GetBid",
					Use:            "get-bid [id]",
					Short:          "Gets a bid by id",
					Alias:          []string{"show-bid"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateAuction",
					Use:            "create-auction [issue] [description] [status] [winner]",
					Short:          "Create auction",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "issue"}, {ProtoField: "description"}, {ProtoField: "status"}, {ProtoField: "winner"}},
				},
				{
					RpcMethod:      "UpdateAuction",
					Use:            "update-auction [id] [issue] [description] [status] [winner]",
					Short:          "Update auction",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "issue"}, {ProtoField: "description"}, {ProtoField: "status"}, {ProtoField: "winner"}},
				},
				{
					RpcMethod:      "DeleteAuction",
					Use:            "delete-auction [id]",
					Short:          "Delete auction",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "CreateBid",
					Use:            "create-bid [auctionId] [bidder] [amount] [description]",
					Short:          "Create bid",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auctionId"}, {ProtoField: "bidder"}, {ProtoField: "amount"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "UpdateBid",
					Use:            "update-bid [id] [auctionId] [bidder] [amount] [description]",
					Short:          "Update bid",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "auctionId"}, {ProtoField: "bidder"}, {ProtoField: "amount"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "DeleteBid",
					Use:            "delete-bid [id]",
					Short:          "Delete bid",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
