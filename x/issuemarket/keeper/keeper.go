package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/moresearch/swechain/x/issuemarket/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema     collections.Schema
	Params     collections.Item[types.Params]
	AuctionSeq collections.Sequence
	Auction    collections.Map[uint64, types.Auction]
	BidSeq     collections.Sequence
	Bid        collections.Map[uint64, types.Bid]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)), AuctionSeq: collections.NewSequence(sb, types.AuctionCountKey, "auction_seq"), Auction: collections.NewMap(sb, types.AuctionKey, "auction", collections.Uint64Key, codec.CollValue[types.Auction](cdc)), BidSeq: collections.NewSequence(sb, types.BidCountKey, "bid_seq"), Bid: collections.NewMap(sb, types.BidKey, "bid", collections.Uint64Key, codec.CollValue[types.Bid](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
