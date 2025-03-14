package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/moresearch/swechain/x/ipfs/types"
)

func (k msgServer) CreateCodingTraj(ctx context.Context, msg *types.MsgCreateCodingTraj) (*types.MsgCreateCodingTrajResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.CodingTraj.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var codingTraj = types.CodingTraj{
		Creator: msg.Creator,
		Index:   msg.Index,
		Title:   msg.Title,
		Data:    msg.Data,
	}

	if err := k.CodingTraj.Set(ctx, codingTraj.Index, codingTraj); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateCodingTrajResponse{}, nil
}

func (k msgServer) UpdateCodingTraj(ctx context.Context, msg *types.MsgUpdateCodingTraj) (*types.MsgUpdateCodingTrajResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.CodingTraj.Get(ctx, msg.Index)
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

	var codingTraj = types.CodingTraj{
		Creator: msg.Creator,
		Index:   msg.Index,
		Title:   msg.Title,
		Data:    msg.Data,
	}

	if err := k.CodingTraj.Set(ctx, codingTraj.Index, codingTraj); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update codingTraj")
	}

	return &types.MsgUpdateCodingTrajResponse{}, nil
}

func (k msgServer) DeleteCodingTraj(ctx context.Context, msg *types.MsgDeleteCodingTraj) (*types.MsgDeleteCodingTrajResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.CodingTraj.Get(ctx, msg.Index)
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

	if err := k.CodingTraj.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove codingTraj")
	}

	return &types.MsgDeleteCodingTrajResponse{}, nil
}
