package keeper_test

import (
	"testing"

	"github.com/moresearch/swechain/testutil/nullify"
	"github.com/moresearch/swechain/x/ipfs/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(), CodingTrajList: []types.CodingTraj{{Index: "0"}, {Index: "1"}},
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
	require.ElementsMatch(t, genesisState.CodingTrajList, got.CodingTrajList)

}
