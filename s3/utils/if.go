package utils

func IfEmpty(a, b string) (c string) {
	c = a
	if a == "" {
		c = b
	}
	return
}

func IfZero(a, b int) (c int) {
	c = a
	if a == 0 {
		c = b
	}
	return
}
