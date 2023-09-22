package utils

// CoalesceStr return the first non-empty string in the list
func CoalesceStr(list ...string) string {
	for _, str := range list {
		if str != "" {
			return str
		}
	}
	return ""
}
