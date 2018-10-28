package orderbook

import (
	"bytes"
	"encoding/binary"
	"fmt"

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

// marshals int64 to a bigendian byte slice so it can be sorted
func Int64ToSortableBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

// returns the reciprocal of an sdk.Dec
func SDKDecReciprocal(dec sdk.Dec) sdk.Dec {
	return sdk.OneDec().Quo(dec)
}

var maxDec, _ = sdk.NewDecFromStr("10000000000")

// Ensures that an sdk.Dec is within the sortable bounds
// sdk.Dec can't have precision of less that 10^-10
func ValidSortableDec(dec sdk.Dec) bool {
	return dec.LTE(maxDec)
}

// Returns a byte slice representation of an sdk.Dec that can be sorted.
// Left and right pads with 0s so there are 10 digits to left and right of decimal point
// For this reason, there is a maximum and minimum value for this
// Prices need to be marshalled using this, and so prices must be within the bounds
// enforced by ValidSortableDec
func SortableSDKDecBytes(dec sdk.Dec) []byte {
	if !ValidSortableDec(dec) {
		panic("dec must be within bounds")
	}
	return []byte(fmt.Sprintf("%020s", dec))
}
