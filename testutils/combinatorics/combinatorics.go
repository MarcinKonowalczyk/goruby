package combinatorics

import (
	"sort"
)

func Combinations[T any](arr []T, k int) [][]T {
	if k == 0 {
		return [][]T{{}}
	}
	if len(arr) < k {
		return nil
	}

	var result [][]T
	for i, v := range arr {
		for _, comb := range Combinations(arr[i+1:], k-1) {
			result = append(result, append([]T{v}, comb...))
		}
	}
	return result
}

func CombinationsSorted[T any](arr []T, k int, less func(a, b T) bool) [][]T {
	combs := Combinations(arr, k)
	// sort each combination
	for i := range combs {
		sort.Slice(combs[i], func(a, b int) bool {
			return less(combs[i][a], combs[i][b])
		})
	}
	// sort the combinations themselves
	sort.Slice(combs, func(i, j int) bool {
		if len(combs[i]) != len(combs[j]) {
			// if lengths are different (which should not happen in this context), sort by length
			return len(combs[i]) < len(combs[j])
		}
		// compare each element
		for idx := 0; idx < len(combs[i]); idx++ {
			ci := combs[i][idx]
			cj := combs[j][idx]
			if less(ci, cj) {
				return true
			} else if less(cj, ci) {
				// NOTE: we don't have equality here, so we need two less checks
				return false
			} else {
				// elements are equal, continue to the next element
			}
		}
		return false
	})
	return combs
}
