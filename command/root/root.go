package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bitvault/command/backup"
	"github.com/bitvault/command/bridge"
	"github.com/bitvault/command/bvbft"
	"github.com/bitvault/command/bvbftsecrets"
	"github.com/bitvault/command/genesis"
	"github.com/bitvault/command/helper"
	"github.com/bitvault/command/ibft"
	"github.com/bitvault/command/license"
	"github.com/bitvault/command/monitor"
	"github.com/bitvault/command/peers"
	"github.com/bitvault/command/regenesis"
	"github.com/bitvault/command/rootchain"
	"github.com/bitvault/command/secrets"
	"github.com/bitvault/command/server"
	"github.com/bitvault/command/status"
	"github.com/bitvault/command/txpool"
	"github.com/bitvault/command/version"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "bitvault is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		rootchain.GetCommand(),
		monitor.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		license.GetCommand(),
		bvbftsecrets.GetCommand(),
		bvbft.GetCommand(),
		bridge.GetCommand(),
		regenesis.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
