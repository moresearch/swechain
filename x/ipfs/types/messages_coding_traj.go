package types

func NewMsgCreateCodingTraj(
	creator string,
	index string,
	title string,
	data string,

) *MsgCreateCodingTraj {
	return &MsgCreateCodingTraj{
		Creator: creator,
		Index:   index,
		Title:   title,
		Data:    data,
	}
}

func NewMsgUpdateCodingTraj(
	creator string,
	index string,
	title string,
	data string,

) *MsgUpdateCodingTraj {
	return &MsgUpdateCodingTraj{
		Creator: creator,
		Index:   index,
		Title:   title,
		Data:    data,
	}
}

func NewMsgDeleteCodingTraj(
	creator string,
	index string,

) *MsgDeleteCodingTraj {
	return &MsgDeleteCodingTraj{
		Creator: creator,
		Index:   index,
	}
}
