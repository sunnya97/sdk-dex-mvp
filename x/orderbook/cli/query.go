package cli

import (
	"fmt"

	"github.com/sunnya97/sdk-dex-mvp/x/orderbook"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetCmdGetOrder queries information about a name
func GetCmdGetOrder(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "order [orderID]",
		Short: "get order by orderID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			orderIDStr := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/order/%s", queryRoute, orderIDStr), nil)
			if err != nil {
				fmt.Printf("could not find order with orderID %s \n", orderIDStr)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// GetCmdWhois queries information about a domain
func GetCmdGetOrderwall(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "orderwall [sellDenom] [buyDenom]",
		Short: "Get orderwall of a specific pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sellDenom := args[0]
			buyDenom := args[1]

			denomPair := orderbook.DenomPair{
				SellDenom: sellDenom,
				BuyDenom:  buyDenom,
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/orderwall/%s", queryRoute, denomPair.String()), nil)
			if err != nil {
				fmt.Printf("could not find orderwall \n")
				return nil
			}

			var orders []orderbook.Order

			cdc.UnmarshalJSON(res, &orders)

			for _, order := range orders {
				fmt.Printf("%d - %s @ %d %s/%s \n",
					order.OrderID,
					order.SellCoins,
					order.Price.Ratio,
					order.Price.NumeratorDenom,
					order.Price.DenomenatorDenom,
				)
			}

			return nil
		},
	}
}
