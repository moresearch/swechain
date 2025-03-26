package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/moresearch/swechain/x/issuemarket/types"
)

func (k msgServer) CreateBid(ctx context.Context, msg *types.MsgCreateBid) (*types.MsgCreateBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	nextId, err := k.BidSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	var bid = types.Bid{
		Id:          nextId,
		Creator:     msg.Creator,
		AuctionId:   msg.AuctionId,
		Bidder:      msg.Bidder,
		Amount:      msg.Amount,
		Description: msg.Description,
	}

	if err = k.Bid.Set(
		ctx,
		nextId,
		bid,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set bid")
	}

	return &types.MsgCreateBidResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdateBid(ctx context.Context, msg *types.MsgUpdateBid) (*types.MsgUpdateBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	var bid = types.Bid{
		Creator:     msg.Creator,
		Id:          msg.Id,
		AuctionId:   msg.AuctionId,
		Bidder:      msg.Bidder,
		Amount:      msg.Amount,
		Description: msg.Description,
	}

	// Checks that the element exists
	val, err := k.Bid.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get bid")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Bid.Set(ctx, msg.Id, bid); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update bid")
	}

	return &types.MsgUpdateBidResponse{}, nil
}

func (k msgServer) DeleteBid(ctx context.Context, msg *types.MsgDeleteBid) (*types.MsgDeleteBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Bid.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get bid")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Bid.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete bid")
	}

	return &types.MsgDeleteBidResponse{}, nil
}
