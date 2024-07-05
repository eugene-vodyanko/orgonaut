package model

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Action string

// Command types
const (
	CREATE Action = "c" // insert
	UPDATE Action = "u" // update
	DELETE Action = "d" // delete
)

// Record describes the internal representation of the modified row from the database
// plus some additional information (such as the type of operation or the name/value of the primary key).
// The row attributes are represented as a map of string values, where the key is the column name.
type Record struct {
	Meta
	Fields map[string]string
}

// Meta information contains auxiliary fields.
// Such fields in xml format have names starting with a double underscore character.
type Meta struct {
	Pk
	Op   Action `xml:"__op" json:"__op"`
	Ts   string `xml:"__ts" json:"__ts"`
	UxTs string `xml:"__ux_ts" json:"__ux_ts"`
}

// Pk represents a primary single-part key with a string representation of the value.
type Pk struct {
	Name  string `xml:"__pk_name"`
	Value string `xml:"__pk_val"`
}

func (r *Record) GetKey() ([]byte, error) {
	err := r.Pk.Validate()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%s=%s", r.Pk.Name, r.Pk.Value)), nil
}

func (r *Record) GetValue() ([]byte, error) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}

	return json.Marshal(r.Fields)
}

func (p *Pk) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Value, validation.Required),
	)
}

func (m *Meta) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.Op, validation.Required, validation.In(UPDATE, DELETE, CREATE)),
	)
}
