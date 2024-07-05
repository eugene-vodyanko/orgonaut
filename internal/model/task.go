package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Task is used for the internal representation of the replication work unit
type Task struct {
	GroupId   string
	PartId    int
	BatchSize int
	Topic     string
	Query
}

type Query struct {
	From     string
	Columns  string
	PkColumn string
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Topic, validation.Required),
		validation.Field(&t.GroupId, validation.Required),
		validation.Field(&t.BatchSize, validation.Required),
		validation.Field(&t.Query),
	)
}

func (q *Query) Validate() error {
	return validation.ValidateStruct(
		q,
		validation.Field(&q.From, validation.Required),
		validation.Field(&q.Columns, validation.Required),
		validation.Field(&q.PkColumn, validation.Required),
	)
}
