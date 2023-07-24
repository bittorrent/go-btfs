package authentication

type AuthErr struct {
}

func (err AuthErr) Error() string {
	return ""
}
