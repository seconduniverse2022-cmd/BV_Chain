## Deposit

Bridge ERC 20 tokens from rootchain to childchain via deposit.

```mermaid
sequenceDiagram
	User->>bitvault: deposit
	bitvault->>RootERC20.sol: approve(RootERC20Predicate)
	bitvault->>RootERC20Predicate.sol: deposit()
	RootERC20Predicate.sol->>RootERC20Predicate.sol: mapToken()
	RootERC20Predicate.sol->>StateSender.sol: syncState(MAP_TOKEN_SIG), recv=ChildERC20Predicate
	RootERC20Predicate.sol-->>bitvault: TokenMapped Event
	StateSender.sol-->>bitvault: StateSynced Event to map tokens on child predicate
	RootERC20Predicate.sol->>StateSender.sol: syncState(DEPOSIT_SIG), recv=ChildERC20Predicate
	StateSender.sol-->>bitvault: StateSynced Event to deposit on child chain
	bitvault->>User: ok
	bitvault->>StateReceiver.sol:commit()
	StateReceiver.sol-->>bitvault: NewCommitment Event
	bitvault->>StateReceiver.sol:execute()
	StateReceiver.sol->>ChildERC20Predicate.sol:onStateReceive()
	ChildERC20Predicate.sol->>ChildERC20.sol: mint()
	StateReceiver.sol-->>bitvault:StateSyncResult Event
```

## Withdraw

Bridge ERC 20 tokens from childchain to rootchain via withdrawal.

```mermaid
sequenceDiagram
	User->>bitvault: withdraw
	bitvault->>ChildERC20Predicate.sol: withdrawTo()
	ChildERC20Predicate.sol->>ChildERC20: burn()
	ChildERC20Predicate.sol->>L2StateSender.sol: syncState(WITHDRAW_SIG), recv=RootERC20Predicate
	bitvault->>User: tx hash
	User->>bitvault: get tx receipt
	bitvault->>User: exit event id
	ChildERC20Predicate.sol-->>bitvault: L2ERC20Withdraw Event
	L2StateSender.sol-->>bitvault: StateSynced Event
	bitvault->>bitvault: Seal block
	bitvault->>CheckpointManager.sol: submit()
```
## Exit

Finalize withdrawal of ERC 20 tokens from childchain to rootchain.

```mermaid
sequenceDiagram
	User->>bitvault: exit, event id:X
	bitvault->>bitvault: bridge_generateExitProof()
	bitvault->>CheckpointManager.sol: getCheckpointBlock()
	CheckpointManager.sol->>bitvault: blockNum
	bitvault->>bitvault: getExitEventsForProof(epochNum, blockNum)
	bitvault->>bitvault: createExitTree(exitEvents)
	bitvault->>bitvault: generateProof()
	bitvault->>ExitHelper.sol: exit()
	ExitHelper.sol->>CheckpointManager.sol: getEventMembershipByBlockNumber()
	ExitHelper.sol->>RootERC20Predicate.sol:onL2StateReceive()
	RootERC20Predicate.sol->>RootERC20: transfer()
	bitvault->>User: ok
	RootERC20Predicate.sol-->>bitvault: ERC20Withdraw Event
	ExitHelper.sol-->>bitvault: ExitProcessed Event
```

