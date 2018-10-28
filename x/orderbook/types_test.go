package orderbook

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPrice_Reciprocal(t *testing.T) {
	type fields struct {
		ratioStr         string
		numeratorDenom   string
		denomenatorDenom string
	}
	tests := []struct {
		name   string
		fields fields
		want   fields
	}{
		{"simple", fields{"0.25", "BTC", "ETH"}, fields{"4", "ETH", "BTC"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio, _ := sdk.NewDecFromStr(tt.fields.ratioStr)
			p := NewPrice(
				ratio,
				tt.fields.numeratorDenom,
				tt.fields.denomenatorDenom,
			)

			ratio2, _ := sdk.NewDecFromStr(tt.want.ratioStr)
			p2 := NewPrice(
				ratio2,
				tt.want.numeratorDenom,
				tt.want.denomenatorDenom,
			)

			if got := p.Reciprocal(); !reflect.DeepEqual(got, p2) {
				t.Errorf("Price.Reciprocal() = %v, want %v", got, p2)
			}
		})
	}
}

func TestMulCoinsPrice(t *testing.T) {
	type args struct {
		coins            sdk.Coin
		priceRatioStr    string
		priceNumerator   string
		priceDenomenator string
	}
	tests := []struct {
		name    string
		args    args
		want    sdk.Coin
		wantErr bool
	}{
		{"simple", args{sdk.NewInt64Coin("BTC", 5), "0.2", "USD", "BTC"}, sdk.NewInt64Coin("USD", 1), false},
		{"invalid", args{sdk.NewInt64Coin("BTC", 5), "0.2", "BTC", "USD"}, sdk.Coin{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ratio, _ := sdk.NewDecFromStr(tt.args.priceRatioStr)
			price := NewPrice(ratio, tt.args.priceNumerator, tt.args.priceDenomenator)

			got, err := MulCoinsPrice(tt.args.coins, price)
			if (err != nil) != tt.wantErr {
				t.Errorf("MulCoinsPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MulCoinsPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
