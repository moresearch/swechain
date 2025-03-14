package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/moresearch/swechain/x/ipfs/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListCodingTraj(ctx context.Context, req *types.QueryAllCodingTrajRequest) (*types.QueryAllCodingTrajResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	codingTrajs, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.CodingTraj,
		req.Pagination,
		func(_ string, value types.CodingTraj) (types.CodingTraj, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCodingTrajResponse{CodingTraj: codingTrajs, Pagination: pageRes}, nil
}

func (q queryServer) GetCodingTraj(ctx context.Context, req *types.QueryGetCodingTrajRequest) (*types.QueryGetCodingTrajResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.CodingTraj.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetCodingTrajResponse{CodingTraj: val}, nil
}
