package helper

type Type int

func RemoveIndex[T any](s []T, index int) []T {
	if s == nil {
		return nil // Return nil if the input slice is nil.
	}

	if index < 0 || index >= len(s) {
		return s // Return the original slice if the index is out of range.
	}

	// Remove the element at the specified index.
	return append(s[:index], s[index+1:]...)
}
