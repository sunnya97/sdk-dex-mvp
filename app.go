package app

import (
	"github.com/sunnya97/sdk-dex-mvp/x/orderbook"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	appName = "Dexter"
)

type DexterApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyOrderbook     *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey

	accountKeeper   auth.AccountKeeper
	feeKeeper       auth.FeeCollectionKeeper
	bankKeeper      bank.Keeper
	orderbookKeeper orderbook.Keeper

	codespacer *sdk.Codespacer
}

func NewDexterApp(logger log.Logger, db dbm.DB) *DexterApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &DexterApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:      sdk.NewKVStoreKey("main"),
		keyAccount:   sdk.NewKVStoreKey("acc"),
		keyOrderbook: sdk.NewKVStoreKey("orderbook"),
	}

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		auth.ProtoBaseAccount,
	)

	app.feeKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper)

	app.orderbookKeeper = orderbook.NewKeeper(
		app.bankKeeper,
		app.keyOrderbook,
		app.cdc,
		app.RegisterCodespace(orderbook.DefaultCodespace),
	)

	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeKeeper))

	app.Router().
		AddRoute("orderbook", orderbook.NewHandler(app.orderbookKeeper)).
		AddRoute("bank", bank.NewHandler(app.bankKeeper))

	app.SetInitChainer(app.initChainer)

	app.MountStoresIAVL(
		app.keyMain,
		app.keyAccount,
		app.keyOrderbook,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

type GenesisState struct {
	Accounts []auth.BaseAccount `json:"accounts"`
}

func (app *DexterApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, &acc)
	}

	return abci.ResponseInitChain{}
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	orderbook.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
