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

package cli

import (
	"fmt"

	"github.com/TTCECO/carnot/x/tcchan"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetCmdOrder queries information about a order
func GetCmdOrder(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "order [orderID]",
		Short: "Query deposit order info by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			orderID := args[0]
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposit/%s", queryRoute, orderID), nil)
			if err != nil {
				fmt.Printf("could not resolve order - %d : %s\n", orderID, err)
				return nil
			}

			var out tcchan.DepositOrder
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdPerson queries information about a address
func GetCmdPerson(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "person [address]",
		Short: "Query person tx records by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/person/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("could not resolve address - %s : %s\n", address, err)
				return nil
			}
			var out tcchan.PersonalOrderRecord
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdCurrent queries information
func GetCmdCurrent(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Query current info",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/current", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get current info : %s\n", err)
				return nil
			}
			var out tcchan.CurrentOrderRecord
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}


// GetCmdConfirm queries information
func GetCmdConfirm(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "confirm",
		Short: "Query confirm info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			confirmID := args[0]
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/confirm/%s", queryRoute, confirmID), nil)
			if err != nil {
				fmt.Printf("could not resolve order - %d : %s\n", confirmID, err)
				return nil
			}

			var out tcchan.WithdrawConfirm
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
