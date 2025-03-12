package keeper_test

import (
	"context"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/moresearch/swechain/testutil/nullify"
	"github.com/moresearch/swechain/x/issuemarket/keeper"
	"github.com/moresearch/swechain/x/issuemarket/types"
)

func createNAuction(keeper keeper.Keeper, ctx context.Context, n int) []types.Auction {
	items := make([]types.Auction, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		_ = keeper.Auction.Set(ctx, iu, items[i])
		_ = keeper.AuctionSeq.Set(ctx, iu)
	}
	return items
}

func TestAuctionQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAuction(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetAuctionRequest
		response *types.QueryGetAuctionResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetAuctionRequest{Id: msgs[0].Id},
			response: &types.QueryGetAuctionResponse{Auction: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetAuctionRequest{Id: msgs[1].Id},
			response: &types.QueryGetAuctionResponse{Auction: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetAuctionRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetAuction(f.ctx, tc.request)
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

func TestAuctionQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAuction(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAuctionRequest {
		return &types.QueryAllAuctionRequest{
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
			resp, err := qs.ListAuction(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auction), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Auction),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAuction(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auction), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Auction),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListAuction(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Auction),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListAuction(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
