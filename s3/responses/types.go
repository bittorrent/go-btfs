package responses

import (
	"encoding/xml"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/services/object"
)

type AccessControlList struct {
	Grant []*s3.Grant `xml:"Grant,omitempty"`
}

type CanonicalUser struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName,omitempty"`
}

// Grant grant
type Grant struct {
	Grantee    Grantee    `xml:"Grantee"`
	Permission Permission `xml:"Permission"`
}

// Grantee grant
type Grantee struct {
	XMLNS       string `xml:"xmlns:xsi,attr"`
	XMLXSI      string `xml:"xsi:type,attr"`
	Type        string `xml:"Type"`
	ID          string `xml:"ID,omitempty"`
	DisplayName string `xml:"DisplayName,omitempty"`
	URI         string `xml:"URI,omitempty"`
}

// Permission May be one of READ, WRITE, READ_ACP, WRITE_ACP, FULL_CONTROL
type Permission string

// ListAllMyBucketsResult  List All Buckets Result
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

type CopyObjectResponse struct {
	CopyObjectResult CopyObjectResult `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CopyObjectResult"`
}

type CopyObjectResult struct {
	LastModified string `xml:"http://s3.amazonaws.com/doc/2006-03-01/ LastModified"`
	ETag         string `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ETag"`
}

// LocationResponse - format for location response.
type LocationResponse struct {
	XMLName  xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ LocationConstraint" json:"-"`
	Location string   `xml:",chardata"`
}

// ListObjectsResponse - format for list objects response.
type ListObjectsResponse struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult" json:"-"`

	Name   string
	Prefix string
	Marker string

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects. Server lists objects in alphabetical
	// order Note: This element is returned only if you have delimiter request parameter
	// specified. If response does not include the NextMaker and it is truncated,
	// you can use the value of the last Key in the response as the marker in the
	// subsequent request to get the next set of object keys.
	NextMarker string `xml:"NextMarker,omitempty"`

	MaxKeys   int
	Delimiter string
	// A flag that indicates whether or not ListObjects returned all of the results
	// that satisfied the search criteria.
	IsTruncated bool

	Contents       []Object
	CommonPrefixes []CommonPrefix

	// Encoding type used to encode object keys in the response.
	EncodingType string `xml:"EncodingType,omitempty"`
}

// ListObjectsV2Response - format for list objects response.
type ListObjectsV2Response struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult" json:"-"`

	Name       string
	Prefix     string
	StartAfter string `xml:"StartAfter,omitempty"`
	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects. Server lists objects in alphabetical
	// order Note: This element is returned only if you have delimiter request parameter
	// specified. If response does not include the NextMaker and it is truncated,
	// you can use the value of the last Key in the response as the marker in the
	// subsequent request to get the next set of object keys.
	ContinuationToken     string `xml:"ContinuationToken,omitempty"`
	NextContinuationToken string `xml:"NextContinuationToken,omitempty"`

	KeyCount  int
	MaxKeys   int
	Delimiter string
	// A flag that indicates whether or not ListObjects returned all of the results
	// that satisfied the search criteria.
	IsTruncated bool

	Contents       []Object
	CommonPrefixes []CommonPrefix

	// Encoding type used to encode object keys in the response.
	EncodingType string `xml:"EncodingType,omitempty"`
}

// Object container for object metadata
type Object struct {
	Key          string
	LastModified string // time string of format "2006-01-02T15:04:05.000Z"
	ETag         string
	CID          string // CID
	Size         int64

	// Owner of the object.
	Owner s3.Owner

	// The class of storage used to store the object.
	StorageClass string

	// UserMetadata user-defined metadata
	UserMetadata StringMap `xml:"UserMetadata,omitempty"`
}

// StringMap is a map[string]string
type StringMap map[string]string

// MarshalXML - StringMap marshals into XML.
func (s StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{}
		t.Name = xml.Name{
			Space: "",
			Local: key,
		}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{
		Name: start.Name,
	})

	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}

// CommonPrefix container for prefix response in ListObjectsResponse
type CommonPrefix struct {
	Prefix string
}

type InitiateMultipartUploadResponse struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ InitiateMultipartUploadResult" json:"-"`

	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

func GenerateInitiateMultipartUploadResponse(bucname, objname, uploadID string) InitiateMultipartUploadResponse {
	return InitiateMultipartUploadResponse{
		Bucket:   bucname,
		Key:      objname,
		UploadID: uploadID,
	}
}

type CompleteMultipartUploadResponse struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CompleteMultipartUploadResult" json:"-"`

	Location string
	Bucket   string
	Key      string
	ETag     string

	ChecksumCRC32  string
	ChecksumCRC32C string
	ChecksumSHA1   string
	ChecksumSHA256 string
}

func GenerateCompleteMultipartUploadResponse(bucname, objname, location string, obj object.Object) CompleteMultipartUploadResponse {
	c := CompleteMultipartUploadResponse{
		Location: location,
		Bucket:   bucname,
		Key:      objname,
		// AWS S3 quotes the ETag in XML, make sure we are compatible here.
		ETag: "\"" + obj.ETag + "\"",
	}
	return c
}

// GenerateListObjectsV2Response Generates an ListObjectsV2 response for the said bucket with other enumerated options.
//func GenerateListObjectsV2Response(bucket, prefix, token, nextToken, startAfter, delimiter, encodingType string, isTruncated bool, maxKeys int, objects []object.Object, prefixes []string) ListObjectsV2Response {
//	contents := make([]Object, 0, len(objects))
//	id := consts.DefaultOwnerID
//	name := consts.DisplayName
//	owner := s3.Owner{
//		ID:          &id,
//		DisplayName: &name,
//	}
//	data := ListObjectsV2Response{}
//
//	for _, object := range objects {
//		content := Object{}
//		if object.Name == "" {
//			continue
//		}
//		content.Key = utils.S3EncodeName(object.Name, encodingType)
//		content.LastModified = object.ModTime.UTC().Format(consts.Iso8601TimeFormat)
//		if object.ETag != "" {
//			content.ETag = "\"" + object.ETag + "\""
//		}
//		content.Size = object.Size
//		content.Owner = owner
//		content.CID = object.CID
//		contents = append(contents, content)
//	}
//	data.Name = bucket
//	data.Contents = contents
//
//	data.EncodingType = encodingType
//	data.StartAfter = utils.S3EncodeName(startAfter, encodingType)
//	data.Delimiter = utils.S3EncodeName(delimiter, encodingType)
//	data.Prefix = utils.S3EncodeName(prefix, encodingType)
//	data.MaxKeys = maxKeys
//	data.ContinuationToken = base64.StdEncoding.EncodeToString([]byte(token))
//	data.NextContinuationToken = base64.StdEncoding.EncodeToString([]byte(nextToken))
//	data.IsTruncated = isTruncated
//
//	commonPrefixes := make([]CommonPrefix, 0, len(prefixes))
//	for _, prefix := range prefixes {
//		prefixItem := CommonPrefix{}
//		prefixItem.Prefix = utils.S3EncodeName(prefix, encodingType)
//		commonPrefixes = append(commonPrefixes, prefixItem)
//	}
//	data.CommonPrefixes = commonPrefixes
//	data.KeyCount = len(data.Contents) + len(data.CommonPrefixes)
//	return data
//}

// generates an ListObjectsV1 response for the said bucket with other enumerated options.
//func GenerateListObjectsV1Response(bucket, prefix, marker, delimiter, encodingType string, maxKeys int, resp object.ObjectsList) ListObjectsResponse {
//	contents := make([]Object, 0, len(resp.Objects))
//	id := consts.DefaultOwnerID
//	name := consts.DisplayName
//	owner := s3.Owner{
//		ID:          &id,
//		DisplayName: &name,
//	}
//	data := ListObjectsResponse{}
//
//	for _, object := range resp.Objects {
//		content := Object{}
//		if object.Name == "" {
//			continue
//		}
//		content.Key = utils.S3EncodeName(object.Name, encodingType)
//		content.LastModified = object.ModTime.UTC().Format(consts.Iso8601TimeFormat)
//		if object.ETag != "" {
//			content.ETag = "\"" + object.ETag + "\""
//		}
//		content.CID = object.CID
//		content.Size = object.Size
//		content.StorageClass = ""
//		content.Owner = owner
//		contents = append(contents, content)
//	}
//	data.Name = bucket
//	data.Contents = contents
//
//	data.EncodingType = encodingType
//	data.Prefix = utils.S3EncodeName(prefix, encodingType)
//	data.Marker = utils.S3EncodeName(marker, encodingType)
//	data.Delimiter = utils.S3EncodeName(delimiter, encodingType)
//	data.MaxKeys = maxKeys
//	data.NextMarker = utils.S3EncodeName(resp.NextMarker, encodingType)
//	data.IsTruncated = resp.IsTruncated
//
//	prefixes := make([]CommonPrefix, 0, len(resp.Prefixes))
//	for _, prefix := range resp.Prefixes {
//		prefixItem := CommonPrefix{}
//		prefixItem.Prefix = utils.S3EncodeName(prefix, encodingType)
//		prefixes = append(prefixes, prefixItem)
//	}
//	data.CommonPrefixes = prefixes
//	return data
//}
//
