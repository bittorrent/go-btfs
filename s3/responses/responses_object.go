package responses

import (
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func WritePutObjectResponse(w http.ResponseWriter, r *http.Request, obj object.Object) {
	setPutObjHeaders(w, obj.ETag, obj.CID, false)
	WriteSuccessResponse(w, nil, "")
}
