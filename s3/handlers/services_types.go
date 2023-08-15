package handlers

import "time"

type AccessKeyRecord struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Enable    bool      `json:"enable"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BucketMetadata contains bucket metadata.
type BucketMetadata struct {
	Name    string
	Region  string
	Owner   string
	Acl     string
	Created time.Time
}

// NewBucketMetadata creates BucketMetadata with the supplied name and Created to Now.
func NewBucketMetadata(name, region, accessKey, acl string) *BucketMetadata {
	return &BucketMetadata{
		Name:    name,
		Region:  region,
		Owner:   accessKey,
		Acl:     acl,
		Created: time.Now().UTC(),
	}
}
