package types

func NewMsgCreateBid(
	creator string,
	index string,
	auctionId string,
	bidder string,
	amount string,
	description string,

) *MsgCreateBid {
	return &MsgCreateBid{
		Creator:     creator,
		Index:       index,
		AuctionId:   auctionId,
		Bidder:      bidder,
		Amount:      amount,
		Description: description,
	}
}

func NewMsgUpdateBid(
	creator string,
	index string,
	auctionId string,
	bidder string,
	amount string,
	description string,

) *MsgUpdateBid {
	return &MsgUpdateBid{
		Creator:     creator,
		Index:       index,
		AuctionId:   auctionId,
		Bidder:      bidder,
		Amount:      amount,
		Description: description,
	}
}

func NewMsgDeleteBid(
	creator string,
	index string,

) *MsgDeleteBid {
	return &MsgDeleteBid{
		Creator: creator,
		Index:   index,
	}
}
