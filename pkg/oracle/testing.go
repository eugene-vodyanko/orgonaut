package oracle

import "testing"

func TestStore(t *testing.T) (*Oracle, func()) {
	t.Helper()

	url := "localhost:1521/orcl"
	schema := "orgon"
	username := "orgon"
	password := "orgon"

	o, err := New(username, password, url, schema, 10, 2, 500, 500)
	if err != nil {
		t.Fatal(err)
	}

	return o, func() {
		_ = o.Close()
	}
}
