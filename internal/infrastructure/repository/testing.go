package repository

import (
	"github.com/eugene-vodyanko/orgonaut/pkg/oracle"
	"testing"
)

func TestOra(t *testing.T) (*oracle.Oracle, string, func()) {
	t.Helper()

	dbSchema := "orgon"
	ora, teardownOra := oracle.TestStore(t)

	return ora, dbSchema,
		func() {
			teardownOra()
		}
}
