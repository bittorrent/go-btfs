package responses

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SetObjectHeaders Write object header
func SetObjectHeaders(w http.ResponseWriter, r *http.Request, objInfo *object.Object) {
	// Set last modified time.
	lastModified := objInfo.ModTime.UTC().Format(http.TimeFormat)
	w.Header().Set(consts.LastModified, lastModified)

	// Set Etag if available.
	if objInfo.ETag != "" {
		w.Header()[consts.ETag] = []string{"\"" + objInfo.ETag + "\""}
	}

	if objInfo.ContentType != "" {
		w.Header().Set(consts.ContentType, objInfo.ContentType)
	}

	if objInfo.ContentEncoding != "" {
		w.Header().Set(consts.ContentEncoding, objInfo.ContentEncoding)
	}

	if !objInfo.Expires.IsZero() {
		w.Header().Set(consts.Expires, objInfo.Expires.UTC().Format(http.TimeFormat))
	}

	// Set content length
	w.Header().Set(consts.ContentLength, strconv.FormatInt(objInfo.Size, 10))

	// Set the relevant version ID as part of the response header.
	if objInfo.VersionID != "" {
		w.Header()[consts.AmzVersionID] = []string{objInfo.VersionID}
	}

}

// SetHeadGetRespHeaders - set any requested parameters as response headers.
func SetHeadGetRespHeaders(w http.ResponseWriter, reqParams url.Values) {
	for k, v := range reqParams {
		if header, ok := supportedHeadGetReqParams[strings.ToLower(k)]; ok {
			w.Header()[header] = v
		}
	}
}

// supportedHeadGetReqParams - supported request parameters for GET and HEAD presigned request.
var supportedHeadGetReqParams = map[string]string{
	"response-expires":             consts.Expires,
	"response-content-type":        consts.ContentType,
	"response-content-encoding":    consts.ContentEncoding,
	"response-content-language":    consts.ContentLanguage,
	"response-content-disposition": consts.ContentDisposition,
}
