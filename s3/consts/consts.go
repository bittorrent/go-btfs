package consts

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dustin/go-humanize"
)

const (
	StreamingContentSHA256   = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"
	EmptySHA256              = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	UnsignedSHA256           = "UNSIGNED-PAYLOAD"
	SlashSeparator           = "/"
	StsAction                = "Action"
	StreamingContentEncoding = "aws-chunked"
	DefaultEncodingType      = "url"
	DefaultContentType       = "binary/octet-stream"
	DefaultServerInfo        = "BTFS"
	DefaultBucketRegion      = "us-east-1"
	DefaultBucketACL         = s3.BucketCannedACLPublicRead
	AllUsersURI              = "http://acs.amazonaws.com/groups/global/AllUsers"
)

var SupportedBucketRegions = map[string]bool{
	DefaultBucketRegion: true,
}

var SupportedBucketACLs = map[string]bool{
	s3.BucketCannedACLPrivate:         true,
	s3.BucketCannedACLPublicRead:      true,
	s3.BucketCannedACLPublicReadWrite: true,
}

// Standard S3 HTTP request constants
const (
	AmzACL           = "x-amz-acl"
	AmzContentSha256 = "X-Amz-Content-Sha256"
	AmzDate          = "X-Amz-Date"
	AmzRequestID     = "x-amz-request-id"
)

// Standard S3 HTTP response constants
const (
	LastModified       = "Last-Modified"
	Date               = "Date"
	ETag               = "ETag"
	ContentType        = "Content-Type"
	ContentMD5         = "Content-Md5"
	ContentEncoding    = "Content-Encoding"
	Expires            = "Expires"
	ContentLength      = "Content-Length"
	ContentLanguage    = "Content-Language"
	ContentRange       = "Content-Range"
	Connection         = "Connection"
	AcceptRanges       = "Accept-Ranges"
	AmzBucketRegion    = "X-Amz-Bucket-Region"
	ServerInfo         = "Server"
	RetryAfter         = "Retry-After"
	Location           = "Location"
	CacheControl       = "Cache-Control"
	ContentDisposition = "Content-Disposition"
	Authorization      = "Authorization"
	Action             = "Action"
	XRequestWith       = "X-Requested-With"
	Range              = "Range"
	UserAgent          = "User-Agent"
	Cid                = "Cid"
)

// Standard HTTP cors headers
const (
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlMaxAge           = "Access-Control-Max-Age"
)

// object const
const (
	MaxXMLBodySize = 5 * humanize.MiByte
	MaxObjectSize  = 5 * humanize.TiByte
	MinPartSize    = 5 * humanize.MiByte
	MaxPartSize    = 5 * humanize.GiByte
	MinPartNumber  = 1
	MaxPartNumber  = 10000
	MaxObjectList  = 1000 // Limit number of objects in a listObjectsResponse/listObjectsVersionsResponse.
	MaxDeleteList  = 1000 // Limit number of objects deleted in a delete call.
)

// Common http query params S3 API
const (
	MaxKeys    = "max-keys"
	PartNumber = "partNumber"
)
