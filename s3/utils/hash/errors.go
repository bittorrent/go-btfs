package hash

import "fmt"

// SHA256Mismatch - when content sha256 does not match with what was sent from client.
type SHA256Mismatch struct {
	ExpectedSHA256   string
	CalculatedSHA256 string
}

func (e SHA256Mismatch) Error() string {
	return "Bad sha256: Expected " + e.ExpectedSHA256 + " does not match calculated " + e.CalculatedSHA256
}

// BadDigest - Content-MD5 you specified did not match what we received.
type BadDigest struct {
	ExpectedMD5   string
	CalculatedMD5 string
}

func (e BadDigest) Error() string {
	return "Bad digest: Expected " + e.ExpectedMD5 + " does not match calculated " + e.CalculatedMD5
}

// ErrSizeMismatch error size mismatch
type ErrSizeMismatch struct {
	Want int64
	Got  int64
}

func (e ErrSizeMismatch) Error() string {
	return fmt.Sprintf("Size mismatch: got %d, want %d", e.Got, e.Want)
}
