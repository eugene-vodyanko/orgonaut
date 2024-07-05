package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var rows =
// language=xml
`<?xml version="1.0"?>
<ROWSET>
 <ROW>
  <OP>u</OP>
  <_pkv>2</_pkv>
  <_dt>2024-06-10 14:45:56</_dt>
  <_ls>2024-06-10 14:45:56.948653</_ls>
  <_ts>2024-06-10T14:45:56.948651 +07:00</_ts>
  <_fl>3.14</_fl>
  <_pkn>2</_pkn>
  <__TS>1718009156929</__TS>
  <ID>2</ID>
  <DT>2024-04-14 22:44:37</DT>
  <STR>str:2</STR>
 </ROW>
 <ROW>
  <OP>u</OP>
  <_pkv>44</_pkv>
  <_dt>2024-06-10 14:45:56</_dt>
  <_ls>2024-06-10 14:45:56.948653</_ls>
  <_ts>2024-06-10T14:45:56.948651 +07:00</_ts>
  <_fl>3.14</_fl>
  <_pkn>44</_pkn>
  <__TS>1718009156929</__TS>
  <ID>44</ID>
  <DT>2024-03-03 22:44:37</DT>
  <STR>str:44</STR>
 </ROW>
</ROWSET>
`

func TestDecoder_decodeRecords(t *testing.T) {
	start := time.Now()

	_, err := decodeRecords(strings.NewReader(rows))

	elapsed := time.Now()

	assert.NoError(t, err)

	fmt.Println("elapsed:", elapsed.Sub(start))
}
