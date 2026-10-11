//go:build test

package godogutil

import (
	"context"
	"strings"

	"github.com/cucumber/godog"
	"github.com/vandad1901/p3s/packages/go/usercontext"
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

func LoginAsUserStepFunc(s *SharedData) func(ctx context.Context, table *godog.Table) context.Context {
	return func(ctx context.Context, table *godog.Table) context.Context {
		currentMap := s.SyncTableToMap(table)
		s.Require.Len(currentMap, 1)

		for _, row := range currentMap {
			return usercontext.CtxWithUser(ctx, ResolveInt64(s, row, "UserID"))
		}

		return ctx
	}
}
