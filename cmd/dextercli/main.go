package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	app "github.com/sunnya97/sdk-dex-mvp"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	orderbookcmd "github.com/sunnya97/sdk-dex-mvp/x/orderbook/cli"
)

const storeAcc = "acc"
const storeNSnames = "orderbook"

var (
	rootCmd = &cobra.Command{
		Use:   "dextercli",
		Short: "Dexter Client",
	}
	DefaultCLIHome = os.ExpandEnv("$HOME/.dextercli")
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd.AddCommand(client.ConfigCmd())
	rpc.AddCommands(rootCmd)

	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(queryCmd, cdc)
	queryCmd.AddCommand(client.LineBreak)
	queryCmd.AddCommand(client.GetCommands(
		authcmd.GetAccountCmd(storeAcc, cdc, authcmd.GetAccountDecoder(cdc)),
		orderbookcmd.GetCmdGetOrder("orderbook", cdc),
		orderbookcmd.GetCmdGetOrderwall("orderbook", cdc),
	)...)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		orderbookcmd.GetCmdMakeOrder(cdc),
		orderbookcmd.GetCmdRemoveOrder(cdc),
	)...)

	rootCmd.AddCommand(
		queryCmd,
		txCmd,
		client.LineBreak,
	)

	rootCmd.AddCommand(
		keys.Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "DEX", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
