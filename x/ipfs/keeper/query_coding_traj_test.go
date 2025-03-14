package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/moresearch/swechain/testutil/nullify"
	"github.com/moresearch/swechain/x/ipfs/keeper"
	"github.com/moresearch/swechain/x/ipfs/types"
)

func createNCodingTraj(keeper keeper.Keeper, ctx context.Context, n int) []types.CodingTraj {
	items := make([]types.CodingTraj, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		_ = keeper.CodingTraj.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestCodingTrajQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNCodingTraj(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetCodingTrajRequest
		response *types.QueryGetCodingTrajResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetCodingTrajRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetCodingTrajResponse{CodingTraj: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetCodingTrajRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetCodingTrajResponse{CodingTraj: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetCodingTrajRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetCodingTraj(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestCodingTrajQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNCodingTraj(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCodingTrajRequest {
		return &types.QueryAllCodingTrajRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListCodingTraj(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.CodingTraj), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.CodingTraj),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListCodingTraj(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.CodingTraj), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.CodingTraj),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListCodingTraj(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.CodingTraj),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListCodingTraj(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
