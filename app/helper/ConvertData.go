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
