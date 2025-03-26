package types

func NewMsgCreateBid(creator string, auctionId string, bidder string, amount string, description string) *MsgCreateBid {
	return &MsgCreateBid{
		Creator:     creator,
		AuctionId:   auctionId,
		Bidder:      bidder,
		Amount:      amount,
		Description: description,
	}
}

func NewMsgUpdateBid(creator string, id uint64, auctionId string, bidder string, amount string, description string) *MsgUpdateBid {
	return &MsgUpdateBid{
		Id:          id,
		Creator:     creator,
		AuctionId:   auctionId,
		Bidder:      bidder,
		Amount:      amount,
		Description: description,
	}
}

func NewMsgDeleteBid(creator string, id uint64) *MsgDeleteBid {
	return &MsgDeleteBid{
		Id:      id,
		Creator: creator,
	}
}
