package helper

import "strconv"

func ParseInt(val string) int {
	result, _ := strconv.Atoi(val)
	return result
}
