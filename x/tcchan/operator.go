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

//curl -H 'content-type:application/json;' -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}' http://47.111.177.215:8511
import (
	"fmt"
	"github.com/TTCECO/gttc/accounts/keystore"
	"github.com/TTCECO/gttc/rpc"
	"github.com/tendermint/tendermint/libs/log"
	"io/ioutil"
)

type Operator struct {
	key    *keystore.Key
	cl     *rpc.Client
	logger log.Logger
}

func NewCrossChainOperator(logger log.Logger, keyfilepath string, password string) *Operator {
	operator := Operator{logger: logger}
	// unlock account
	if keyJson, err := ioutil.ReadFile(keyfilepath); err == nil {
		if operator.key, err = keystore.DecryptKey(keyJson, password); err == nil {
			logger.Info("Account unlock success ", "address", operator.key.Address.Hex())
		} else {
			logger.Error("Account unlock fail ", "error", err)
		}
	} else {
		fmt.Println("Keystore load fail ", "error", err)
	}
	// dial rpc
	if client, err := rpc.Dial(RPC_URL); err == nil {
		operator.cl = client
		logger.Info("Dial rpc success ", "url", RPC_URL)
	} else {
		fmt.Println("Dial rpc fail", "error", err)
	}
	return &operator
}
