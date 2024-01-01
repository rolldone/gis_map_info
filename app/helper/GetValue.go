package helper

var MapEmpty map[string]interface{}

func GetValue(val1 interface{}, alt interface{}) interface{} {
	if val1 != nil {
		return val1
	}
	return alt
}
