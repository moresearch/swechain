package ipfs

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/moresearch/swechain/x/ipfs/types"
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
					RpcMethod: "ListCodingTraj",
					Use:       "list-coding-traj",
					Short:     "List all coding_traj",
				},
				{
					RpcMethod:      "GetCodingTraj",
					Use:            "get-coding-traj [id]",
					Short:          "Gets a coding_traj",
					Alias:          []string{"show-coding-traj"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
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
					RpcMethod:      "CreateCodingTraj",
					Use:            "create-coding-traj [index] [title] [data]",
					Short:          "Create a new coding_traj",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "title"}, {ProtoField: "data"}},
				},
				{
					RpcMethod:      "UpdateCodingTraj",
					Use:            "update-coding-traj [index] [title] [data]",
					Short:          "Update coding_traj",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "title"}, {ProtoField: "data"}},
				},
				{
					RpcMethod:      "DeleteCodingTraj",
					Use:            "delete-coding-traj [index]",
					Short:          "Delete coding_traj",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
