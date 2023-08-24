package responses

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/gorilla/mux"
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

type RESTErrorResponse struct {
	XMLName    xml.Name `xml:"Error" json:"-"`
	Code       string   `xml:"Code" json:"Code"`
	Message    string   `xml:"Message" json:"Message"`
	Resource   string   `xml:"Resource" json:"Resource"`
	RequestID  string   `xml:"RequestId" json:"RequestId"`
	Key        string   `xml:"Key,omitempty" json:"Key,omitempty"`
	BucketName string   `xml:"BucketName,omitempty" json:"BucketName,omitempty"`
}

func getRESTErrorResponse(err *Error, resource string, bucket, object string) RESTErrorResponse {
	return RESTErrorResponse{
		Code:       err.Code(),
		BucketName: bucket,
		Key:        object,
		Message:    err.Description(),
		Resource:   resource,
		RequestID:  fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}

func WriteErrorResponseHeadersOnly(w http.ResponseWriter, r *http.Request, err error) {
	var rerr *Error
	if !errors.As(err, &rerr) {
		rerr = ErrInternalError
	}
	writeResponse(w, r, rerr.HTTPStatusCode(), nil, mimeNone)
}

// WriteErrorResponse write ErrorResponse
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var rerr *Error
	if !errors.As(err, &rerr) {
		rerr = ErrInternalError
	}
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]
	errorResponse := RESTErrorResponse{
		Code:       rerr.Code(),
		BucketName: bucket,
		Key:        object,
		Message:    rerr.Description(),
		Resource:   r.URL.Path,
		RequestID:  fmt.Sprintf("%d", time.Now().UnixNano()),
	}
	WriteXMLResponse(w, r, rerr.HTTPStatusCode(), errorResponse)
}

// WriteSuccessResponse write SuccessResponseHeadersOnly
func WriteSuccessResponse(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, r, http.StatusOK, nil, mimeNone)
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
func WriteErrorResponseJSON(w http.ResponseWriter, err error, reqURL *url.URL, host string) {
	var rerr *Error
	if !errors.As(err, &rerr) {
		rerr = ErrInternalError
	}
	// Generate error response.
	errorResponse := getAPIErrorResponse(rerr, reqURL.Path, w.Header().Get(consts.AmzRequestID), host)
	encodedErrorResponse := encodeResponseJSON(errorResponse)
	writeResponseSimple(w, rerr.HTTPStatusCode(), encodedErrorResponse, mimeJSON)
}

// getErrorResponse gets in standard error and resource value and
// provides a encodable populated response values
func getAPIErrorResponse(err *Error, resource, requestID, hostID string) APIErrorResponse {
	return APIErrorResponse{
		Code:      err.Code(),
		Message:   err.Description(),
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
