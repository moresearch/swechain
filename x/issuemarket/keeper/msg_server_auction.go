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

func (k msgServer) CreateAuction(ctx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	nextId, err := k.AuctionSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	var auction = types.Auction{
		Id:          nextId,
		Creator:     msg.Creator,
		Issue:       msg.Issue,
		Description: msg.Description,
		Status:      msg.Status,
		Winner:      msg.Winner,
	}

	if err = k.Auction.Set(
		ctx,
		nextId,
		auction,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set auction")
	}

	return &types.MsgCreateAuctionResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdateAuction(ctx context.Context, msg *types.MsgUpdateAuction) (*types.MsgUpdateAuctionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	var auction = types.Auction{
		Creator:     msg.Creator,
		Id:          msg.Id,
		Issue:       msg.Issue,
		Description: msg.Description,
		Status:      msg.Status,
		Winner:      msg.Winner,
	}

	// Checks that the element exists
	val, err := k.Auction.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get auction")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Auction.Set(ctx, msg.Id, auction); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update auction")
	}

	return &types.MsgUpdateAuctionResponse{}, nil
}

func (k msgServer) DeleteAuction(ctx context.Context, msg *types.MsgDeleteAuction) (*types.MsgDeleteAuctionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Auction.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get auction")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Auction.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete auction")
	}

	return &types.MsgDeleteAuctionResponse{}, nil
}
