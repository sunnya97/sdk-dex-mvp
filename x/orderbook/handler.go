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
	orderID := keeper.GetNextOrderID(ctx)

	order := Order{
		OrderID:        orderID,
		Owner:          msg.OwnerAddr,
		SellCoins:      msg.SellCoins,
		Price:          msg.Price,
		ExpirationTime: msg.ExpirationTime,
	}

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, order.Owner, sdk.Coins{order.SellCoins})
	if err != nil {
		return err.Result()
	}

	consumed, err := keeper.AddNewOrder(ctx, order)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data: keeper.cdc.MustMarshalBinaryBare(consumed),
	}
}

// Handle MsgRemoveOrder
func handleMsgRemoveOrder(ctx sdk.Context, keeper Keeper, msg MsgRemoveOrder) sdk.Result {
	removedOrder := keeper.RemoveOrder(ctx, msg.OrderID)
	keeper.coinKeeper.AddCoins(ctx, removedOrder.Owner, sdk.Coins{removedOrder.SellCoins})

	return sdk.Result{}
}
