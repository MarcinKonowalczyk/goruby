package walk

import (
	"fmt"
	"sort"
)

func deleteIndices[T any](
	slice []T,
	indices []int,
) []T {
	sorted := sort.IntSlice(indices)
	sorted.Sort()

	unique := make([]int, 0, len(sorted))
	for _, index := range sorted {
		if len(unique) == 0 || unique[len(unique)-1] != index {
			unique = append(unique, index)
		}
	}

	for i := len(unique) - 1; i >= 0; i-- {
		index := sorted[i]
		if index < 0 || index >= len(slice) {
			panic(fmt.Sprintf("deleteIndices: index %d out of bounds for slice of length %d", index, len(slice)))
		}
		slice = append(slice[:index], slice[index+1:]...)
	}
	return slice
}
