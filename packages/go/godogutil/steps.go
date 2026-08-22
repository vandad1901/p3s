//go:build test

package godogutil

import (
	"strings"

	"github.com/cucumber/godog"
)

func SyncToDataMapStep(s *SharedData) func(table *godog.Table) {
	return func(table *godog.Table) {
		s.SyncTableToMap(table)
	}
}

func AddToHeaderStep(s *SharedData, itemsKey string) func(*godog.Table) {
	return func(table *godog.Table) {
		currentMap := s.SyncTableToMap(table)
		for _, row := range currentMap {
			headerValues := s.DataMap[row["HeaderKey"]][itemsKey]
			headerValuesSlice := commaSeparatedValueToSlice(headerValues)
			headerValuesSlice = append(headerValuesSlice, row["Key"])
			s.DataMap[row["HeaderKey"]][itemsKey] = strings.Join(headerValuesSlice, ",")
		}
	}
}

func AssertErrorStep(s *SharedData) func(errorKey string) {
	return func(errorKey string) {
		s.Require.Error(s.ReturnedError)
		s.Require.EqualError(s.ReturnedError, errorKey)
	}
}
