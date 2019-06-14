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

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TTCECO/carnot/x/tcchan"
	"github.com/TTCECO/carnot/x/tcchan/contract"
	"github.com/TTCECO/gttc/accounts/abi/bind"
	"github.com/TTCECO/gttc/common"
	"github.com/TTCECO/gttc/core/types"
	"github.com/TTCECO/gttc/crypto"
	"github.com/TTCECO/gttc/ethclient"
	"github.com/TTCECO/gttc/rlp"
	"github.com/TTCECO/gttc/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
	"time"
)

func main() {
	viper.New()
	rootCmd := &cobra.Command{
		Use:   "carnotcd",
		Short: "Carnot Contract Deploy Tools",
	}
	rootCmd.AddCommand(DeployCmd())
	rootCmd.AddCommand(InitContractSettingCmd())
	rootCmd.AddCommand(WithdrawTxCmd())
	rootCmd.AddCommand(WithdrawTokenCmd())
	rootCmd.AddCommand(DisplayTxCmd())

	rootCmd.Execute()
}

func DeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy contract by default settings",
		RunE: func(_ *cobra.Command, _ []string) error {
			client, err := rpc.Dial(tcchan.RPCURL)
			if err != nil {
				fmt.Println("rpc.Dial err", err)
				return err
			}
			var result string
			err = client.Call(&result, "eth_getTransactionCount", Address, "latest")
			if err != nil {
				fmt.Println("client.nonce err", err)
				return err
			}
			nonce, err := strconv.ParseInt(result[2:], 16, 64)
			if err != nil {
				fmt.Println("nonce parse fail", err)
				return err
			}
			fmt.Printf("nonce : %d\n", nonce)

			err = client.Call(&result, "net_version")
			if err != nil {
				fmt.Println("get chain id fail", err)
				return err
			}
			fmt.Printf("chainId: %s\n", result)

			chainID, err := strconv.ParseInt(result, 10, 64)
			if err != nil {
				fmt.Println("parse chain id fail", err)
				return err
			}

			privateKey, err := crypto.HexToECDSA(PrivateKey)
			if err != nil {
				fmt.Println("create private key err :", err)
				return err
			}

			res := make(map[string]interface{})
			err = json.Unmarshal([]byte(byteCode), &res)
			if err != nil {
				fmt.Println("json unmarshal fail", err)
				return err
			}
			byteData, err := hex.DecodeString(res["object"].(string))
			if err != nil {
				fmt.Println("decode string fail ", err)
				return err
			}
			tx := types.NewContractCreation(uint64(nonce), big.NewInt(0), uint64(30000000), big.NewInt(21000000), byteData)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
			txData, err := rlp.EncodeToBytes(signedTx)
			if err != nil {
				fmt.Println("rlp Encode fail", err)
				return err
			}
			if err := client.Call(&result, "eth_sendRawTransaction", common.ToHex(txData)); err != nil {
				fmt.Println("send Transaction fail", err)
				return err
			}
			fmt.Println("send Transaction tx : ", result)
			waitSeconds := 20
			fmt.Println("wait ", waitSeconds, " seconds for receipt ")
			for i := 0; i < waitSeconds; i++ {
				time.Sleep(time.Second)
				fmt.Println(waitSeconds - i)
			}
			var receiptResult types.Receipt
			if err := client.Call(&receiptResult, "eth_getTransactionReceipt", result); err != nil {
				fmt.Println("get Transaction Receipt fail", err)
				return err
			}

			if receiptResult.Status > 0 {
				fmt.Println("contract deploy address : ", receiptResult.ContractAddress.Hex())
				fmt.Println("set ContractAddress in carnot/x/tcchan/params by this. ")
			} else {
				fmt.Println("contract deploy fail")
				return err
			}

			return nil
		},
	}
	return cmd
}

func InitContractSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [contract address]",
		Short: "Initialize contract ",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) error {

			client, err := rpc.Dial(tcchan.RPCURL)
			if err != nil {
				fmt.Println("rpc.Dial err", err)
				return err
			}

			privateKey, err := crypto.HexToECDSA(PrivateKey)
			if err != nil {
				fmt.Println("create private key err :", err)
				return err
			}

			ctx := context.Background()
			cl := ethclient.NewClient(client)

			contractAddress := tcchan.ContractAddress
			if len(args) > 0 {
				contractAddress = args[0]
			}
			contract, err := contract.NewContract(common.HexToAddress(contractAddress), cl)
			if err != nil {
				fmt.Println("initialize contract fail : ", err)
				return err
			}

			for _, v := range validators {
				tx, err := contract.AddValidator(bind.NewKeyedTransactor(privateKey), common.HexToAddress(v))
				if err != nil {
					return err
				}
				receipt, err := bind.WaitMined(ctx, cl, tx)
				if err != nil {
					return err
				}
				if receipt.Status != 1 {
					return err
				}
				fmt.Println("add validator sucess : ", v)
			}

			minConfirmNum := big.NewInt(2)
			if _, err := contract.SetMinConfirmNum(bind.NewKeyedTransactor(privateKey), minConfirmNum); err != nil {
				fmt.Println("setMinConfirmNum fail : ", err)
				return err
			} else {
				fmt.Println("setMinConfirmNum success to : ", minConfirmNum)
			}

			if _, err := contract.AddSupportToken(bind.NewKeyedTransactor(privateKey), TokenName, common.HexToAddress(TokenAddress)); err != nil {
				fmt.Println("AddSupportToken fail : ", err)
				return err
			} else {
				fmt.Println("AddSupportToken success : ", TokenName)
			}

			return nil
		},
	}

	return cmd
}

func WithdrawTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [target accAddress] [amount]",
		Short: "Call withdraw func on contract (send ttc from ttc mainnet to cosmos)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {

			targetAddress := args[0]
			ttcAmount, ok := new(big.Int).SetString(args[1], 10)
			if !ok {
				fmt.Println("amount is not number")
				return errors.New("is not number")
			}
			fmt.Println(ttcAmount)

			client, err := rpc.Dial(tcchan.RPCURL)
			if err != nil {
				fmt.Println("rpc.Dial err", err)
				return err
			}
			privateKey, err := crypto.HexToECDSA(PrivateKey)
			if err != nil {
				fmt.Println("create private key err :", err)
				return err
			}
			ctx := context.Background()
			cl := ethclient.NewClient(client)
			contract, err := contract.NewContract(common.HexToAddress(tcchan.ContractAddress), cl)
			if err != nil {
				fmt.Println("initialize contract fail : ", err)
				return err
			}

			txOpts := bind.NewKeyedTransactor(privateKey)
			txOpts.Value = ttcAmount
			tx, err := contract.CrossChainTransactionCoin(txOpts, targetAddress)
			if err != nil {
				return err
			}
			receipt, err := bind.WaitMined(ctx, cl, tx)
			if err != nil {
				return err
			}
			if receipt.Status != 1 {
				return err
			}
			fmt.Println("send coin success .")
			return nil
		},
	}

	return cmd
}

func WithdrawTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawToken [target accAddress] [amount]",
		Short: "Call withdraw func on contract (send acn from ttc mainnet to cosmos)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {

			targetAddress := args[0]
			tokenAmount, ok := new(big.Int).SetString(args[1], 10)
			if !ok {
				fmt.Println("amount is not number")
				return errors.New("is not number")
			}
			fmt.Println(tokenAmount)

			client, err := rpc.Dial(tcchan.RPCURL)
			if err != nil {
				fmt.Println("rpc.Dial err", err)
				return err
			}
			privateKey, err := crypto.HexToECDSA(acnPrivateKey)
			if err != nil {
				fmt.Println("create private key err :", err)
				return err
			}
			ctx := context.Background()
			cl := ethclient.NewClient(client)
			contract, err := contract.NewContract(common.HexToAddress(tcchan.ContractAddress), cl)
			if err != nil {
				fmt.Println("initialize contract fail : ", err)
				return err
			}

			txOpts := bind.NewKeyedTransactor(privateKey)
			tx, err := contract.CrossChainTransactionToken(txOpts, common.HexToAddress(TokenAddress), TokenName, targetAddress, tokenAmount)
			if err != nil {
				return err
			}
			receipt, err := bind.WaitMined(ctx, cl, tx)
			if err != nil {
				return err
			}
			if receipt.Status != 1 {
				return err
			}
			fmt.Println("send coin success .")
			return nil
		},
	}

	return cmd
}

func DisplayTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display",
		Short: "Display the information on contract",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, args []string) error {

			client, err := rpc.Dial(tcchan.RPCURL)
			if err != nil {
				fmt.Println("rpc.Dial err", err)
				return err
			}
			cl := ethclient.NewClient(client)
			contract, err := contract.NewContract(common.HexToAddress(tcchan.ContractAddress), cl)
			if err != nil {
				fmt.Println("initialize contract fail : ", err)
				return err
			}

			callOpts := &bind.CallOpts{}
			if res, err := contract.MinConfirmNum(callOpts); err != nil {
				return err
			} else {
				fmt.Println("MinConfirmNum is ", res)
			}

			if res, err := contract.WithdrawOrderID(callOpts); err != nil {
				return err
			} else {
				fmt.Println("WithdrawOrderID is ", res)
				for i := big.NewInt(1); i.Cmp(res) <= 0; i = i.Add(i, big.NewInt(1)) {
					if res, err := contract.WithdrawRecords(callOpts, i); err != nil {
						break
					} else {
						fmt.Println("WithdrawOrders ", i, " : ", res)
					}
				}
			}

			if res, err := contract.DepositFee(callOpts); err != nil {
				return err
			} else {
				fmt.Println("DepositFee is ", res)
			}
			for i := big.NewInt(0); ; {
				if res, err := contract.DepositKeys(callOpts, i); err != nil {
					break
				} else {
					fmt.Println("DepositKeys ", i, " : ", res)
					i = i.Add(i, big.NewInt(1))
				}
			}
			return nil
		},
	}

	return cmd
}
