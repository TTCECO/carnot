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
	CoinTTC = "CTTC" // TTC on Cosmos
	CoinACN = "CACN" // ACN on Cosmos, ERC20

	RouterName = "tcchan"
	StoreAcc   = "acc"
	StoreTCC   = "tcchan"

	prefixOrder   = "order"
	prefixPerson  = "person"
	prefixCurrent = "current"

	//rpcUrl = "http://rpc-tokyo.ttcnet.io" // Mainnet
	rpcUrl = "http://47.111.177.215:8511" // Testnet
	//defaultChainID = 8848                 // Mainnet
	defaultChainID = 8341 // Testnet

	contractAddress = "t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89"
)
