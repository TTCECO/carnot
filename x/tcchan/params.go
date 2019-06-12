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
	//RPCURL= "http://rpc-tokyo.ttcnet.io" // Mainnet
	RPCURL = "http://47.111.177.215:8511" // Testnet
	//defaultChainID = 8848                 // Mainnet
	defaultChainID = 8341 // Testnet

	ContractAddress = "t0C4B13EA1219a6FfA1A092707db7071377c0E4374" // test address

	minValidatorCount = 2

	blockDelay = 15
)
