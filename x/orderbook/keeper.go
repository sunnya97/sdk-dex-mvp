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

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
		codespace:  codespace,
	}
}

// AddNewOrder - Adds a new order into the proper orderbook
func (k Keeper) AddNewOrder(ctx sdk.Context, order Order) (consumed bool, err sdk.Error) {
	if !ValidSortableDec(order.Price.Ratio) {
		return false, ErrInvalidPriceRange(k.codespace, order.Price.Ratio)
	}

	// First run order against opposing order wall
	order, consumed = k.ExecuteOrderAgainstOrderWall(ctx, order)

	// if the order hasn't been fully executed, add it to its own order wall
	if !consumed {
		k.SetOrder(ctx, order)
		k.InsertOrderwallOrder(ctx, order)
	}
	return consumed, nil
}

// Updates the amount of SellCoins left in an order
func (k Keeper) DecreaseOrderBidAmount(ctx sdk.Context, orderID int64, newAmount sdk.Coin) {
	order, found := k.GetOrder(ctx, orderID)
	if !found || !order.SellCoins.SameDenomAs(newAmount) || !newAmount.IsNotNegative() {
		return
	}

	if !newAmount.IsZero() {
		k.RemoveOrder(ctx, orderID)
		return
	}

	order.SellCoins = newAmount
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

// Executes an order against an orderwall until either the order is fully consumed, there are no more order left in the wall,
// or there is a spread (the prices don't overlap)
func (k Keeper) ExecuteOrderAgainstOrderWall(ctx sdk.Context, order Order) (remainingOrder Order, consumed bool) {
	opposingPair := order.Pair().ReversePair()

	// while the order hasn't been fully consumed
	for order.SellCoins.IsPositive() {

		// get the first order in the opposing wall
		// if there are no more orders in the opposing wall, end by placing the remaining incoming order in its own wall
		peekWallOrder, found := k.PeekOrderwallOrder(ctx, opposingPair)
		if !found {
			k.DecreaseOrderBidAmount(ctx, order.OrderID, order.SellCoins)
			return order, false
		}

		// get the asking price of peekedOrder
		askPrice := peekWallOrder.Price.Reciprocal()

		// If the asking price is greater than the incoming order is willing to pay, break out of the loop and end
		if order.Price.LT(askPrice) {
			break
		}

		// get the amount the taker has to pay to execute the entire peekedOrder *at the maker's price*
		bidAtAskingPrice, _ := MulCoinsPrice(order.SellCoins, askPrice)

		// if the peeked order can't fulfill my entire order, execute as much as possible (the entire peeked order)
		// and remove the peeked order
		if bidAtAskingPrice.IsGTE(peekWallOrder.SellCoins) {
			// the amount that the taker has to pay to complete the peekedOrder
			executeAmount, _ := MulCoinsPrice(peekWallOrder.SellCoins, askPrice.Reciprocal())

			// Remove executeAmount from the incoming order's sellCoins and send them to the peekOrder's maker
			k.coinKeeper.AddCoins(ctx, peekWallOrder.Owner, sdk.Coins{executeAmount})
			order.SellCoins = order.SellCoins.Minus(executeAmount)

			// Send the full sellCoins of the peekedOrder to the incoming order's owner (the taker)
			// and remove the peeked order from state
			k.coinKeeper.AddCoins(ctx, order.Owner, sdk.Coins{peekWallOrder.SellCoins})
			k.RemoveOrder(ctx, peekWallOrder.OrderID)
		} else {
			// scenario that peekedOrder is larger than the incoming taker order

			// amount that the peekedOrder trades to fully execute the incoming order
			executeAmount, _ := MulCoinsPrice(order.SellCoins, askPrice)

			// Remove executeAmount from the peekedOrder's sellCoins and send them to the taker (the incoming order's owner)
			k.coinKeeper.AddCoins(ctx, order.Owner, sdk.Coins{executeAmount})
			peekWallOrder.SellCoins = peekWallOrder.SellCoins.Minus(executeAmount)

			// send all the coins in the taker's order to the maker,
			// remove the taker's order as it's been completely fulfilled,
			// and return with consumed as true, as the entire incoming order has been consumed
			k.coinKeeper.AddCoins(ctx, peekWallOrder.Owner, sdk.Coins{order.SellCoins})
			order.SellCoins = order.SellCoins.Minus(order.SellCoins)
			k.RemoveOrder(ctx, order.OrderID)
			return order, true
		}
	}

	// Breaks out of loop when the order hasn't been completely executed, but there are no more
	// overlapping price orders

	// Set the decreased coins left in state
	k.DecreaseOrderBidAmount(ctx, order.OrderID, order.SellCoins)
	// return that the order has not been completely consumed
	return order, false
}
