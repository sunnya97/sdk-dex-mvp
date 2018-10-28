package orderbook

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var expirationQueuePrefix = []byte("expirationQueue")

func ExpirationQueueTimesliceKey(timestamp time.Time) []byte {
	return AppendWithSeperator(expirationQueuePrefix, sdk.FormatTimeBytes(timestamp))
}

// gets a specific expiration queue timeslice. A timeslice is a slice of orderIDs of orders
// that expire at a certain time.
func (k Keeper) GetExpirationQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (orderIDs []int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(ExpirationQueueTimesliceKey(timestamp))
	if bz == nil {
		return []int64{}
	}
	k.cdc.MustUnmarshalBinary(bz, &orderIDs)
	return orderIDs
}

// Sets a specific expiration queue timeslice.
func (k Keeper) SetExpirationQueueTimeSlice(ctx sdk.Context, timestamp time.Time, orderIDs []int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(orderIDs)
	store.Set(ExpirationQueueTimesliceKey(timestamp), bz)
}

// Insert an orderID into the appropriate timeslice in the expiration queue
func (k Keeper) InsertExpirationQueue(ctx sdk.Context, orderID int64, expirationTime time.Time) {
	timeSlice := k.GetExpirationQueueTimeSlice(ctx, expirationTime)
	if len(timeSlice) == 0 {
		k.SetExpirationQueueTimeSlice(ctx, expirationTime, []int64{orderID})
	} else {
		timeSlice = append(timeSlice, orderID)
		k.SetExpirationQueueTimeSlice(ctx, expirationTime, timeSlice)
	}
}

// Returns all the expiration queue timeslices from time 0 until endTime
func (k Keeper) ExpirationQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(expirationQueuePrefix, sdk.InclusiveEndBytes(ExpirationQueueTimesliceKey(endTime)))
}

// Returns a concatenated list of all the timeslices before currTime, and deletes the timeslices from the queue
func (k Keeper) GetAllExpiredOrdersQueue(ctx sdk.Context, currTime time.Time) (expiredOrderIDs []int64) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	expiredTimesliceIterator := k.ExpirationQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; expiredTimesliceIterator.Valid(); expiredTimesliceIterator.Next() {
		timeslice := []int64{}
		k.cdc.MustUnmarshalBinary(expiredTimesliceIterator.Value(), &timeslice)
		expiredOrderIDs = append(expiredOrderIDs, timeslice...)
	}
	return expiredOrderIDs
}

// // Unbonds all the unbonding validators that have finished their unbonding period
// func (k Keeper) DeleteAllExpiredOrders(ctx sdk.Context) {
// 	expiredTimesliceIterator := k.ExpirationQueueIterator(ctx, ctx.BlockHeader().Time)
// 	for ; expiredTimesliceIterator.Valid(); expiredTimesliceIterator.Next() {
// 		timeslice := []int64{}
// 		k.cdc.MustUnmarshalBinary(expiredTimesliceIterator.Value(), &timeslice)
// 		for _, orderID := range timeslice {
// 			order, found := k.GetOrder(ctx, orderID)
// 			if !found {
// 				continue
// 			}
// 			k.DeleteOrder(orderID)
// 		}
// 		store.Delete(expiredTimesliceIterator.Key())
// 	}
// }
