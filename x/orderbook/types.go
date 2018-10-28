package orderbook

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var minPrice, _ = sdk.NewDecFromStr("0.0000000001")
var maxPrice, _ = sdk.NewDecFromStr("1000000000")

func ValidPriceRatio(ratio sdk.Dec) bool {
	return ratio.GTE(minPrice) && ratio.LTE(maxPrice)
}

type Price struct {
	ratio            sdk.Dec
	numeratorDenom   string
	denomenatorDenom string
}

func NewPrice(ratio sdk.Dec, numeratorDenom, denomenatorDenom string) (price Price) {
	if !ValidPriceRatio(ratio) {
		return price
	}
	return Price{
		ratio:            ratio,
		numeratorDenom:   numeratorDenom,
		denomenatorDenom: denomenatorDenom,
	}
}

func (p Price) Reciprocal() Price {
	return Price{
		ratio:            SDKDecReciprocal(p.ratio),
		numeratorDenom:   p.denomenatorDenom,
		denomenatorDenom: p.numeratorDenom,
	}
}

func (p Price) Equal(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.Equal(p2.ratio)
}

func (p Price) LT(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.LT(p2.ratio)
}

func (p Price) GT(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.GT(p2.ratio)
}

func (p Price) LTE(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.LTE(p2.ratio)
}

func (p Price) GTE(p2 Price) bool {
	if p.numeratorDenom != p2.numeratorDenom || p.denomenatorDenom != p2.denomenatorDenom {
		panic("cannot compare prices of different units")
	}
	return p.ratio.GTE(p2.ratio)
}

func MulCoinsPrice(coins sdk.Coin, price Price) (sdk.Coin, error) {
	if coins.Denom != price.denomenatorDenom {
		return sdk.Coin{}, errors.New("Price and Coins incompatible")
	}
	return sdk.Coin{
		Amount: sdk.NewDecFromInt(coins.Amount).Mul(price.ratio).RoundInt(),
		Denom:  price.numeratorDenom,
	}, nil
}

type Order struct {
	orderId        int64
	owner          sdk.AccAddress
	sellCoins      sdk.Coin
	buyDenom       string
	price          Price
	expirationTime time.Time
}

func (o Order) Pair() DenomPair {
	return DenomPair{
		SellDenom: o.sellCoins.Denom,
		BuyDenom:  o.buyDenom,
	}
}

func (o Order) BidAmountAtAskPrice(askedPrice sdk.Dec) sdk.Coin {
	return sdk.Coin{
		Denom:  o.buyDenom,
		Amount: sdk.NewDecFromInt(o.sellCoins.Amount).Mul(askedPrice).RoundInt(),
	}
}

func (o Order) MatchOffer(offerAmount sdk.Coin) sdk.Coin {
	return sdk.Coin{
		Denom:  o.sellCoins.Denom,
		Amount: o.price.ratio.MulInt(offerAmount.Amount).RoundInt(),
	}
}

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

func (denomPair DenomPair) String() string {
	return fmt.Sprintf("%s - %s", denomPair.SellDenom, denomPair.BuyDenom)
}

func (denomPair DenomPair) ReversePair() DenomPair {
	return DenomPair{
		SellDenom: denomPair.BuyDenom,
		BuyDenom:  denomPair.SellDenom,
	}
}
