package withdraw

import (
	"fmt"
	"time"

	"github.com/bitvault/command"
	"github.com/bitvault/command/bvbftsecrets"
	"github.com/bitvault/command/helper"
	rootHelper "github.com/bitvault/command/rootchain/helper"
	sidechainHelper "github.com/bitvault/command/sidechain"
	"github.com/bitvault/consensus/bvbft/contractsapi"
	"github.com/bitvault/txrelayer"
	"github.com/bitvault/types"
	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
)

var params withdrawParams

func GetCommand() *cobra.Command {
	withdrawCmd := &cobra.Command{
		Use:     "withdraw-root",
		Short:   "Withdraws sender's withdrawable amount to specified address on the root chain",
		PreRunE: runPreRun,
		RunE:    runCommand,
	}

	setFlags(withdrawCmd)

	return withdrawCmd
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
		&params.addressTo,
		addressToFlag,
		"",
		"address where to withdraw withdrawable amount",
	)

	cmd.Flags().StringVar(
		&params.stakeManagerAddr,
		rootHelper.StakeManagerFlag,
		"",
		rootHelper.StakeManagerFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.amount,
		sidechainHelper.AmountFlag,
		"",
		"amount to withdraw",
	)

	cmd.MarkFlagsMutuallyExclusive(bvbftsecrets.AccountDirFlag, bvbftsecrets.AccountConfigFlag)
	helper.RegisterJSONRPCFlag(cmd)
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

	txRelayer, err := txrelayer.NewTxRelayer(txrelayer.WithIPAddress(params.jsonRPC),
		txrelayer.WithReceiptTimeout(150*time.Millisecond))
	if err != nil {
		return err
	}

	withdrawFn := &contractsapi.WithdrawStakeStakeManagerFn{
		To:     types.StringToAddress(params.addressTo),
		Amount: params.amountValue,
	}

	encoded, err := withdrawFn.EncodeAbi()
	if err != nil {
		return err
	}

	stakeManagerAddr := ethgo.Address(types.StringToAddress(params.stakeManagerAddr))
	txn := rootHelper.CreateTransaction(validatorAccount.Ecdsa.Address(), &stakeManagerAddr, encoded, nil, true)

	receipt, err := txRelayer.SendTransaction(txn, validatorAccount.Ecdsa)
	if err != nil {
		return err
	}

	if receipt.Status == uint64(types.ReceiptFailed) {
		return fmt.Errorf("withdraw transaction failed on block %d", receipt.BlockNumber)
	}

	result := &withdrawResult{
		ValidatorAddress: validatorAccount.Ecdsa.Address().String(),
	}

	var (
		withdrawalEvent contractsapi.StakeWithdrawnEvent
		foundLog        bool
	)

	for _, log := range receipt.Logs {
		doesMatch, err := withdrawalEvent.ParseLog(log)
		if !doesMatch {
			continue
		}

		if err != nil {
			return err
		}

		result.Amount = withdrawalEvent.Amount.Uint64()
		result.WithdrawnTo = withdrawalEvent.Recipient.String()
		foundLog = true

		break
	}

	if !foundLog {
		return fmt.Errorf("could not find an appropriate log in receipt that withdrawal happened")
	}

	outputter.WriteCommandResult(result)

	return nil
}
