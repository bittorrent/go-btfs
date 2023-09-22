package responses

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
)

func WriteCreateMultipartUploadResponse(w http.ResponseWriter, r *http.Request, multipart *object.Multipart) {
	output := new(s3.CreateMultipartUploadOutput)
	output.SetBucket(multipart.Bucket)
	output.SetKey(multipart.Object)
	output.SetUploadId(multipart.UploadID)
	WriteSuccessResponse(w, output, "InitiateMultipartUploadResult")
}

func WriteUploadPartResponse(w http.ResponseWriter, r *http.Request, part *object.Part) {
	output := new(s3.UploadPartOutput)
	output.SetETag(`"` + part.ETag + `"`)
	w.Header().Set(consts.Cid, part.CID)
	WriteSuccessResponse(w, output, "")
}

func WriteAbortMultipartUploadResponse(w http.ResponseWriter, r *http.Request) {
	output := new(s3.AbortMultipartUploadOutput)
	WriteSuccessResponse(w, output, "")
}

func WriteCompleteMultipartUploadResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.CompleteMultipartUploadOutput)
	output.SetBucket(obj.Bucket)
	output.SetKey(obj.Name)
	output.SetETag(`"` + obj.ETag + `"`)
	w.Header().Set(consts.Cid, obj.CID)
	WriteSuccessResponse(w, output, "CompleteMultipartUploadResult")
}
