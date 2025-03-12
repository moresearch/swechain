package types

func NewMsgCreateAuction(creator string, issue string, description string, status string, winner string) *MsgCreateAuction {
	return &MsgCreateAuction{
		Creator:     creator,
		Issue:       issue,
		Description: description,
		Status:      status,
		Winner:      winner,
	}
}

func NewMsgUpdateAuction(creator string, id uint64, issue string, description string, status string, winner string) *MsgUpdateAuction {
	return &MsgUpdateAuction{
		Id:          id,
		Creator:     creator,
		Issue:       issue,
		Description: description,
		Status:      status,
		Winner:      winner,
	}
}

func NewMsgDeleteAuction(creator string, id uint64) *MsgDeleteAuction {
	return &MsgDeleteAuction{
		Id:      id,
		Creator: creator,
	}
}
