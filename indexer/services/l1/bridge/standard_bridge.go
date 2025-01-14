package bridge

import (
	"context"

	"github.com/ethereum-optimism/optimism/indexer/db"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type StandardBridge struct {
	name     string
	ctx      context.Context
	address  common.Address
	client   bind.ContractFilterer
	filterer *bindings.L1StandardBridgeFilterer
}

func (s *StandardBridge) Address() common.Address {
	return s.address
}

func (s *StandardBridge) GetDepositsByBlockRange(start, end uint64) (DepositsMap, error) {
	depositsByBlockhash := make(DepositsMap)

	iter, err := FilterERC20DepositInitiatedWithRetry(s.ctx, s.filterer, &bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		logger.Error("Error fetching filter", "err", err)
	}

	for iter.Next() {
		depositsByBlockhash[iter.Event.Raw.BlockHash] = append(
			depositsByBlockhash[iter.Event.Raw.BlockHash], db.Deposit{
				TxHash:      iter.Event.Raw.TxHash,
				L1Token:     iter.Event.L1Token,
				L2Token:     iter.Event.L2Token,
				FromAddress: iter.Event.From,
				ToAddress:   iter.Event.To,
				Amount:      iter.Event.Amount,
				Data:        iter.Event.ExtraData,
				LogIndex:    iter.Event.Raw.Index,
			})
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}

	return depositsByBlockhash, nil
}

func (s *StandardBridge) GetWithdrawalsByBlockRange(start, end uint64) (WithdrawalsMap, error) {
	withdrawalsByBlockHash := make(WithdrawalsMap)

	iter, err := FilterERC20WithdrawalFinalizedWithRetry(s.ctx, s.filterer, &bind.FilterOpts{
		Start: start,
		End:   &end,
	})

	if err != nil {
		logger.Error("Error fetching filter", "err", err)
	}

	for iter.Next() {
		withdrawalsByBlockHash[iter.Event.Raw.BlockHash] = append(
			withdrawalsByBlockHash[iter.Event.Raw.BlockHash], db.Withdrawal{
				TxHash:      iter.Event.Raw.TxHash,
				L1Token:     iter.Event.L1Token,
				L2Token:     iter.Event.L2Token,
				FromAddress: iter.Event.From,
				ToAddress:   iter.Event.To,
				Amount:      iter.Event.Amount,
				Data:        iter.Event.ExtraData,
				LogIndex:    iter.Event.Raw.Index,
			})
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}

	return withdrawalsByBlockHash, nil
}

func (s *StandardBridge) String() string {
	return s.name
}
