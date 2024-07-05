package service

import (
	"context"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
)

type (
	Relayer interface {
		Relay(context.Context, *model.Task) (uint16, error)
	}

	Repository interface {
		GetRecords(context.Context, *model.Task) ([]*model.Record, error)
	}

	Broker interface {
		SendRecords(context.Context, string, []*model.Record) error
	}
)
