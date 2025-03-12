package keeper

import (
	"context"

	"github.com/moresearch/swechain/x/issuemarket/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.AuctionList {
		if err := k.Auction.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.AuctionSeq.Set(ctx, genState.AuctionCount); err != nil {
		return err
	}
	for _, elem := range genState.BidList {
		if err := k.Bid.Set(ctx, elem.Index, elem); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = k.Auction.Walk(ctx, nil, func(key uint64, elem types.Auction) (bool, error) {
		genesis.AuctionList = append(genesis.AuctionList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.AuctionCount, err = k.AuctionSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.Bid.Walk(ctx, nil, func(_ string, val types.Bid) (stop bool, err error) {
		genesis.BidList = append(genesis.BidList, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
