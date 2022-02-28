package core

// define some helper function used in core

// CloneSlice will clone the source slice
func CloneSlice[T any](srcSlice []T) []T {
	dst := make([]T, len(srcSlice))
	copy(dst, srcSlice)
	return dst
}
