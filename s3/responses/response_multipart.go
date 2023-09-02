package responses

import (
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func WriteCreateMultipartUploadResponse(w http.ResponseWriter, r *http.Request, bucname, objname, uploadID string) {
	resp := GenerateInitiateMultipartUploadResponse(bucname, objname, uploadID)
	WriteSuccessResponseXML(w, r, resp)
}

func WriteAbortMultipartUploadResponse(w http.ResponseWriter, r *http.Request) {
	WriteSuccessNoContent(w)
}

func WriteUploadPartResponse(w http.ResponseWriter, r *http.Request, part object.Part) {
	setPutObjHeaders(w, part.ETag, part.CID, false)
	WriteSuccessResponseHeadersOnly(w, r)
}

func WriteCompleteMultipartUploadResponse(w http.ResponseWriter, r *http.Request, bucname, objname, region string, obj object.Object) {
	resp := GenerateCompleteMultipartUploadResponse(bucname, objname, region, obj)
	setPutObjHeaders(w, obj.ETag, obj.CID, false)
	WriteSuccessResponseXML(w, r, resp)
}
