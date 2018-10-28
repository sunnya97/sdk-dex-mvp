package cli

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunnya97/sdk-dex-mvp/x/orderbook"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

// GetCmdMakeOrder is the CLI command for sending a MakeOrder transaction
func GetCmdMakeOrder(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "make-order [sellcoins] @ [priceratio] [numerDenom] / [denomDenom]",
		Short: "make an order for selling coins for another coin at a certain price",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			sellCoins, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			priceRatio, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			numerDenom := args[2]
			denomDenom := args[3]

			if numerDenom != sellCoins.Denom && denomDenom != sellCoins.Denom || numerDenom == denomDenom {
				return orderbook.ErrInvalidDenomPair(orderbook.DefaultCodespace)
			}

			price := orderbook.NewPrice(priceRatio, args[2], args[3])

			if denomDenom != sellCoins.Denom {
				price = price.Reciprocal()
			}

			msg := orderbook.NewMsgMakeOrder(account, sellCoins, price, time.Time{})
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

// GetCmdMakeOrder is the CLI command for sending a MakeOrder transaction
func GetCmdRemoveOrder(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-order [orderID]",
		Short: "make an order for selling coins for another coin at a certain price",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			orderID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := orderbook.NewMsgRemoveOrder(account, orderID)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
