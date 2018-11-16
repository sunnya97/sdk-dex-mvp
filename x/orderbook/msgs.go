package orderbook

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg for creating a new order
// Price must be in units of BuyDenom/SellDenom
type MsgMakeOrder struct {
	OwnerAddr      sdk.AccAddress
	SellCoins      sdk.Coin
	Price          Price
	ExpirationTime time.Time
}

func NewMsgMakeOrder(ownerAddr sdk.AccAddress, sellCoins sdk.Coin, price Price, expirationTime time.Time) MsgMakeOrder {
	return MsgMakeOrder{
		OwnerAddr:      ownerAddr,
		SellCoins:      sellCoins,
		Price:          price,
		ExpirationTime: expirationTime,
	}
}

// Implements Msg.
func (msg MsgMakeOrder) Route() string { return "orderbook" }
func (msg MsgMakeOrder) Type() string  { return "add_order" }

// Implements Msg.
func (msg MsgMakeOrder) ValidateBasic() sdk.Error {
	if msg.OwnerAddr.Empty() {
		return sdk.ErrInvalidAddress(msg.OwnerAddr.String())
	}

	if !msg.SellCoins.IsPositive() {
		return sdk.ErrInvalidCoins(msg.SellCoins.String())
	}

	if len(msg.Price.NumeratorDenom) == 0 {
		return sdk.ErrInvalidCoins(msg.Price.NumeratorDenom)
	}

	if len(msg.Price.DenomenatorDenom) == 0 {
		return sdk.ErrInvalidCoins(msg.Price.DenomenatorDenom)
	}

	// Price must be in units of BuyDenom/SellDenom
	if msg.SellCoins.Denom != msg.Price.DenomenatorDenom {
		return ErrInvalidPriceFormat(DefaultCodespace, msg.Price)
	}

	if !ValidSortableDec(msg.Price.Ratio) {
		return ErrInvalidPriceRange(DefaultCodespace, msg.Price.Ratio)
	}

	return nil
}

// Implements Msg.
func (msg MsgMakeOrder) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgMakeOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddr}
}

type MsgRemoveOrder struct {
	OwnerAddr sdk.AccAddress
	OrderID   int64
}

func NewMsgRemoveOrder(ownerAddr sdk.AccAddress, orderID int64) MsgRemoveOrder {
	return MsgRemoveOrder{
		OwnerAddr: ownerAddr,
		OrderID:   orderID,
	}
}

// Implements Msg.
func (msg MsgRemoveOrder) Route() string { return "orderbook" }
func (msg MsgRemoveOrder) Type() string  { return "remove_order" }

// Implements Msg.
func (msg MsgRemoveOrder) ValidateBasic() sdk.Error {
	if msg.OwnerAddr.Empty() {
		return sdk.ErrInvalidAddress(msg.OwnerAddr.String())
	}

	if msg.OrderID < 0 {
		return sdk.ErrInternal(fmt.Sprintf("%d", msg.OrderID))
	}

	return nil
}

// Implements Msg.
func (msg MsgRemoveOrder) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgRemoveOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddr}
}
