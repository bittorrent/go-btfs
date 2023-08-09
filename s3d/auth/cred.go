package auth

import (
	"github.com/bittorrent/go-btfs/s3/apierrors"
	"time"
)

var timeSentinel = time.Unix(0, 0).UTC()

// Credentials holds access and secret keys.
type Credentials struct {
	AccessKey    string    `xml:"AccessKeyId" json:"accessKey,omitempty"`
	SecretKey    string    `xml:"SecretAccessKey" json:"secretKey,omitempty"`
	CreateTime   time.Time `xml:"CreateTime" json:"createTime,omitempty"`
	Expiration   time.Time `xml:"Expiration" json:"expiration,omitempty"`
	SessionToken string    `xml:"SessionToken" json:"sessionToken"`
	Status       string    `xml:"-" json:"status,omitempty"`
	ParentUser   string    `xml:"-" json:"parentUser,omitempty"`
}

// IsValid - returns whether credential is valid or not.
func (cred *Credentials) IsValid() bool {
	return true
}

// IsExpired - returns whether Credential is expired or not.
func (cred *Credentials) IsExpired() bool {
	return false
}

func CheckAccessKeyValid(accessKey string) (*Credentials, apierrors.ErrorCode) {

	////check it
	//cred, bl: =  mp[accessKey]
	//if bl {
	//	return cred, nil
	//} else {
	//	return nil, errors.New("node found accessKey! ")
	//}

	return &Credentials{AccessKey: accessKey}, apierrors.ErrNone
}

const (
	// Minimum length for  access key.
	accessKeyMinLen = 3
	
	// Maximum length for  access key.
	// There is no max length enforcement for access keys
	accessKeyMaxLen = 20

	// Minimum length for  secret key for both server and gateway mode.
	secretKeyMinLen = 8

	// Maximum secret key length , this
	// is used when autogenerating new credentials.
	// There is no max length enforcement for secret keys
	secretKeyMaxLen = 40
)

// IsAccessKeyValid - validate access key for right length.
func IsAccessKeyValid(accessKey string) bool {
	return len(accessKey) >= accessKeyMinLen
}

// IsSecretKeyValid - validate secret key for right length.
func IsSecretKeyValid(secretKey string) bool {
	return len(secretKey) >= secretKeyMinLen
}
