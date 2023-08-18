package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func WriteErrorResponseHeadersOnly(w http.ResponseWriter, r *http.Request, err ErrorCode) {
	writeResponse(w, r, GetAPIError(err).HTTPStatusCode, nil, mimeNone)
}

// WriteErrorResponse write ErrorResponse
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, errorCode ErrorCode) {
	fmt.Println("response errcode: ", errorCode)
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]

	apiError := GetAPIError(errorCode)
	errorResponse := getRESTErrorResponse(apiError, r.URL.Path, bucket, object)
	WriteXMLResponse(w, r, apiError.HTTPStatusCode, errorResponse)
}

func getRESTErrorResponse(err APIError, resource string, bucket, object string) RESTErrorResponse {
	return RESTErrorResponse{
		Code:       err.Code,
		BucketName: bucket,
		Key:        object,
		Message:    err.Description,
		Resource:   resource,
		RequestID:  fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}

// NotFoundHandler If none of the http routes match respond with MethodNotAllowed
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
}
