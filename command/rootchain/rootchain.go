package rootchain

import (
	"github.com/spf13/cobra"

	"github.com/bitvault/command/rootchain/deploy"
	"github.com/bitvault/command/rootchain/fund"
	"github.com/bitvault/command/rootchain/premine"
	"github.com/bitvault/command/rootchain/server"
)

// GetCommand creates "rootchain" helper command
func GetCommand() *cobra.Command {
	rootchainCmd := &cobra.Command{
		Use:   "rootchain",
		Short: "Top level rootchain helper command.",
	}

	rootchainCmd.AddCommand(
		// rootchain server
		server.GetCommand(),
		// rootchain deploy
		deploy.GetCommand(),
		// rootchain fund
		fund.GetCommand(),
		// rootchain premine
		premine.GetCommand(),
	)

	return rootchainCmd
}
