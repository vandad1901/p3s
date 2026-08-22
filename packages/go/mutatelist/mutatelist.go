package mutatelist

import "errors"

var ErrInvalidID = errors.New("invalid ID found")

type IDType interface {
	~int64 | ~string
}

type MutateList[T any, D IDType] struct {
	Added   []T
	Edited  []T
	Removed []D
}

type IDer[D IDType] interface {
	GetID() D
}

func MergeItems[T IDer[D], D IDType](originalList []T, mutateList *MutateList[T, D]) ([]T, error) {
	removedMap := make(map[D]struct{}, len(mutateList.Removed))
	removedFound := 0

	for _, item := range mutateList.Removed {
		removedMap[item] = struct{}{}
	}

	editedMap := make(map[D]T, len(mutateList.Edited))
	editedFound := 0

	for _, item := range mutateList.Edited {
		editedMap[item.GetID()] = item
	}

	finalItems := make([]T, 0, len(originalList)+len(mutateList.Added))
	for _, item := range originalList {
		if _, ok := removedMap[item.GetID()]; ok {
			removedFound += 1

			continue
		}

		if newItem, ok := editedMap[item.GetID()]; ok {
			editedFound += 1
			item = newItem
		}

		finalItems = append(finalItems, item)
	}

	if removedFound != len(mutateList.Removed) || editedFound != len(mutateList.Edited) {
		return nil, ErrInvalidID
	}

	finalItems = append(finalItems, mutateList.Added...)

	return finalItems, nil
}
