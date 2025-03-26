package types

import

// DefaultGenesis returns the default genesis state
"fmt"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(), AuctionList: []Auction{}, BidList: []Bid{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	auctionIdMap := make(map[uint64]bool)
	auctionCount := gs.GetAuctionCount()
	for _, elem := range gs.AuctionList {
		if _, ok := auctionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for auction")
		}
		if elem.Id >= auctionCount {
			return fmt.Errorf("auction id should be lower or equal than the last id")
		}
		auctionIdMap[elem.Id] = true
	}
	bidIdMap := make(map[uint64]bool)
	bidCount := gs.GetBidCount()
	for _, elem := range gs.BidList {
		if _, ok := bidIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for bid")
		}
		if elem.Id >= bidCount {
			return fmt.Errorf("bid id should be lower or equal than the last id")
		}
		bidIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
