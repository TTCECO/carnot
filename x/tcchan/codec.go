package tcchan

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "tcchan/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "tcchan/BuyName", nil)
}
