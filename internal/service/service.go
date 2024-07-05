package service

import (
	"context"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
)

// RelayService is the main engine of the application
type RelayService struct {
	source Repository
	dest   Broker
	tx     Transactor
}

func New(sourceBroker Repository, tx Transactor, destBroker Broker) *RelayService {
	return &RelayService{sourceBroker, destBroker, tx}
}

// Relay requests the next entries in the database and sends them to the broker
// and returns the number of processed records.
//
// Processing implies transactional semantics: records are marked processed (transaction is committed)
// if the function completes without errors. Otherwise, the transaction is rolled back.
// Processing will not progress until the cause of the error is resolved.
// Delivery guarantees can be understood as at least once.
func (s *RelayService) Relay(ctx context.Context, task *model.Task) (uint16, error) {

	var amount int
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		items, err := s.source.GetRecords(txCtx, task)
		if err != nil {
			return fmt.Errorf("service - get records error: %w", err)
		}

		amount = len(items)

		if amount > 0 {
			err = s.dest.SendRecords(ctx, task.Topic, items)
			if err != nil {
				return fmt.Errorf("service - send records: %w", err)
			}
		}

		return nil
	})

	return uint16(amount), err
}
