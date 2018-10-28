package orderbook

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryOrder     = "order"
	QueryOrderwall = "orderwall"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryOrder:
			return queryOrder(ctx, path[1:], req, keeper)
		case QueryOrderwall:
			return queryOrderwall(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown orderbook query endpoint")
		}
	}
}

// nolint: unparam
func queryOrder(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	orderIDStr := path[0]

	orderID, err2 := strconv.ParseInt(orderIDStr, 10, 64)
	if err2 != nil {
		return res, ErrInvalidOrderID(keeper.codespace)
	}

	order, found := keeper.GetOrder(ctx, orderID)

	if !found {
		return res, ErrOrderNotFound(keeper.codespace, orderID)
	}

	res, err2 = codec.MarshalJSONIndent(keeper.cdc, order)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

// Whois represents a name -> value lookup
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// nolint: unparam
func queryOrderwall(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	denomPairStr := path[0]

	denomPair, err2 := DenomPairFromStr(denomPairStr)
	if err2 != nil {
		return res, ErrInvalidDenomPair(keeper.codespace)
	}

	var orderwall []Order

	orderwallIterator := keeper.OrderWallIterator(ctx, denomPair)

	for ; orderwallIterator.Valid(); orderwallIterator.Next() {
		var orderID int64
		keeper.cdc.MustUnmarshalBinary(orderwallIterator.Value(), &orderID)

		order, found := keeper.GetOrder(ctx, orderID)
		if found {
			orderwall = append(orderwall, order)
		}
	}

	orderwallIterator.Close()

	res, err2 = codec.MarshalJSONIndent(keeper.cdc, orderwall)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
