package orderbook

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var orderwallPrefix = []byte("orderwalls")

// returns a prefix for storing all orders in the orderwall of a specific DenomPair
func OrderwallPrefix(pair DenomPair) []byte {
	return AppendWithSeperator(orderwallPrefix, []byte(pair.String()))
}

// Returns the key for getting an orderID in an orderWall
func OrderwallOrderKey(pair DenomPair, price Price, orderID int64) []byte {
	return AppendWithSeperator(AppendWithSeperator(OrderwallPrefix(pair), SortableSDKDecBytes(price.Ratio)), Int64ToSortableBytes(orderID))
}

// Returns an iterator for all the orders in an orderwall by price
func (k Keeper) OrderWallIterator(ctx sdk.Context, pair DenomPair) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(OrderwallPrefix(pair), sdk.PrefixEndBytes(OrderwallPrefix(pair)))
}

// peeks at the next lowest orderwall order
func (k Keeper) PeekOrderwallOrder(ctx sdk.Context, pair DenomPair) (order Order, found bool) {
	orderWall := k.OrderWallIterator(ctx, pair)
	if !orderWall.Valid() {
		return order, false
	}

	var orderID int64
	k.cdc.MustUnmarshalBinaryBare(orderWall.Value(), &orderID)

	orderWall.Close()

	return k.GetOrder(ctx, orderID)
}

// pops the next lowest orderwall order
func (k Keeper) PopOrderwallOrder(ctx sdk.Context, pair DenomPair) (order Order, found bool) {
	orderWall := k.OrderWallIterator(ctx, pair)
	if !orderWall.Valid() {
		return order, false
	}

	var orderID int64
	k.cdc.MustUnmarshalBinaryBare(orderWall.Value(), &orderID)

	orderWall.Close()

	store := ctx.KVStore(k.storeKey)
	store.Delete(orderWall.Key())

	return k.GetOrder(ctx, orderID)
}

// Insert an orderID into the appropriate timeslice in the expiration queue
func (k Keeper) InsertOrderwallOrder(ctx sdk.Context, order Order) {
	store := ctx.KVStore(k.storeKey)
	store.Set(OrderwallOrderKey(order.Pair(), order.Price, order.OrderID), k.cdc.MustMarshalBinaryBare(order.OrderID))
}

// Insert an orderID into the appropriate timeslice in the expiration queue
func (k Keeper) DeleteOrderwallOrder(ctx sdk.Context, order Order) {
	store := ctx.KVStore(k.storeKey)
	store.Set(OrderwallOrderKey(order.Pair(), order.Price, order.OrderID), nil)
}
