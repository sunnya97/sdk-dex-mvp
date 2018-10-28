package orderbook

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Price object that has an sdk.Dec for the price along with units
// Supports comparision operators, that panic if not used to compare prices of the same units
type Price struct {
	Ratio            sdk.Dec
	NumeratorDenom   string
	DenomenatorDenom string
}

// Checks the make sure that the ratio sdk.Dec is within sortable bounds
func NewPrice(ratio sdk.Dec, numeratorDenom, denomenatorDenom string) (price Price) {
	if !ValidSortableDec(ratio) {
		return price
	}
	return Price{
		Ratio:            ratio,
		NumeratorDenom:   numeratorDenom,
		DenomenatorDenom: denomenatorDenom,
	}
}

// Returns the reciprocal price (1/ratio) as well as flips the denoms
func (p Price) Reciprocal() Price {
	return Price{
		Ratio:            SDKDecReciprocal(p.Ratio),
		NumeratorDenom:   p.DenomenatorDenom,
		DenomenatorDenom: p.NumeratorDenom,
	}
}

// nolint
func (p Price) Equal(p2 Price) bool {
	if p.NumeratorDenom != p2.NumeratorDenom || p.DenomenatorDenom != p2.DenomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.Ratio.Equal(p2.Ratio)
}

// nolint
func (p Price) LT(p2 Price) bool {
	if p.NumeratorDenom != p2.NumeratorDenom || p.DenomenatorDenom != p2.DenomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.Ratio.LT(p2.Ratio)
}

// nolint
func (p Price) GT(p2 Price) bool {
	if p.NumeratorDenom != p2.NumeratorDenom || p.DenomenatorDenom != p2.DenomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.Ratio.GT(p2.Ratio)
}

// nolint
func (p Price) LTE(p2 Price) bool {
	if p.NumeratorDenom != p2.NumeratorDenom || p.DenomenatorDenom != p2.DenomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.Ratio.LTE(p2.Ratio)
}

// nolint
func (p Price) GTE(p2 Price) bool {
	if p.NumeratorDenom != p2.NumeratorDenom || p.DenomenatorDenom != p2.DenomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.Ratio.GTE(p2.Ratio)
}

// Returns the result of converting an sdk.Coin using a conversion ratio price
func MulCoinsPrice(coins sdk.Coin, price Price) (sdk.Coin, error) {
	if coins.Denom != price.DenomenatorDenom {
		return sdk.Coin{}, errors.New("Price and Coins incompatible")
	}
	return sdk.Coin{
		Amount: sdk.NewDecFromInt(coins.Amount).Mul(price.Ratio).RoundInt(),
		Denom:  price.NumeratorDenom,
	}, nil
}

// ------------------------------------------------------------

// Order
type Order struct {
	OrderID        int64
	Owner          sdk.AccAddress
	SellCoins      sdk.Coin
	BuyDenom       string
	Price          Price
	ExpirationTime time.Time
}

// Returns the DenomPair of (BuyDenom, SellDenom).  Used for assigning order to the proper orderbook
func (o Order) Pair() DenomPair {
	return DenomPair{
		SellDenom: o.SellCoins.Denom,
		BuyDenom:  o.BuyDenom,
	}
}

// DenomPair is a tuple of two denoms
type DenomPair struct {
	SellDenom string
	BuyDenom  string
}

func NewDenomPair(sellDenom, buyDenom string) DenomPair {
	return DenomPair{
		SellDenom: sellDenom,
		BuyDenom:  buyDenom,
	}
}

// Returns a DenomPair from its string representation
func DenomPairFromStr(str string) (DenomPair, error) {
	denomStrings := strings.Split(str, "|")
	if len(denomStrings) != 2 {
		return DenomPair{}, errors.New("Incorrectly formatted DenomPair string")
	}
	return DenomPair{
		SellDenom: denomStrings[0],
		BuyDenom:  denomStrings[1],
	}, nil
}

// returns a string form of Denom Pair, is just the two denom strings, separated by " - " delimiter
func (denomPair DenomPair) String() string {
	return fmt.Sprintf("%s|%s", denomPair.SellDenom, denomPair.BuyDenom)
}

// Returns the inverse DenomPair with BuyDenom and SellDenom switched
func (denomPair DenomPair) ReversePair() DenomPair {
	return DenomPair{
		SellDenom: denomPair.BuyDenom,
		BuyDenom:  denomPair.SellDenom,
	}
}
