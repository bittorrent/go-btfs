package handlers

import (
	"encoding/xml"
	"github.com/bittorrent/go-btfs/s3/apierrors"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/bittorrent/go-btfs/s3/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"path"
)

type RequestBinder interface {
	Bind(r *http.Request) (err error)
}

//type PutObjectRequest struct {
//	Bucket string
//	Object string
//	Body   io.Reader
//}
//
//func (req *PutObjectRequest) Bind(r *http.Request) (err error) {
//	return
//}

func checkPermissionType(s string) bool {
	switch s {
	case policy.PublicRead:
		return true
	case policy.PublicReadWrite:
		return true
	case policy.Private:
		return true
	}
	return false
}

type PutBucketRequest struct {
	Bucket string
	ACL    string
	Region string
	Body   io.Reader
}

func (req *PutBucketRequest) Bind(r *http.Request) (err error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	region, _ := parseLocationConstraint(r)

	acl := r.Header.Get(consts.AmzACL)

	//set request
	req.Bucket = bucket
	req.ACL = acl
	req.Region = region
	return
}

// Parses location constraint from the incoming reader.
func parseLocationConstraint(r *http.Request) (location string, s3Error apierrors.ErrorCode) {
	// If the request has no body with content-length set to 0,
	// we do not have to validate location constraint. Bucket will
	// be created at default region.
	locationConstraint := createBucketLocationConfiguration{}
	err := utils.XmlDecoder(r.Body, &locationConstraint, r.ContentLength)
	if err != nil && r.ContentLength != 0 {
		// Treat all other failures as XML parsing errors.
		return "", apierrors.ErrMalformedXML
	} // else for both err as nil or io.EOF
	location = locationConstraint.Location
	if location == "" {
		location = consts.DefaultRegion
	}
	return location, apierrors.ErrNone
}

// createBucketConfiguration container for bucket configuration request from client.
// Used for parsing the location from the request body for Makebucket.
type createBucketLocationConfiguration struct {
	XMLName  xml.Name `xml:"CreateBucketConfiguration" json:"-"`
	Location string   `xml:"LocationConstraint"`
}

// pathClean is like path.Clean but does not return "." for
// empty inputs, instead returns "empty" as is.
func pathClean(p string) string {
	cp := path.Clean(p)
	if cp == "." {
		return ""
	}
	return cp
}

func unmarshalXML(reader io.Reader, isObject bool) (*store.Tags, error) {
	tagging := &store.Tags{
		TagSet: &store.TagSet{
			TagMap:   make(map[string]string),
			IsObject: isObject,
		},
	}

	if err := xml.NewDecoder(reader).Decode(tagging); err != nil {
		return nil, err
	}

	return tagging, nil
}
