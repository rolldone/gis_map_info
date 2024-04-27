package helper

func ReverseArray[T any](arr []T) []T {
	result := make([]T, len(arr))
	length := len(arr)
	for index := range arr {
		result[(length-1)-index] = arr[index]
	}
	return result
}
