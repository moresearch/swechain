package types_test

import (
	"testing"

	"github.com/moresearch/swechain/x/issuemarket/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{AuctionList: []types.Auction{{Id: 0}, {Id: 1}}, AuctionCount: 2, BidList: []types.Bid{{Index: "0"}, {Index: "1"}}},
			valid:    true,
		}, {desc: "duplicated auction",

			genState: &types.GenesisState{AuctionList: []types.Auction{{Id: 0},
				{

					Id: 0}}, BidList: []types.Bid{{Index: "0"}, {Index: "1"}},
			}, valid: false}, {desc: "invalid auction count",

			genState: &types.GenesisState{AuctionList: []types.Auction{{Id: 1}}, AuctionCount: 0, BidList: []types.Bid{{Index: "0"}, {Index: "1"}}}, valid: false}, {desc: "duplicated bid",

			genState: &types.GenesisState{BidList: []types.Bid{{Index: "0"}, {Index: "0"}}}, valid: false},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
