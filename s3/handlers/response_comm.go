package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var log = logging.Logger("resp")

type mimeType string

const (
	mimeNone mimeType = ""
	mimeJSON mimeType = "application/json"
	//mimeXML application/xml UTF-8
	mimeXML mimeType = " application/xml"
)

// APIErrorResponse - error response format
type APIErrorResponse struct {
	XMLName   xml.Name `xml:"Error" json:"-"`
	Code      string
	Message   string
	Resource  string
	RequestID string `xml:"RequestId" json:"RequestId"`
	HostID    string `xml:"HostId" json:"HostId"`
}

// WriteSuccessResponseXML Write Success Response XML
func WriteSuccessResponseXML(w http.ResponseWriter, r *http.Request, response interface{}) {
	WriteXMLResponse(w, r, http.StatusOK, response)
}

// WriteXMLResponse Write XMLResponse
func WriteXMLResponse(w http.ResponseWriter, r *http.Request, statusCode int, response interface{}) {
	writeResponse(w, r, statusCode, encodeXMLResponse(response), mimeXML)
}

func writeResponse(w http.ResponseWriter, r *http.Request, statusCode int, response []byte, mType mimeType) {
	setCommonHeaders(w, r)
	if response != nil {
		w.Header().Set(consts.ContentLength, strconv.Itoa(len(response)))
	}
	if mType != mimeNone {
		w.Header().Set(consts.ContentType, string(mType))
	}
	w.WriteHeader(statusCode)
	if response != nil {
		log.Debugf("status %d %s: %s", statusCode, mType, string(response))
		_, err := w.Write(response)
		if err != nil {
			log.Errorf("write err: %v", err)
		}
		w.(http.Flusher).Flush()
	}
}

func setCommonHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(consts.ServerInfo, "FDS")
	w.Header().Set(consts.AmzRequestID, fmt.Sprintf("%d", time.Now().UnixNano()))
	w.Header().Set(consts.AcceptRanges, "bytes")
	if r.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

// encodeXMLResponse Encodes the response headers into XML format.
func encodeXMLResponse(response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(response)
	return bytesBuffer.Bytes()
}

// WriteErrorResponseJSON - writes error response in JSON format;
// useful for admin APIs.
func WriteErrorResponseJSON(w http.ResponseWriter, err APIError, reqURL *url.URL, host string) {
	// Generate error response.
	errorResponse := getAPIErrorResponse(err, reqURL.Path, w.Header().Get(consts.AmzRequestID), host)
	encodedErrorResponse := encodeResponseJSON(errorResponse)
	writeResponseSimple(w, err.HTTPStatusCode, encodedErrorResponse, mimeJSON)
}

// getErrorResponse gets in standard error and resource value and
// provides a encodable populated response values
func getAPIErrorResponse(err APIError, resource, requestID, hostID string) APIErrorResponse {
	return APIErrorResponse{
		Code:      err.Code,
		Message:   err.Description,
		Resource:  resource,
		RequestID: requestID,
		HostID:    hostID,
	}
}

// Encodes the response headers into JSON format.
func encodeResponseJSON(response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	e := json.NewEncoder(&bytesBuffer)
	e.Encode(response)
	return bytesBuffer.Bytes()
}

// WriteSuccessResponseJSON writes success headers and response if any,
// with content-type set to `application/json`.
func WriteSuccessResponseJSON(w http.ResponseWriter, response []byte) {
	writeResponseSimple(w, http.StatusOK, response, mimeJSON)
}

func writeResponseSimple(w http.ResponseWriter, statusCode int, response []byte, mType mimeType) {
	if mType != mimeNone {
		w.Header().Set(consts.ContentType, string(mType))
	}
	w.Header().Set(consts.ContentLength, strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	if response != nil {
		w.Write(response)
	}
}

// WriteSuccessNoContent writes success headers with http status 204
func WriteSuccessNoContent(w http.ResponseWriter) {
	writeResponseSimple(w, http.StatusNoContent, nil, mimeNone)
}

// ListAllMyBucketsResult  List All Buckets Result
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

// WriteSuccessResponseHeadersOnly write SuccessResponseHeadersOnly
func WriteSuccessResponseHeadersOnly(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, r, http.StatusOK, nil, mimeNone)
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

//
//// DeleteError structure.
//type DeleteError struct {
//	Code      string
//	Message   string
//	Key       string
//	VersionID string `xml:"VersionId"`
//}
//
//// DeleteObjectsResponse container for multiple object deletes.
//type DeleteObjectsResponse struct {
//	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ DeleteResult" json:"-"`
//
//	// Collection of all deleted objects
//	DeletedObjects []datatypes.DeletedObject `xml:"Deleted,omitempty"`
//
//	// Collection of errors deleting certain objects.
//	Errors []DeleteError `xml:"Error,omitempty"`
//}
//
//// GenerateListObjectsV2Response Generates an ListObjectsV2 response for the said bucket with other enumerated options.
//func GenerateListObjectsV2Response(bucket, prefix, token, nextToken, startAfter, delimiter, encodingType string, isTruncated bool, maxKeys int, objects []store.ObjectInfo, prefixes []string) ListObjectsV2Response {
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
//
//// generates an ListObjectsV1 response for the said bucket with other enumerated options.
//func GenerateListObjectsV1Response(bucket, prefix, marker, delimiter, encodingType string, maxKeys int, resp store.ListObjectsInfo) ListObjectsResponse {
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
//// generate multi objects delete response.
//func GenerateMultiDeleteResponse(quiet bool, deletedObjects []datatypes.DeletedObject, errs []DeleteError) DeleteObjectsResponse {
//	deleteResp := DeleteObjectsResponse{}
//	if !quiet {
//		deleteResp.DeletedObjects = deletedObjects
//	}
//	deleteResp.Errors = errs
//	return deleteResp
//}
//
//// InitiateMultipartUploadResponse container for InitiateMultiPartUpload response, provides uploadID to start MultiPart upload
//type InitiateMultipartUploadResponse struct {
//	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ InitiateMultipartUploadResult" json:"-"`
//
//	Bucket   string
//	Key      string
//	UploadID string `xml:"UploadId"`
//}
//
//// CompleteMultipartUploadResponse container for completed multipart upload response
//type CompleteMultipartUploadResponse struct {
//	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CompleteMultipartUploadResult" json:"-"`
//
//	Location string
//	Bucket   string
//	Key      string
//	ETag     string
//
//	ChecksumCRC32  string
//	ChecksumCRC32C string
//	ChecksumSHA1   string
//	ChecksumSHA256 string
//}
//
//// Part container for part metadata.
//type Part struct {
//	PartNumber   int
//	LastModified string
//	ETag         string
//	Size         int64
//
//	// Checksum values
//	ChecksumCRC32  string
//	ChecksumCRC32C string
//	ChecksumSHA1   string
//	ChecksumSHA256 string
//}
//
//// Initiator inherit from Owner struct, fields are same
//type Initiator s3.Owner
//
//// ListPartsResponse - format for list parts response.
//type ListPartsResponse struct {
//	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListPartsResult" json:"-"`
//
//	Bucket   string
//	Key      string
//	UploadID string `xml:"UploadId"`
//
//	Initiator Initiator
//	Owner     s3.Owner
//
//	// The class of storage used to store the object.
//	StorageClass string
//
//	PartNumberMarker     int
//	NextPartNumberMarker int
//	MaxParts             int
//	IsTruncated          bool
//
//	ChecksumAlgorithm string
//	// List of parts.
//	Parts []Part `xml:"Part"`
//}
//
//// ListMultipartUploadsResponse - format for list multipart uploads response.
//type ListMultipartUploadsResponse struct {
//	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListMultipartUploadsResult" json:"-"`
//
//	Bucket             string
//	KeyMarker          string
//	UploadIDMarker     string `xml:"UploadIdMarker"`
//	NextKeyMarker      string
//	NextUploadIDMarker string `xml:"NextUploadIdMarker"`
//	Delimiter          string
//	Prefix             string
//	EncodingType       string `xml:"EncodingType,omitempty"`
//	MaxUploads         int
//	IsTruncated        bool
//
//	// List of pending uploads.
//	Uploads []Upload `xml:"Upload"`
//
//	// Delimed common prefixes.
//	CommonPrefixes []CommonPrefix
//}
//
//// Upload container for in progress multipart upload
//type Upload struct {
//	Key          string
//	UploadID     string `xml:"UploadId"`
//	Initiator    Initiator
//	Owner        s3.Owner
//	StorageClass string
//	Initiated    string
//}
//
//// generates InitiateMultipartUploadResponse for given bucket, key and uploadID.
//func GenerateInitiateMultipartUploadResponse(bucket, key, uploadID string) InitiateMultipartUploadResponse {
//	return InitiateMultipartUploadResponse{
//		Bucket:   bucket,
//		Key:      key,
//		UploadID: uploadID,
//	}
//}
//
//// generates CompleteMultipartUploadResponse for given bucket, key, location and ETag.
//func GenerateCompleteMultpartUploadResponse(bucket, key, location string, oi store.ObjectInfo) CompleteMultipartUploadResponse {
//	c := CompleteMultipartUploadResponse{
//		Location: location,
//		Bucket:   bucket,
//		Key:      key,
//		// AWS S3 quotes the ETag in XML, make sure we are compatible here.
//		ETag: "\"" + oi.ETag + "\"",
//	}
//	return c
//}
//
//// generates ListPartsResponse from ListPartsInfo.
//func GenerateListPartsResponse(partsInfo store.ListPartsInfo, encodingType string) ListPartsResponse {
//	resp := ListPartsResponse{}
//	resp.Bucket = partsInfo.Bucket
//	resp.Key = utils.S3EncodeName(partsInfo.Object, encodingType)
//	resp.UploadID = partsInfo.UploadID
//	resp.StorageClass = consts.DefaultStorageClass
//
//	// Dumb values not meaningful
//	resp.Initiator = Initiator{
//		ID:          aws.String(consts.DefaultOwnerID),
//		DisplayName: aws.String(consts.DisplayName),
//	}
//	resp.Owner = s3.Owner{
//		ID:          aws.String(consts.DefaultOwnerID),
//		DisplayName: aws.String(consts.DisplayName),
//	}
//
//	resp.MaxParts = partsInfo.MaxParts
//	resp.PartNumberMarker = partsInfo.PartNumberMarker
//	resp.IsTruncated = partsInfo.IsTruncated
//	resp.NextPartNumberMarker = partsInfo.NextPartNumberMarker
//	resp.ChecksumAlgorithm = partsInfo.ChecksumAlgorithm
//
//	resp.Parts = make([]Part, len(partsInfo.Parts))
//	for index, part := range partsInfo.Parts {
//		newPart := Part{}
//		newPart.PartNumber = part.Number
//		newPart.ETag = "\"" + part.ETag + "\""
//		newPart.Size = part.Size
//		newPart.LastModified = part.ModTime.UTC().Format(consts.Iso8601TimeFormat)
//		resp.Parts[index] = newPart
//	}
//	return resp
//}
//
//// generates ListMultipartUploadsResponse for given bucket and ListMultipartsInfo.
//func GenerateListMultipartUploadsResponse(bucket string, multipartsInfo store.ListMultipartsInfo, encodingType string) ListMultipartUploadsResponse {
//	resp := ListMultipartUploadsResponse{}
//	resp.Bucket = bucket
//	resp.Delimiter = utils.S3EncodeName(multipartsInfo.Delimiter, encodingType)
//	resp.IsTruncated = multipartsInfo.IsTruncated
//	resp.EncodingType = encodingType
//	resp.Prefix = utils.S3EncodeName(multipartsInfo.Prefix, encodingType)
//	resp.KeyMarker = utils.S3EncodeName(multipartsInfo.KeyMarker, encodingType)
//	resp.NextKeyMarker = utils.S3EncodeName(multipartsInfo.NextKeyMarker, encodingType)
//	resp.MaxUploads = multipartsInfo.MaxUploads
//	resp.NextUploadIDMarker = multipartsInfo.NextUploadIDMarker
//	resp.UploadIDMarker = multipartsInfo.UploadIDMarker
//	resp.CommonPrefixes = make([]CommonPrefix, len(multipartsInfo.CommonPrefixes))
//	for index, commonPrefix := range multipartsInfo.CommonPrefixes {
//		resp.CommonPrefixes[index] = CommonPrefix{
//			Prefix: utils.S3EncodeName(commonPrefix, encodingType),
//		}
//	}
//	resp.Uploads = make([]Upload, len(multipartsInfo.Uploads))
//	for index, upload := range multipartsInfo.Uploads {
//		newUpload := Upload{}
//		newUpload.UploadID = upload.UploadID
//		newUpload.Key = utils.S3EncodeName(upload.Object, encodingType)
//		newUpload.Initiated = upload.Initiated.UTC().Format(consts.Iso8601TimeFormat)
//		resp.Uploads[index] = newUpload
//	}
//	return resp
//}
