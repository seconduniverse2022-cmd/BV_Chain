Before proceeding, ensure that your system meets the necessary [system requirements](operate/system.md).

To access the pre-built releases, visit the [GitHub releases page](https://github.com/bitvault/releases). The client provides cross-compiled AMD64/ARM64 binaries for Darwin and Linux.

To locally run the bitvault test environment, run the following command from the project's root:

  ```bash
  ./scripts/cluster bvbft
  ```

**That's it! You should have successfully been able to start a local bitvault-powered chain just by running the script.**

> - Stop the network: "CTRL/Command C" or `./scripts/cluster bvbft stop`.
> - Destroy the network: `./scripts/cluster bvbft destroy`.


## Deployment script details

The script is available under the "scripts" directory of the client.
These are the optional configuration parameters you can pass to the script:

### Flags

| Flag | Description | Default Value |
|------|-------------|---------------|
| --block-gas-limit | Maximum gas allowed for a block. | 10000000 |
| --premine | Address and amount of tokens to premine in the genesis block. | 0x85da99c8a7c2c95964c8efd687e95e632fc533d6:1000000000000000000000 |
| --epoch-size | Number of blocks per epoch. | 10 |
| --data-dir | Directory to store chain data. | test-chain- |
| --num | Number of nodes in the network. | 4 |
| --bootnode | Bootstrap node address in multiaddress format. | /ip4/127.0.0.1/tcp/30301/p2p/... |
| --insecure | Disable TLS. | |
| --log-level | Logging level for validators. | INFO |
| --seal | Enable block sealing. | |
| --help | Print usage information. | |

After running the command, the test network will be initialized with bvbft consensus engine and the genesis file will be created. Then, the four validators will start running, and their log outputs will be displayed in the terminal.

By default, this will start an bitvault network with bvbft consensus engine, four validators, and premine of 1 billion tokens at address `0x85da99c8a7c2c95964c8efd687e95e632fc533d6`.

The nodes will continue to run until stopped manually. To stop the network, open a new session and use the following command, or, simply press "CTRL/Command C" in the CLI:

  ```bash
  ./scripts/cluster bvbft stop
  ```

If you want to destroy the environment, use the following command:

  ```bash
  ./scripts/cluster bvbft destroy
  ```

### Explanation of the deployment script

The deployment script is a wrapper script for starting an bitvault test network with bvbft consensus engine. 
It offers the following functionality:

- Initialize the network with either IBFT or bvbft consensus engine.
- Create the genesis file for the test network.
- Start the validators on four separate ports.
- Write the logs to separate log files for each validator.
- Stop and destroy the environment when no longer needed.
- The script also allows you to choose between running the environment from a local binary or a Docker container.

## Deployment script

<details>
<summary>Click to open</summary>

```sh
#!/usr/bin/env bash

dp_error_flag=0

# Check if jq is installed
if [[ "$1" == "bvbft" ]] && ! command -v jq >/dev/null 2>&1; then
  echo "jq is not installed."
  echo "Manual installation instructions: Visit https://jqlang.github.io/jq/ for more information."
  dp_error_flag=1
fi

# Check if curl is installed
if [[ "$1" == "bvbft" ]] && ! command -v curl >/dev/null 2>&1; then
  echo "curl is not installed."
  echo "Manual installation instructions: Visit https://everything.curl.dev/get/ for more information."
  dp_error_flag=1
fi

# Check if docker-compose is installed
if [[ "$2" == "--docker" ]] && ! command -v docker-compose >/dev/null 2>&1; then
  echo "docker-compose is not installed."
  echo "Manual installation instructions: Visit https://docs.docker.com/compose/install/ for more information."
  dp_error_flag=1
fi

# Stop script if any of the dependencies have failed
if [[ "$dp_error_flag" -eq 1 ]]; then
  echo "Missing dependencies. Please install them and run the script again."
  exit 1
fi

function showhelp(){
  echo "Usage: cluster {consensus} [{command}] [{flags}]"
  echo "Consensus:"
  echo "  ibft            Start Supernets test environment locally with ibft consensus"
  echo "  bvbft         Start Supernets test environment locally with bvbft consensus"
  echo "Commands:"
  echo "  stop            Stop the running environment"
  echo "  destroy         Destroy the running environment"
  echo "  write-logs      Writes STDOUT and STDERR output to log file. Not applicable when using --docker flag."
  echo "Flags:"
  echo "  --docker        Run using Docker (requires docker-compose)."
  echo "  --help          Display this help information"
  echo "Examples:"
  echo "  cluster bvbft -- Run the script with the bvbft consensus"
  echo "  cluster bvbft --docker -- Run the script with the bvbft consensus using docker"
  echo "  cluster bvbft stop -- Stop the running environment"
}

function initIbftConsensus() {
  echo "Running with ibft consensus"
  ./bitvault secrets init --insecure --data-dir test-chain- --num 4

  node1_id=$(./bitvault secrets output --data-dir test-chain-1 | grep Node | head -n 1 | awk -F ' ' '{print $4}')
  node2_id=$(./bitvault secrets output --data-dir test-chain-2 | grep Node | head -n 1 | awk -F ' ' '{print $4}')

  genesis_params="--consensus ibft --ibft-validators-prefix-path test-chain- \
    --bootnode /ip4/127.0.0.1/tcp/30301/p2p/$node1_id \
    --bootnode /ip4/127.0.0.1/tcp/30302/p2p/$node2_id"
}

function initbvbftConsensus() {
  echo "Running with bvbft consensus"
  genesis_params="--consensus bvbft"

  address1=$(./bitvault bvbft-secrets --insecure --data-dir test-chain-1 | grep Public | head -n 1 | awk -F ' ' '{print $5}')
  address2=$(./bitvault bvbft-secrets --insecure --data-dir test-chain-2 | grep Public | head -n 1 | awk -F ' ' '{print $5}')
  address3=$(./bitvault bvbft-secrets --insecure --data-dir test-chain-3 | grep Public | head -n 1 | awk -F ' ' '{print $5}')
  address4=$(./bitvault bvbft-secrets --insecure --data-dir test-chain-4 | grep Public | head -n 1 | awk -F ' ' '{print $5}')
}

function createGenesis() {
  ./bitvault genesis $genesis_params \
    --block-gas-limit 10000000 \
    --premine 0x85da99c8a7c2c95964c8efd687e95e632fc533d6:1000000000000000000000 \
    --premine 0x0000000000000000000000000000000000000000 \
    --epoch-size 10 \
    --reward-wallet 0xDEADBEEF:1000000 \
    --native-token-config "bitvault:MATIC:18:true:$address1" \
    --burn-contract 0:0x0000000000000000000000000000000000000000 \
    --proxy-contracts-admin 0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed
}

function initRootchain() {
  echo "Initializing rootchain"

  if [ "$1" == "write-logs" ]; then
    echo "Writing rootchain server logs to the file..."
    ./bitvault rootchain server 2>&1 | tee ./rootchain-server.log &
  else
    ./bitvault rootchain server >/dev/null &
  fi

  set +e
  while true; do
    if curl -sSf -o /dev/null http://127.0.0.1:8545; then
      break
    fi
    sleep 1
  done
  set -e

  proxyContractsAdmin=0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed

  ./bitvault bvbft stake-manager-deploy \
    --jsonrpc http://127.0.0.1:8545 \
    --proxy-contracts-admin ${proxyContractsAdmin} \
    --test

  stakeManagerAddr=$(cat genesis.json | jq -r '.params.engine.bvbft.bridge.stakeManagerAddr')
  stakeToken=$(cat genesis.json | jq -r '.params.engine.bvbft.bridge.stakeTokenAddr')

  ./bitvault rootchain deploy \
    --stake-manager ${stakeManagerAddr} \
    --stake-token ${stakeToken} \
    --proxy-contracts-admin ${proxyContractsAdmin} \
    --test

  customSupernetManagerAddr=$(cat genesis.json | jq -r '.params.engine.bvbft.bridge.customSupernetManagerAddr')
  supernetID=$(cat genesis.json | jq -r '.params.engine.bvbft.supernetID')

  ./bitvault rootchain fund \
    --stake-token ${stakeToken} \
    --mint \
    --addresses ${address1},${address2},${address3},${address4} \
    --amounts 1000000000000000000000000,1000000000000000000000000,1000000000000000000000000,1000000000000000000000000

  ./bitvault bvbft whitelist-validators \
    --addresses ${address1},${address2},${address3},${address4} \
    --supernet-manager ${customSupernetManagerAddr} \
    --private-key aa75e9a7d427efc732f8e4f1a5b7646adcc61fd5bae40f80d13c8419c9f43d6d \
    --jsonrpc http://127.0.0.1:8545

  counter=1
  while [ $counter -le 4 ]; do
    echo "Registering validator: ${counter}"

    ./bitvault bvbft register-validator \
      --supernet-manager ${customSupernetManagerAddr} \
      --data-dir test-chain-${counter} \
      --jsonrpc http://127.0.0.1:8545

    ./bitvault bvbft stake \
      --data-dir test-chain-${counter} \
      --amount 1000000000000000000000000 \
      --supernet-id ${supernetID} \
      --stake-manager ${stakeManagerAddr} \
      --stake-token ${stakeToken} \
      --jsonrpc http://127.0.0.1:8545

    ((counter++))
  done

  ./bitvault bvbft supernet \
    --private-key aa75e9a7d427efc732f8e4f1a5b7646adcc61fd5bae40f80d13c8419c9f43d6d \
    --supernet-manager ${customSupernetManagerAddr} \
    --stake-manager ${stakeManagerAddr} \
    --finalize-genesis-set \
    --enable-staking \
    --jsonrpc http://127.0.0.1:8545
}

function startServerFromBinary() {
  if [ "$1" == "write-logs" ]; then
    echo "Writing validators logs to the files..."
    ./bitvault server --data-dir ./test-chain-1 --chain genesis.json \
      --grpc-address :10000 --libp2p :30301 --jsonrpc :10002 --relayer \
      --num-block-confirmations 2 --seal --log-level DEBUG 2>&1 | tee ./validator-1.log &
    ./bitvault server --data-dir ./test-chain-2 --chain genesis.json \
      --grpc-address :20000 --libp2p :30302 --jsonrpc :20002 \
      --num-block-confirmations 2 --seal --log-level DEBUG 2>&1 | tee ./validator-2.log &
    ./bitvault server --data-dir ./test-chain-3 --chain genesis.json \
      --grpc-address :30000 --libp2p :30303 --jsonrpc :30002 \
      --num-block-confirmations 2 --seal --log-level DEBUG 2>&1 | tee ./validator-3.log &
    ./bitvault server --data-dir ./test-chain-4 --chain genesis.json \
      --grpc-address :40000 --libp2p :30304 --jsonrpc :40002 \
      --num-block-confirmations 2 --seal --log-level DEBUG 2>&1 | tee ./validator-4.log &
    wait
  else
    ./bitvault server --data-dir ./test-chain-1 --chain genesis.json \
      --grpc-address :10000 --libp2p :30301 --jsonrpc :10002 --relayer \
      --num-block-confirmations 2 --seal --log-level DEBUG &
    ./bitvault server --data-dir ./test-chain-2 --chain genesis.json \
      --grpc-address :20000 --libp2p :30302 --jsonrpc :20002 \
      --num-block-confirmations 2 --seal --log-level DEBUG &
    ./bitvault server --data-dir ./test-chain-3 --chain genesis.json \
      --grpc-address :30000 --libp2p :30303 --jsonrpc :30002 \
      --num-block-confirmations 2 --seal --log-level DEBUG &
    ./bitvault server --data-dir ./test-chain-4 --chain genesis.json \
      --grpc-address :40000 --libp2p :30304 --jsonrpc :40002 \
      --num-block-confirmations 2 --seal --log-level DEBUG &
    wait
  fi
}

function startServerFromDockerCompose() {
  if [ "$1" != "bvbft" ]; then
    export bitvault_CONSENSUS="$1"
  fi

  docker-compose -f ./docker/local/docker-compose.yml up -d --build
}

function destroyDockerEnvironment() {
  docker-compose -f ./docker/local/docker-compose.yml down -v
}

function stopDockerEnvironment() {
  docker-compose -f ./docker/local/docker-compose.yml stop
}

set -e

# Show help if help flag is entered or no arguments are provided
if [[ "$1" == "--help" ]] || [[ $# -eq 0 ]]; then
  showhelp
  exit 0
fi

# Reset test-dirs
rm -rf test-chain-*
rm -f genesis.json

# Build binary
go build -o bitvault .

# If --docker flag is set run docker environment otherwise run from binary
case "$2" in
"--docker")
  # cluster {consensus} --docker destroy
  if [ "$3" == "destroy" ]; then
    destroyDockerEnvironment
    echo "Docker $1 environment destroyed!"
    exit 0
  # cluster {consensus} --docker stop
  elif [ "$3" == "stop" ]; then
    stopDockerEnvironment
    echo "Docker $1 environment stopped!"
    exit 0
  fi

  # cluster {consensus} --docker
  echo "Running $1 docker environment..."
  startServerFromDockerCompose $1
  echo "Docker $1 environment deployed."
  exit 0
  ;;
# cluster {consensus}
*)
  echo "Running $1 environment from local binary..."
  # Initialize ibft or bvbft consensus
  if [ "$1" == "ibft" ]; then
    # Initialize ibft consensus
    initIbftConsensus
    # Create genesis file and start the server from binary
    createGenesis
    startServerFromBinary $2
    exit 0
  elif [ "$1" == "bvbft" ]; then
    # Initialize bvbft consensus
    initbvbftConsensus
    # Create genesis file and start the server from binary
    createGenesis
    initRootchain $2
    startServerFromBinary $2
    exit 0
  else
    echo "Unsupported consensus mode. Supported modes are: ibft and bvbft."
    showhelp
    exit 1
  fi
  ;;
esac
```

</details>
<br />
