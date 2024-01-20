package helper

import "encoding/json"

func ToStructFromMap(c interface{}, v any) error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, &v)
	return nil
}

// SliceOfInt converts the value into a slice of ints.
// It works with interface{}, []interface{}, []string, []int, and int values.
// If the passed value cannot be converted, then an empty slice is returned.
func SliceOfInt64(value2 interface{}) []int64 {

	switch value := value2.(type) {

	case []int:
		return nil

	case []interface{}:
		int64Slice := make([]int64, len(value))
		// Convert each element to int64
		for i, v := range value {
			int64Slice[i] = int64(v.(float64))
		}
		return int64Slice

	case []string:
		return nil

	case int:
		return nil
	}

	return nil
}
