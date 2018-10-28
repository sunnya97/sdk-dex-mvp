package orderbook

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgMakeOrder{}, "orderbook/MakeOrder", nil)
	cdc.RegisterConcrete(MsgRemoveOrder{}, "orderbook/RemoveOrder", nil)
}
