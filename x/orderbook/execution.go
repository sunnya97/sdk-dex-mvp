package orderbook

import sdk "github.com/cosmos/cosmos-sdk/types"

// Executes an order against an orderwall until either the order is fully consumed, there are no more order left in the wall,
// or there is a spread (the prices don't overlap)
func (k Keeper) ExecuteOrderAgainstOrderWall(ctx sdk.Context, order Order) (remainingOrder Order, consumed bool) {
	opposingPair := order.Pair().ReversePair()

	// while the order hasn't been fully consumed
	for order.sellCoins.IsPositive() {

		// get the first order in the opposing wall
		// if there are no more orders in the opposing wall, end by placing the remaining incoming order in its own wall
		peekWallOrder, found := k.PeekOrderwallOrder(ctx, opposingPair)
		if !found {
			k.DecreaseOrderBidAmount(ctx, order.orderID, order.sellCoins)
			return order, false
		}

		// get the asking price of peekedOrder
		askPrice := peekWallOrder.price.Reciprocal()

		// If the asking price is greater than the incoming order is willing to pay, break out of the loop and end
		if order.price.LT(askPrice) {
			break
		}

		// get the amount the taker has to pay to execute the entire peekedOrder *at the maker's price*
		bidAtAskingPrice, _ := MulCoinsPrice(order.sellCoins, askPrice)

		// if the peeked order can't fulfill my entire order, execute as much as possible (the entire peeked order)
		// and remove the peeked order
		if bidAtAskingPrice.IsGTE(peekWallOrder.sellCoins) {
			// the amount that the taker has to pay to complete the peekedOrder
			executeAmount, _ := MulCoinsPrice(peekWallOrder.sellCoins, askPrice.Reciprocal())

			// Remove executeAmount from the incoming order's sellCoins and send them to the peekOrder's maker
			k.coinKeeper.AddCoins(ctx, peekWallOrder.owner, sdk.Coins{executeAmount})
			order.sellCoins = order.sellCoins.Minus(executeAmount)

			// Send the full sellCoins of the peekedOrder to the incoming order's owner (the taker)
			// and remove the peeked order from state
			k.coinKeeper.AddCoins(ctx, order.owner, sdk.Coins{peekWallOrder.sellCoins})
			k.RemoveOrder(ctx, peekWallOrder.orderID)
		} else {
			// scenario that peekedOrder is larger than the incoming taker order

			// amount that the peekedOrder trades to fully execute the incoming order
			executeAmount, _ := MulCoinsPrice(order.sellCoins, askPrice)

			// Remove executeAmount from the peekedOrder's sellCoins and send them to the taker (the incoming order's owner)
			k.coinKeeper.AddCoins(ctx, order.owner, sdk.Coins{executeAmount})
			peekWallOrder.sellCoins = peekWallOrder.sellCoins.Minus(executeAmount)

			// send all the coins in the taker's order to the maker,
			// remove the taker's order as it's been completely fulfilled,
			// and return with consumed as true, as the entire incoming order has been consumed
			k.coinKeeper.AddCoins(ctx, peekWallOrder.owner, sdk.Coins{order.sellCoins})
			order.sellCoins = order.sellCoins.Minus(order.sellCoins)
			k.RemoveOrder(ctx, order.orderID)
			return order, true
		}
	}

	// Breaks out of loop when the order hasn't been completely executed, but there are no more
	// overlapping price orders

	// Set the decreased coins left in state
	k.DecreaseOrderBidAmount(ctx, order.orderID, order.sellCoins)
	// return that the order has not been completely consumed
	return order, false
}
