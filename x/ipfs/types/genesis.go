package types

import

// DefaultGenesis returns the default genesis state
"fmt"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(), CodingTrajList: []CodingTraj{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	codingTrajIndexMap := make(map[string]struct{})

	for _, elem := range gs.CodingTrajList {
		index := fmt.Sprint(elem.Index)
		if _, ok := codingTrajIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for codingTraj")
		}
		codingTrajIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
