package orderbook

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// get key in store to get an Order
func OrderKey(orderID int64) []byte {
	return AppendWithSeperator(ordersPrefix, Int64ToSortableBytes(orderID))
}

// Gets an Order from the Store
func (k Keeper) GetOrder(ctx sdk.Context, orderID int64) (order Order, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(OrderKey(orderID))
	if bz == nil {
		return order, false
	}
	err := k.cdc.UnmarshalBinary(bz, &order)
	if err != nil {
		return order, false
	}
	return order, true
}

// Sets an Order int the Store
func (k Keeper) SetOrder(ctx sdk.Context, order Order) {
	store := ctx.KVStore(k.storeKey)
	store.Set(OrderKey(order.orderId), k.cdc.MustMarshalBinary(order))
}

// Sets an Order int the Store
func (k Keeper) DeleteOrder(ctx sdk.Context, orderID int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(OrderKey(orderID), nil)
}

// Gets the last orderID that was assigned
func (k Keeper) GetLastOrderID(ctx sdk.Context) (lastOrderID int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastOrderIDKey)
	if bz == nil {
		k.SetLastOrderID(ctx, 0)
		return 0
	}
	k.cdc.UnmarshalBinary(bz, &lastOrderID)
	return lastOrderID
}

// Sets the last orderID that was assigned
func (k Keeper) SetLastOrderID(ctx sdk.Context, orderID int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(orderID)
	store.Set(lastOrderIDKey, bz)
}

// Gets the next unassigned orderID (and increments lastOrderID)
func (k Keeper) GetNextOrderID(ctx sdk.Context) (nextOrderID int64) {
	lastOrderID := k.GetLastOrderID(ctx)
	nextOrderID = lastOrderID + 1
	k.SetLastOrderID(ctx, nextOrderID)
	return nextOrderID
}
