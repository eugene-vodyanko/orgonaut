package oracle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOracle_Ping(t *testing.T) {
	db, teardown := TestStore(t)
	defer teardown()

	err := db.Ping()

	assert.NoError(t, err)
}

func TestOracle_Close(t *testing.T) {
	_, teardown := TestStore(t)
	defer teardown()

	assert.NoError(t, nil)
}
