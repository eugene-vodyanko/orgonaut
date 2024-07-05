package repository

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
	"io"
	"strings"
)

const rowElementName = "ROW"

func makeRecords(input []byte) ([]*model.Record, error) {

	if input != nil {
		gzipReader, err := getGZipReader(input)
		if err != nil {
			return nil, fmt.Errorf("decompress error: %w", err)
		}

		records, err := decodeRecords(gzipReader)
		if err != nil {
			return nil, fmt.Errorf("decode error %w", err)
		}

		return records, nil
	}

	return nil, nil
}

type row struct {
	model.Meta
	Fields []byte `xml:",innerxml"`
}

func decodeRecords(r io.Reader) ([]*model.Record, error) {
	var rows []*model.Record
	d := xml.NewDecoder(r)
	for {
		t, err := d.Token()
		if t == nil || err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("rows token error: %w", err)
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == rowElementName {
				var row row
				err = d.DecodeElement(&row, &se)
				if err != nil {
					return nil, err
				}

				m, err := parseFields(bytes.NewReader(row.Fields))
				if err != nil {
					return nil, fmt.Errorf("field token error: %w", err)
				}

				rows = append(rows, &model.Record{
					Meta:   row.Meta,
					Fields: m,
				})

			}
		}

	}
	return rows, nil
}

func parseFields(s io.Reader) (map[string]string, error) {
	r := make(map[string]string)
	d := xml.NewDecoder(s)
	for t, err := d.Token(); err == nil; t, err = d.Token() {

		if se, ok := t.(xml.StartElement); ok {
			name := se.Name.Local
			token, err := d.Token()
			if err != nil {
				return nil, err
			}

			if cdata, ok := token.(xml.CharData); ok {
				r[strings.ToLower(name)] = string(cdata)
			}
		}
	}
	return r, nil
}

func getGZipReader(input []byte) (io.Reader, error) {
	if input != nil {
		bytesReader := bytes.NewReader(input)
		if bytesReader != nil {
			gzipReader, err := gzip.NewReader(bytesReader)
			if err != nil {
				return nil, err

			}

			return gzipReader, nil
		}
	}

	return nil, nil
}
