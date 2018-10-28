package orderbook

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "orderexecution" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMakeOrder:
			return handleMsgMakeOrder(ctx, keeper, msg)
		case MsgRemoveOrder:
			return handleMsgRemoveOrder(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized orderbook Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgMakeOrder
func handleMsgMakeOrder(ctx sdk.Context, keeper Keeper, msg MsgMakeOrder) sdk.Result {
	order := Order{
		orderID:        keeper.GetNextOrderID(ctx),
		owner:          msg.MakerAddr,
		sellCoins:      msg.SellCoins,
		price:          msg.Price,
		expirationTime: msg.ExpirationTime,
	}

	keeper.coinKeeper.SubtractCoins(ctx, order.owner, sdk.Coins{order.sellCoins})
	keeper.AddNewOrder(ctx, order)

	return sdk.Result{}
}

// Handle MsgRemoveOrder
func handleMsgRemoveOrder(ctx sdk.Context, keeper Keeper, msg MsgRemoveOrder) sdk.Result {
	removedOrder := keeper.RemoveOrder(ctx, msg.OrderID)
	keeper.coinKeeper.AddCoins(ctx, removedOrder.owner, sdk.Coins{removedOrder.sellCoins})

	return sdk.Result{}
}
