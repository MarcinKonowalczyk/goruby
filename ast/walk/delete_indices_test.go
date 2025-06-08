package walk

import (
	"testing"

	"github.com/MarcinKonowalczyk/assert"
)

func TestDeleteIndices(t *testing.T) {

	t.Run("delete indices from slice", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		indices := []int{1, 3}
		result := deleteIndices(arr, indices)
		expected := []int{1, 3, 5}

		assert.EqualArrays(t, expected, result)
	})

	t.Run("delete indices with duplicates", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		indices := []int{1, 3, 3}
		result := deleteIndices(arr, indices)
		expected := []int{1, 3, 5}

		assert.EqualArrays(t, expected, result)
	})

	t.Run("delete indices out of order", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		indices := []int{3, 1}
		result := deleteIndices(arr, indices)
		expected := []int{1, 3, 5}

		assert.EqualArrays(t, expected, result)
	})

	t.Run("delete indices with negative index", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		indices := []int{-1, 1}
		assert.Panic(t, func() { deleteIndices(arr, indices) }, func(t testing.TB, rec any) {
			assert.ContainsString(
				t,
				assert.Type[string](t, rec),
				"index -1 out of bounds",
			)
		})
	})

	t.Run("delete indices with out of bounds index", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		indices := []int{5} // index 5 is out of bounds for a slice of length 5
		assert.Panic(t, func() { deleteIndices(arr, indices) }, func(t testing.TB, rec any) {
			assert.ContainsString(
				t,
				assert.Type[string](t, rec),
				"index 5 out of bounds",
			)
		})
	})
	t.Run("delete indices with empty slice", func(t *testing.T) {
		arr := []int{}
		indices := []int{0} // index 0 is out of bounds for an empty slice
		assert.Panic(t, func() { deleteIndices(arr, indices) }, func(t testing.TB, rec any) {
			assert.ContainsString(
				t,
				assert.Type[string](t, rec),
				"index 0 out of bounds",
			)
		})
	})
	t.Run("delete indices with empty indices", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		indices := []int{} // no indices to delete
		result := deleteIndices(arr, indices)
		expected := []int{1, 2, 3, 4, 5}

		assert.EqualArrays(t, expected, result)
	})
}
