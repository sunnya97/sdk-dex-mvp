package orderbook

import (
	"bytes"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var keySeperator = []byte("/")

// joins two byte slices using the keySeperator as a delimter
func AppendWithSeperator(byteslice1, byteslice2 []byte) []byte {
	return append(append(byteslice1, keySeperator...), byteslice2...)
}

// Splits a key path using the delimter keySeperator
func SplitKeyAlongSeperator(fullKey []byte) [][]byte {
	return bytes.Split(fullKey, keySeperator)
}

func Int64ToSortableBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func SDKDecReciprocal(dec sdk.Dec) sdk.Dec {
	return sdk.OneDec().Quo(dec)
}
