package orderbook

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Price object that has an sdk.Dec for the price along with units
// Supports comparision operators, that panic if not used to compare prices of the same units
type Price struct {
	ratio            sdk.Dec
	numeratorDenom   string
	denomenatorDenom string
}

// Checks the make sure that the ratio sdk.Dec is within sortable bounds
func NewPrice(ratio sdk.Dec, numeratorDenom, denomenatorDenom string) (price Price) {
	if !ValidSortableDec(ratio) {
		return price
	}
	return Price{
		ratio:            ratio,
		numeratorDenom:   numeratorDenom,
		denomenatorDenom: denomenatorDenom,
	}
}

// Returns the reciprocal price (1/ratio) as well as flips the denoms
func (p Price) Reciprocal() Price {
	return Price{
		ratio:            SDKDecReciprocal(p.ratio),
		numeratorDenom:   p.denomenatorDenom,
		denomenatorDenom: p.numeratorDenom,
	}
}

// nolint
func (p Price) Equal(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.Equal(p2.ratio)
}

// nolint
func (p Price) LT(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.LT(p2.ratio)
}

// nolint
func (p Price) GT(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.GT(p2.ratio)
}

// nolint
func (p Price) LTE(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.LTE(p2.ratio)
}

// nolint
func (p Price) GTE(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.GTE(p2.ratio)
}

// Returns the result of converting an sdk.Coin using a conversion ratio price
func MulCoinsPrice(coins sdk.Coin, price Price) (sdk.Coin, error) {
	if coins.Denom != price.denomenatorDenom {
		return sdk.Coin{}, errors.New("Price and Coins incompatible")
	}
	return sdk.Coin{
		Amount: sdk.NewDecFromInt(coins.Amount).Mul(price.ratio).RoundInt(),
		Denom:  price.numeratorDenom,
	}, nil
}

// ------------------------------------------------------------

// Order
type Order struct {
	orderID        int64
	owner          sdk.AccAddress
	sellCoins      sdk.Coin
	buyDenom       string
	price          Price
	expirationTime time.Time
}

// Returns the DenomPair of (BuyDenom, SellDenom).  Used for assigning order to the proper orderbook
func (o Order) Pair() DenomPair {
	return DenomPair{
		SellDenom: o.sellCoins.Denom,
		BuyDenom:  o.buyDenom,
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

// returns a string form of Denom Pair, is just the two denom strings, separated by " - " delimiter
func (denomPair DenomPair) String() string {
	return fmt.Sprintf("%s - %s", denomPair.SellDenom, denomPair.BuyDenom)
}

// Returns the inverse DenomPair with BuyDenom and SellDenom switched
func (denomPair DenomPair) ReversePair() DenomPair {
	return DenomPair{
		SellDenom: denomPair.BuyDenom,
		BuyDenom:  denomPair.SellDenom,
	}
}
