package utils

func CloneMapSS(src map[string]string) map[string]string {
	r := make(map[string]string, len(src))
	for k, v := range src {
		r[k] = v
	}
	return r
}
