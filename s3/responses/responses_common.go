package responses

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
	"path"
	"time"
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

func getRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func setCommonHeader(w http.ResponseWriter, requestId string) {
	w.Header().Set(consts.ServerInfo, consts.DefaultServerInfo)
	w.Header().Set(consts.AmzRequestID, requestId)
	w.Header().Set(consts.AcceptRanges, "bytes")
}

type ErrorOutput struct {
	_         struct{} `type:"structure"`
	Code      string   `locationName:"Code" type:"string"`
	Message   string   `locationName:"Message" type:"string"`
	Resource  string   `locationName:"Resource" type:"string"`
	RequestID string   `locationName:"RequestID" type:"string"`
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, rerr *Error) {
	reqID := getRequestID()
	setCommonHeader(w, reqID)
	output := &ErrorOutput{
		Code:      rerr.Code(),
		Message:   rerr.Description(),
		Resource:  pathClean(r.URL.Path),
		RequestID: reqID,
	}
	err := WriteResponse(w, rerr.HTTPStatusCode(), output, "Error")
	if err != nil {
		fmt.Println("write response: ", err)
	}
}

func WriteSuccessResponse(w http.ResponseWriter, output interface{}, locationName string) {
	setCommonHeader(w, getRequestID())
	err := WriteResponse(w, http.StatusOK, output, locationName)
	if err != nil {
		fmt.Println("write response: ", err)
	}
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
