package utils

func Coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return "" // return an empty string if all values are empty
}

func GetHeaderVarToString(value any, exists bool) string {

	if !exists {
		return ""
	}

	return value.(string)
}
