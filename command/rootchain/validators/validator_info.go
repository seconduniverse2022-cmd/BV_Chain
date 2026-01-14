package validators

import (
	"fmt"

	"github.com/bitvault/command"
	"github.com/bitvault/command/bvbftsecrets"
	"github.com/bitvault/command/helper"
	rootHelper "github.com/bitvault/command/rootchain/helper"
	sidechainHelper "github.com/bitvault/command/sidechain"
	"github.com/bitvault/txrelayer"
	"github.com/bitvault/types"
	"github.com/spf13/cobra"
)

var (
	params validatorInfoParams
)

func GetCommand() *cobra.Command {
	validatorInfoCmd := &cobra.Command{
		Use:     "validator-info",
		Short:   "Gets validator info",
		PreRunE: runPreRun,
		RunE:    runCommand,
	}

	helper.RegisterJSONRPCFlag(validatorInfoCmd)
	setFlags(validatorInfoCmd)

	return validatorInfoCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.accountDir,
		bvbftsecrets.AccountDirFlag,
		"",
		bvbftsecrets.AccountDirFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.accountConfig,
		bvbftsecrets.AccountConfigFlag,
		"",
		bvbftsecrets.AccountConfigFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.supernetManagerAddress,
		rootHelper.SupernetManagerFlag,
		"",
		rootHelper.SupernetManagerFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.stakeManagerAddress,
		rootHelper.StakeManagerFlag,
		"",
		rootHelper.StakeManagerFlagDesc,
	)

	cmd.Flags().Int64Var(
		&params.chainID,
		bvbftsecrets.ChainIDFlag,
		0,
		bvbftsecrets.ChainIDFlagDesc,
	)

	cmd.MarkFlagsMutuallyExclusive(bvbftsecrets.AccountDirFlag, bvbftsecrets.AccountConfigFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	params.jsonRPC = helper.GetJSONRPCAddress(cmd)

	return params.validateFlags()
}

func runCommand(cmd *cobra.Command, _ []string) error {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	validatorAccount, err := sidechainHelper.GetAccount(params.accountDir, params.accountConfig)
	if err != nil {
		return err
	}

	txRelayer, err := txrelayer.NewTxRelayer(txrelayer.WithIPAddress(params.jsonRPC))
	if err != nil {
		return err
	}

	validatorAddr := validatorAccount.Ecdsa.Address()
	supernetManagerAddr := types.StringToAddress(params.supernetManagerAddress)
	stakeManagerAddr := types.StringToAddress(params.stakeManagerAddress)

	validatorInfo, err := rootHelper.GetValidatorInfo(validatorAddr,
		supernetManagerAddr, stakeManagerAddr, params.chainID, txRelayer)
	if err != nil {
		return fmt.Errorf("failed to get validator info for %s: %w", validatorAddr, err)
	}

	outputter.WriteCommandResult(&validatorsInfoResult{
		Address:     validatorInfo.Address.String(),
		Stake:       validatorInfo.Stake.Uint64(),
		Active:      validatorInfo.IsActive,
		Whitelisted: validatorInfo.IsWhitelisted,
	})

	return nil
}
