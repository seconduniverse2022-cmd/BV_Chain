package bvbft

import (
	"github.com/bitvault/command/rootchain/registration"
	"github.com/bitvault/command/rootchain/staking"
	"github.com/bitvault/command/rootchain/supernet"
	"github.com/bitvault/command/rootchain/supernet/stakemanager"
	"github.com/bitvault/command/rootchain/validators"
	"github.com/bitvault/command/rootchain/whitelist"
	"github.com/bitvault/command/rootchain/withdraw"
	"github.com/bitvault/command/sidechain/rewards"
	"github.com/bitvault/command/sidechain/unstaking"
	sidechainWithdraw "github.com/bitvault/command/sidechain/withdraw"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	bvbftCmd := &cobra.Command{
		Use:   "bvbft",
		Short: "bvbft command",
	}

	bvbftCmd.AddCommand(
		// sidechain (validator set) command to unstake on child chain
		unstaking.GetCommand(),
		// sidechain (validator set) command to withdraw stake on child chain
		sidechainWithdraw.GetCommand(),
		// sidechain (reward pool) command to withdraw pending rewards
		rewards.GetCommand(),
		// rootchain (stake manager) command to withdraw stake
		withdraw.GetCommand(),
		// rootchain (supernet manager) command that queries validator info
		validators.GetCommand(),
		// rootchain (supernet manager) whitelist validator
		whitelist.GetCommand(),
		// rootchain (supernet manager) register validator
		registration.GetCommand(),
		// rootchain (stake manager) stake command
		staking.GetCommand(),
		// rootchain (supernet manager) command for finalizing genesis
		// validator set and enabling staking
		supernet.GetCommand(),
		// rootchain command for deploying stake manager
		stakemanager.GetCommand(),
	)

	return bvbftCmd
}
