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

	// Check if the value already exists
	ok, err := k.Bid.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var bid = types.Bid{
		Creator:     msg.Creator,
		Index:       msg.Index,
		AuctionId:   msg.AuctionId,
		Bidder:      msg.Bidder,
		Amount:      msg.Amount,
		Description: msg.Description,
	}

	if err := k.Bid.Set(ctx, bid.Index, bid); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateBidResponse{}, nil
}

func (k msgServer) UpdateBid(ctx context.Context, msg *types.MsgUpdateBid) (*types.MsgUpdateBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Bid.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var bid = types.Bid{
		Creator:     msg.Creator,
		Index:       msg.Index,
		AuctionId:   msg.AuctionId,
		Bidder:      msg.Bidder,
		Amount:      msg.Amount,
		Description: msg.Description,
	}

	if err := k.Bid.Set(ctx, bid.Index, bid); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update bid")
	}

	return &types.MsgUpdateBidResponse{}, nil
}

func (k msgServer) DeleteBid(ctx context.Context, msg *types.MsgDeleteBid) (*types.MsgDeleteBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Bid.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Bid.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove bid")
	}

	return &types.MsgDeleteBidResponse{}, nil
}
