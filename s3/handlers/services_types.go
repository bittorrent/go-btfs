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

type ObjectMetadata struct {
}
