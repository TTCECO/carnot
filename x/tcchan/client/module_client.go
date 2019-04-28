package client

import (
	tcchancmd "github.com/TTCECO/ttc-cosmos-channal/x/tcchan/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group tcchan queries under a subcommand
	namesvcQueryCmd := &cobra.Command{
		Use:   "tcchan",
		Short: "Querying commands for the tcchan module",
	}

	namesvcQueryCmd.AddCommand(client.GetCommands(
		tcchancmd.GetCmdResolveName(mc.storeKey, mc.cdc),
		tcchancmd.GetCmdWhois(mc.storeKey, mc.cdc),
		tcchancmd.GetCmdNames(mc.storeKey, mc.cdc),
	)...)

	return namesvcQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	namesvcTxCmd := &cobra.Command{
		Use:   "tcchan",
		Short: "tcchan transactions subcommands",
	}

	namesvcTxCmd.AddCommand(client.PostCommands(
		tcchancmd.GetCmdBuyName(mc.cdc),
		tcchancmd.GetCmdSetName(mc.cdc),
	)...)

	return namesvcTxCmd
}
