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

const (
	CoinTTC = "ttc" // TTC on Cosmos
	CoinACN = "acn" // ACN on Cosmos, ERC20

	RouterName = "tcchan"
	StoreTCC   = "tcchan"

	prefixDeposit = "deposit"
	prefixPerson  = "person"
	prefixCurrent = "current"
	prefixConfirm = "confirm"

	minBalanceValue = 1e+18 // for gas
	//rpcUrl = "http://rpc-tokyo.ttcnet.io" // Mainnet
	rpcUrl = "http://47.111.177.215:8511" // Testnet
	//defaultChainID = 8848                 // Mainnet
	defaultChainID = 8341 // Testnet

	contractAddress = "t0f6086b9588d4a74636d8df61f4bd9ab2e8eeb7f9" // test address

	minValidatorCount = 2

	blockDelay = 15
)
