package utils

func Coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return "" // return an empty string if all values are empty
}
