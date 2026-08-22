//go:build test

package godogutil

import (
	"maps"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/require"
)

type SharedData struct {
	ReturnedError error
	DataMap       map[string]map[string]string
	Require       *require.Assertions
}

func InitBaseData(t *testing.T) *SharedData {
	t.Helper()

	return &SharedData{
		DataMap: make(map[string]map[string]string),
		Require: require.New(t),
	}
}

func (s *SharedData) SyncTableToMap(table *godog.Table) []map[string]string {
	tableMap := s.getMapFromTable(table)

	for i, row := range tableMap {
		key := row["Key"]
		if s.DataMap[key] == nil {
			s.DataMap[key] = row
		} else {
			maps.Copy(s.DataMap[key], row)
		}

		tableMap[i] = s.DataMap[key]
	}

	return tableMap
}

func (s *SharedData) getMapFromTable(table *godog.Table) []map[string]string {
	var (
		res    = make([]map[string]string, max(1, len(table.Rows)-1))
		keyMap = make(map[string]struct{}, len(res))
	)

	s.Require.Greater(len(table.Rows), 1, "Table must have at least one row of data")

	headerRow := table.Rows[0]

	for i := 1; i < len(table.Rows); i++ {
		rowMap := make(map[string]string)

		for i, cell := range table.Rows[i].Cells {
			header := headerRow.Cells[i].Value
			rowMap[header] = cell.Value
		}

		rowKey, ok := rowMap["Key"]
		s.Require.True(ok, "Key column is missing in the table")

		_, exists := keyMap[rowKey]
		s.Require.False(exists, "Duplicate key found in the table: %s", rowKey)

		res[i-1] = rowMap
		keyMap[rowKey] = struct{}{}
	}

	return res
}
