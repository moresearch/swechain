package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/moresearch/swechain/x/issuemarket/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListBid(ctx context.Context, req *types.QueryAllBidRequest) (*types.QueryAllBidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	bids, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Bid,
		req.Pagination,
		func(_ uint64, value types.Bid) (types.Bid, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllBidResponse{Bid: bids, Pagination: pageRes}, nil
}

func (q queryServer) GetBid(ctx context.Context, req *types.QueryGetBidRequest) (*types.QueryGetBidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	bid, err := q.k.Bid.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetBidResponse{Bid: bid}, nil
}
