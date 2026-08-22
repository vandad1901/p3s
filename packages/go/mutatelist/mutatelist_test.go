package mutatelist_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/vandad1901/p3s/packages/go/mutatelist"
)

type testItem struct {
	ID    int64
	Value string
}

func (t testItem) GetID() int64 {
	return t.ID
}

type testStructure struct {
	name        string
	original    []testItem
	mutations   mutatelist.MutateList[testItem, int64]
	expectedRet []testItem
	expectedErr error
}

func TestMergeItems(t *testing.T) {
	t.Parallel()

	tests := []testStructure{
		testEmptyOriginalAndMutations(),
		testOriginalWithNoMutations(),
		testOriginalWithAllMutations(),
		testInvalidDeletedID(),
		testInvalidUpdatedID(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := mutatelist.MergeItems(tt.original, &tt.mutations)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got=%v want=%v", err, tt.expectedErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.expectedRet, got) {
				t.Fatalf("result mismatch: want=%v got=%v", tt.expectedRet, got)
			}
		})
	}
}

func testEmptyOriginalAndMutations() testStructure {
	return testStructure{
		name:     "empty original and empty mutations",
		original: []testItem{},
		mutations: mutatelist.MutateList[testItem, int64]{
			Added:   []testItem{},
			Edited:  []testItem{},
			Removed: []int64{},
		},
		expectedRet: []testItem{},
	}
}

func testOriginalWithNoMutations() testStructure {
	return testStructure{
		name: "some items in original and empty mutations",
		original: []testItem{
			{ID: 1, Value: "a"},
			{ID: 2, Value: "b"},
			{ID: 3, Value: "b"},
		},
		mutations: mutatelist.MutateList[testItem, int64]{
			Added:   []testItem{},
			Edited:  []testItem{},
			Removed: []int64{},
		},
		expectedRet: []testItem{
			{ID: 1, Value: "a"},
			{ID: 2, Value: "b"},
			{ID: 3, Value: "b"},
		},
	}
}

func testOriginalWithAllMutations() testStructure {
	return testStructure{
		name: "some items in original and all mutations",
		original: []testItem{
			{ID: 1, Value: "a"},
			{ID: 2, Value: "b"},
			{ID: 3, Value: "c"},
		},
		mutations: mutatelist.MutateList[testItem, int64]{
			Added:   []testItem{{ID: 4, Value: "d"}},
			Edited:  []testItem{{ID: 2, Value: "b-updated"}},
			Removed: []int64{3},
		},
		expectedRet: []testItem{
			{ID: 1, Value: "a"},
			{ID: 2, Value: "b-updated"},
			{ID: 4, Value: "d"},
		},
	}
}

func testInvalidDeletedID() testStructure {
	return testStructure{
		name: "returns InvalidIDErr when deleted contains unknown id",
		original: []testItem{
			{ID: 1, Value: "a"},
			{ID: 2, Value: "b"},
		},
		mutations: mutatelist.MutateList[testItem, int64]{
			Removed: []int64{999},
		},
		expectedErr: mutatelist.ErrInvalidID,
	}
}

func testInvalidUpdatedID() testStructure {
	return testStructure{
		name: "returns InvalidIDErr when updated contains unknown id",
		original: []testItem{
			{ID: 1, Value: "a"},
			{ID: 2, Value: "b"},
		},
		mutations: mutatelist.MutateList[testItem, int64]{
			Edited: []testItem{{ID: 999, Value: "unknown"}},
		},
		expectedErr: mutatelist.ErrInvalidID,
	}
}
