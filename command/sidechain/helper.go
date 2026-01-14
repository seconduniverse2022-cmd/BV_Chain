package sidechain

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitvault/command/bvbftsecrets"
	rootHelper "github.com/bitvault/command/rootchain/helper"
	"github.com/bitvault/consensus/bvbft"
	"github.com/bitvault/consensus/bvbft/contractsapi"
	"github.com/bitvault/consensus/bvbft/wallet"
	"github.com/bitvault/contracts"
	"github.com/bitvault/helper/common"
	"github.com/bitvault/txrelayer"
	"github.com/bitvault/types"
	"github.com/umbracle/ethgo"
)

const (
	AmountFlag = "amount"
)

func CheckIfDirectoryExist(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("provided directory '%s' doesn't exist", dir)
	}

	return nil
}

func ValidateSecretFlags(dataDir, config string) error {
	if config == "" {
		if dataDir == "" {
			return bvbftsecrets.ErrInvalidParams
		} else {
			return CheckIfDirectoryExist(dataDir)
		}
	}

	return nil
}

// GetAccount resolves secrets manager and returns an account object
func GetAccount(accountDir, accountConfig string) (*wallet.Account, error) {
	// resolve secrets manager instance and allow usage of insecure local secrets manager
	secretsManager, err := bvbftsecrets.GetSecretsManager(accountDir, accountConfig, true)
	if err != nil {
		return nil, err
	}

	return wallet.NewAccountFromSecret(secretsManager)
}

// GetAccountFromDir returns an account object from local secrets manager
func GetAccountFromDir(accountDir string) (*wallet.Account, error) {
	return GetAccount(accountDir, "")
}

// GetValidatorInfo queries CustomSupernetManager, StakeManager and RewardPool smart contracts
// to retrieve validator info for given address
func GetValidatorInfo(validatorAddr ethgo.Address, supernetManager, stakeManager types.Address,
	chainID int64, rootRelayer, childRelayer txrelayer.TxRelayer) (*bvbft.ValidatorInfo, error) {
	validatorInfo, err := rootHelper.GetValidatorInfo(validatorAddr, supernetManager, stakeManager,
		chainID, rootRelayer)
	if err != nil {
		return nil, err
	}

	withdrawableFn := contractsapi.RewardPool.Abi.GetMethod("pendingRewards")

	encode, err := withdrawableFn.Encode([]interface{}{validatorAddr})
	if err != nil {
		return nil, err
	}

	response, err := childRelayer.Call(ethgo.ZeroAddress, ethgo.Address(contracts.RewardPoolContract), encode)
	if err != nil {
		return nil, err
	}

	withdrawableRewards, err := common.ParseUint256orHex(&response)
	if err != nil {
		return nil, err
	}

	validatorInfo.WithdrawableRewards = withdrawableRewards

	return validatorInfo, nil
}
