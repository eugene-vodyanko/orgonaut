package kafkakit

import (
	"testing"
)

func TestWriter(t *testing.T, topic string) *Writer {
	t.Helper()

	writer, err := NewWriter([]string{"localhost:9092"}, topic, false,
		50, 10, "one", true, 1048576)
	if err != nil {
		t.Fatal(err)
	}

	return writer
}
