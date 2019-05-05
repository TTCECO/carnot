// Copyright 2019 The TTC Authors
// This file is part of the TTC library.
//
// The TTC library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The TTC library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the TTC library. If not, see <http://www.gnu.org/licenses/>.

package tcchan

import (
	"github.com/TTCECO/gttc/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// Returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// CCToC is the struct that contains cross chain transaction from TTC Mainnet to cosmos SDK
type CCToC struct {
	From  common.Address `json:"from"`
	To    sdk.AccAddress `json:"to"`
	Value sdk.Coin       `json:"value"`
}

// CCFromC is the struct that contains cross chain transation from cosmos SDK to TTC Mainnet
type CCFromC struct {
	From  sdk.AccAddress `json:"from"`
	To    common.Address `json:"to"`
	Value sdk.Coin       `json:"value"`
}
