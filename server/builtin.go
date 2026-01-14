package server

import (
	"github.com/bitvault/chain"
	"github.com/bitvault/consensus"
	consensusbvbft "github.com/bitvault/consensus/bvbft"
	consensusDev "github.com/bitvault/consensus/dev"
	consensusDummy "github.com/bitvault/consensus/dummy"
	consensusIBFT "github.com/bitvault/consensus/ibft"
	"github.com/bitvault/forkmanager"
	"github.com/bitvault/secrets"
	"github.com/bitvault/secrets/awsssm"
	"github.com/bitvault/secrets/gcpssm"
	"github.com/bitvault/secrets/hashicorpvault"
	"github.com/bitvault/secrets/local"
	"github.com/bitvault/state"
)

type GenesisFactoryHook func(config *chain.Chain, engineName string) func(*state.Transition) error

type ConsensusType string

type ForkManagerFactory func(forks *chain.Forks) error

type ForkManagerInitialParamsFactory func(config *chain.Chain) (*forkmanager.ForkParams, error)

const (
	DevConsensus   ConsensusType = "dev"
	IBFTConsensus  ConsensusType = "ibft"
	bvbftConsensus ConsensusType = consensusbvbft.ConsensusName
	DummyConsensus ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:   consensusDev.Factory,
	IBFTConsensus:  consensusIBFT.Factory,
	bvbftConsensus: consensusbvbft.Factory,
	DummyConsensus: consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

var genesisCreationFactory = map[ConsensusType]GenesisFactoryHook{
	bvbftConsensus: consensusbvbft.GenesisPostHookFactory,
}

var forkManagerFactory = map[ConsensusType]ForkManagerFactory{
	bvbftConsensus: consensusbvbft.ForkManagerFactory,
}

var forkManagerInitialParamsFactory = map[ConsensusType]ForkManagerInitialParamsFactory{
	bvbftConsensus: consensusbvbft.ForkManagerInitialParamsFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
