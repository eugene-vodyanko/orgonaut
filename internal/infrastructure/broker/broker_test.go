package broker

import (
	"context"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
	"github.com/eugene-vodyanko/orgonaut/pkg/kafka/kafkakit"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestBroker_SendRecords(t *testing.T) {
	start := time.Now()

	writer := NewBroker(
		kafkakit.TestWriter(t, ""),
	)

	var records = []*model.Record{
		{
			Meta: model.Meta{
				Pk: model.Pk{
					Name:  "id",
					Value: "42",
				},
				Op: "u",
			},
			Fields: map[string]string{"col_name": "col_value"},
		},
	}
	err := writer.SendRecords(context.Background(), "test_tab", records)

	assert.NoError(t, err)

	elapsed := time.Now()

	t.Logf("elapsed: %v", elapsed.Sub(start))
}
