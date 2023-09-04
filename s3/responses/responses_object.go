package responses

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func WritePutObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.PutObjectOutput)
	output.SetETag(`"` + obj.ETag + `"`)
	w.Header().Set(consts.CID, obj.CID)
	WriteSuccessResponse(w, output, "")
}
