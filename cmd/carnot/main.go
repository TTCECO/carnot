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
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"strconv"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/TTCECO/carnot/app"
	carnotInit "github.com/TTCECO/carnot/init"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// carnot custom flags
const flagAssertInvariantsBlockly = "assert-invariants-blockly"

var assertInvariantsBlockly bool

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "carnot",
		Short:             "Carnot Daemon (server)",
		PersistentPreRunE: app.PersistentPreRunEFn(ctx),
	}
	rootCmd.AddCommand(carnotInit.InitCmd(ctx, cdc))
	rootCmd.AddCommand(carnotInit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(carnotInit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(CrossChainStartCmd(ctx, newApp))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GN", app.DefaultNodeHome)
	rootCmd.PersistentFlags().BoolVar(&assertInvariantsBlockly, flagAssertInvariantsBlockly,
		false, "Assert registered invariants on a blockly basis")
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewCarnotApp(
		logger, db, traceStore, true, assertInvariantsBlockly,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		gApp := app.NewCarnotApp(logger, db, traceStore, false, false)
		err := gApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return gApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	gApp := app.NewCarnotApp(logger, db, traceStore, true, false)
	return gApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func CrossChainStartCmd(ctx *server.Context, appCreator server.AppCreator) *cobra.Command {
	cmd := server.StartCmd(ctx, appCreator)
	cmd.Use = "cc-start [keystore.json] [password] [validator name] [validator pass]"
	cmd.Short = "Run the full node (with unlock TTC address for cross chain action)"
	cmd.Args = cobra.ExactArgs(4)
	oldRunE := cmd.RunE
	cmd.RunE = func(cc *cobra.Command, args []string) error {

		app.InitKeystore = args[0]
		app.InitPassword = args[1]
		app.ValidatorName = args[2]
		app.ValidatorPass = args[3]
		app.KeyPath = viper.GetString("home-client")
		listenAddress := viper.GetString("rpc.laddr")
		if res := strings.Split(listenAddress, ":"); len(res) > 0 {
			if port, err := strconv.Atoi(res[len(res)-1]); err == nil {
				app.RPCPort = port
			}
		}
		return oldRunE(cc, args)
	}
	cmd.Flags().String("home-client", app.DefaultCLIHome, "client's home directory")
	return cmd

}
