package responses

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
	"path"
)

func owner(accessKey string) *s3.Owner {
	return new(s3.Owner).SetID(accessKey).SetDisplayName(accessKey)
}

func ownerFullControlGrant(accessKey string) *s3.Grant {
	return new(s3.Grant).SetGrantee(new(s3.Grantee).SetType(s3.TypeCanonicalUser).SetID(accessKey).SetDisplayName(accessKey)).SetPermission(s3.PermissionFullControl)
}

var (
	allUsersReadGrant  = new(s3.Grant).SetGrantee(new(s3.Grantee).SetType(s3.TypeGroup).SetURI(consts.AllUsersURI)).SetPermission(s3.PermissionRead)
	allUsersWriteGrant = new(s3.Grant).SetGrantee(new(s3.Grantee).SetType(s3.TypeGroup).SetURI(consts.AllUsersURI)).SetPermission(s3.PermissionWrite)
)

type ErrorOutput struct {
	_         struct{} `type:"structure"`
	Code      string   `locationName:"Code" type:"string"`
	Message   string   `locationName:"Message" type:"string"`
	Resource  string   `locationName:"Resource" type:"string"`
	RequestID string   `locationName:"RequestID" type:"string"`
}

func NewErrOutput(r *http.Request, rerr *Error) *ErrorOutput {
	return &ErrorOutput{
		Code:      rerr.Code(),
		Message:   rerr.Description(),
		Resource:  pathClean(r.URL.Path),
		RequestID: "", // this field value will be automatically filled
	}
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, rerr *Error) {
	output := NewErrOutput(r, rerr)
	_ = WriteResponse(w, rerr.HTTPStatusCode(), output, "Error")
}

func WriteSuccessResponse(w http.ResponseWriter, output interface{}, locationName string) {
	_ = WriteResponse(w, http.StatusOK, output, locationName)
}

func setPutObjHeaders(w http.ResponseWriter, etag, cid string, delete bool) {
	if etag != "" && !delete {
		w.Header()[consts.ETag] = []string{`"` + etag + `"`}
	}
	if cid != "" {
		w.Header()[consts.CID] = []string{cid}
	}
}

func pathClean(p string) string {
	cp := path.Clean(p)
	if cp == "." {
		return ""
	}
	return cp
}
