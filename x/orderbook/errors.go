//nolint
package orderbook

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 431

	CodeInvalidPrice sdk.CodeType = 1
)

//----------------------------------------
// Error constructors

func ErrInvalidPrice(codespace sdk.CodespaceType, priceRatio sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrice, fmt.Sprintf("Invalid Price %d. Must be between 10^10 & 10^-10.", priceRatio))
}
