package orderbook

// import (
// 	"bytes"
// 	"log"
// 	"sort"
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	abci "github.com/tendermint/tendermint/abci/types"
// 	"github.com/tendermint/tendermint/crypto"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/bank"
// 	"github.com/cosmos/cosmos-sdk/x/mock"
// )

// // initialize the mock application for this module
// func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, bank.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
// 	mapp := mock.NewApp()

// 	RegisterCodec(mapp.Cdc)

// 	keyOrderbook := sdk.NewKVStoreKey("orderbook")

// 	ck := bank.NewBaseKeeper(mapp.AccountKeeper)
// 	keeper := NewKeeper(ck, keyOrderbook, mapp.Cdc, DefaultCodespace)

// 	mapp.Router().AddRoute("orderbook", NewHandler(keeper))

// 	mapp.SetInitChainer(getInitChainer(mapp))

// 	require.NoError(t, mapp.CompleteSetup(keyOrderbook))

// 	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin("steak", 42)})

// 	mock.SetGenesis(mapp, genAccs)

// 	return mapp, keeper, ck, addrs, pubKeys, privKeys
// }

// // gov and stake initchainer
// func getInitChainer(mapp *mock.App) sdk.InitChainer {
// 	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
// 		mapp.InitChainer(ctx, req)
// 	}
// }

// // TODO: Remove once address interface has been implemented (ref: #2186)
// func SortValAddresses(addrs []sdk.ValAddress) {
// 	var byteAddrs [][]byte
// 	for _, addr := range addrs {
// 		byteAddrs = append(byteAddrs, addr.Bytes())
// 	}

// 	SortByteArrays(byteAddrs)

// 	for i, byteAddr := range byteAddrs {
// 		addrs[i] = byteAddr
// 	}
// }

// // Sorts Addresses
// func SortAddresses(addrs []sdk.AccAddress) {
// 	var byteAddrs [][]byte
// 	for _, addr := range addrs {
// 		byteAddrs = append(byteAddrs, addr.Bytes())
// 	}
// 	SortByteArrays(byteAddrs)
// 	for i, byteAddr := range byteAddrs {
// 		addrs[i] = byteAddr
// 	}
// }

// // implement `Interface` in sort package.
// type sortByteArrays [][]byte

// func (b sortByteArrays) Len() int {
// 	return len(b)
// }

// func (b sortByteArrays) Less(i, j int) bool {
// 	// bytes package already implements Comparable for []byte.
// 	switch bytes.Compare(b[i], b[j]) {
// 	case -1:
// 		return true
// 	case 0, 1:
// 		return false
// 	default:
// 		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
// 		return false
// 	}
// }

// func (b sortByteArrays) Swap(i, j int) {
// 	b[j], b[i] = b[i], b[j]
// }

// // Public
// func SortByteArrays(src [][]byte) [][]byte {
// 	sorted := sortByteArrays(src)
// 	sort.Sort(sorted)
// 	return sorted
// }
