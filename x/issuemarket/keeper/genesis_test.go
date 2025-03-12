package keeper_test

import (
	"testing"

	"github.com/moresearch/swechain/testutil/nullify"
	"github.com/moresearch/swechain/x/issuemarket/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(), AuctionList: []types.Auction{{Id: 0}, {Id: 1}}, AuctionCount: 2, BidList: []types.Bid{{Index: "0"}, {Index: "1"}},
	}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.Params, got.Params)
	require.ElementsMatch(t, genesisState.AuctionList, got.AuctionList)
	require.Equal(t, genesisState.AuctionCount, got.AuctionCount)
	require.ElementsMatch(t, genesisState.BidList, got.BidList)

}
