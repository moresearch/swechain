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
	"github.com/moresearch/swechain/x/issuemarket/keeper"
	"github.com/moresearch/swechain/x/issuemarket/types"
)

func createNBid(keeper keeper.Keeper, ctx context.Context, n int) []types.Bid {
	items := make([]types.Bid, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		_ = keeper.Bid.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestBidQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNBid(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetBidRequest
		response *types.QueryGetBidResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetBidRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetBidResponse{Bid: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetBidRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetBidResponse{Bid: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetBidRequest{
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
			response, err := qs.GetBid(f.ctx, tc.request)
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

func TestBidQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNBid(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllBidRequest {
		return &types.QueryAllBidRequest{
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
			resp, err := qs.ListBid(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bid), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bid),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListBid(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bid), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bid),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListBid(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Bid),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListBid(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
