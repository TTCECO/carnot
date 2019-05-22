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

package client

import (
	"github.com/TTCECO/carnot/x/tcchan"
	tcchancmd "github.com/TTCECO/carnot/x/tcchan/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

// NewModuleClient create the new client
func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group tcchan queries under a subcommand
	tcchanQueryCmd := &cobra.Command{
		Use:   tcchan.RouterName,
		Short: "Querying commands for the tcchan module",
	}

	tcchanQueryCmd.AddCommand(client.GetCommands(
		tcchancmd.GetCmdOrder(mc.storeKey, mc.cdc),
		tcchancmd.GetCmdPerson(mc.storeKey, mc.cdc),
		tcchancmd.GetCmdCurrent(mc.storeKey, mc.cdc),
		tcchancmd.GetCmdConfirm(mc.storeKey, mc.cdc),
	)...)

	return tcchanQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	tcchanTxCmd := &cobra.Command{
		Use:   tcchan.RouterName,
		Short: "Transactions subcommands",
	}

	tcchanTxCmd.AddCommand(client.PostCommands(
		tcchancmd.GetCmdDeposit(mc.cdc),
		tcchancmd.GetCmdWithdrawConfirm(mc.cdc),
	)...)

	return tcchanTxCmd
}
