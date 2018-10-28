package orderbook

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper - handlers sets/gets of custom variables for your module
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	// Reserved codespace
	codespace sdk.CodespaceType
}

var lastOrderIDKey = []byte("lastOrderID")
var ordersPrefix = []byte("orders")

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// AddNewOrder - Adds a new order into the proper orderbook
func (k Keeper) AddNewOrder(ctx sdk.Context, order Order) sdk.Error {
	if !ValidSortableDec(order.price.ratio) {
		return ErrInvalidPriceRange(k.codespace, order.price.ratio)
	}

	k.SetOrder(ctx, order)
	k.InsertOrderwallOrder(ctx, order)
	return nil
}

// Updates the amount of SellCoins left in an order
func (k Keeper) DecreaseOrderBidAmount(ctx sdk.Context, orderID int64, newAmount sdk.Coin) {
	order, found := k.GetOrder(ctx, orderID)
	if !found || !order.sellCoins.SameDenomAs(newAmount) || !newAmount.IsNotNegative() {
		return
	}

	if !newAmount.IsZero() {
		k.RemoveOrder(ctx, orderID)
		return
	}

	order.sellCoins = newAmount
	k.SetOrder(ctx, order)
}

// Removes an order from state and from its orderwall
func (k Keeper) RemoveOrder(ctx sdk.Context, orderID int64) Order {
	order, found := k.GetOrder(ctx, orderID)
	if !found {
		return Order{}
	}
	k.DeleteOrderwallOrder(ctx, order)
	k.DeleteOrder(ctx, orderID)

	return order
}
