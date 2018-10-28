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
		OrderID:        keeper.GetNextOrderID(ctx),
		Owner:          msg.OwnerAddr,
		SellCoins:      msg.SellCoins,
		Price:          msg.Price,
		ExpirationTime: msg.ExpirationTime,
	}

	keeper.coinKeeper.SubtractCoins(ctx, order.Owner, sdk.Coins{order.SellCoins})
	keeper.AddNewOrder(ctx, order)

	return sdk.Result{}
}

// Handle MsgRemoveOrder
func handleMsgRemoveOrder(ctx sdk.Context, keeper Keeper, msg MsgRemoveOrder) sdk.Result {
	removedOrder := keeper.RemoveOrder(ctx, msg.OrderID)
	keeper.coinKeeper.AddCoins(ctx, removedOrder.Owner, sdk.Coins{removedOrder.SellCoins})

	return sdk.Result{}
}
