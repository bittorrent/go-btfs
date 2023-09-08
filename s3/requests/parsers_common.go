package requests

import (
	"encoding/xml"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"path"
)

func ParseBucketAndObject(r *http.Request) (bucket string, object string, rerr *responses.Error) {
	bucket, rerr = ParseBucket(r)
	if rerr != nil {
		return
	}
	object, rerr = ParseObject(r)
	return
}

func ParseBucket(r *http.Request) (bucket string, rerr *responses.Error) {
	bucket = mux.Vars(r)["Bucket"]
	err := s3utils.CheckValidBucketNameStrict(bucket)
	if err != nil {
		rerr = responses.ErrInvalidBucketName
	}
	return
}

func ParseObject(r *http.Request) (object string, rerr *responses.Error) {
	object, err := unescapePath(mux.Vars(r)["Key"])
	if err != nil {
		rerr = responses.ErrInvalidRequestParameter
	}
	return
}

func ParseLocation(r *http.Request) (location string, rerr *responses.Error) {
	if r.ContentLength != 0 {
		locationCfg := s3.CreateBucketConfiguration{}
		decoder := xml.NewDecoder(r.Body)
		err := xmlutil.UnmarshalXML(&locationCfg, decoder, "")
		if err != nil {
			rerr = responses.ErrMalformedXML
			return
		}
		location = *locationCfg.LocationConstraint
	}
	if len(location) == 0 {
		location = consts.DefaultLocation
	}
	if !consts.SupportedLocations[location] {
		rerr = responses.ErrNotImplemented
	}

	return
}

func ParseBucketACL(r *http.Request) (acl string, rerr *responses.Error) {
	acl = r.Header.Get(consts.AmzACL)
	if len(acl) == 0 {
		acl = consts.DefaultBucketACL
	}
	if !consts.SupportedBucketACLs[acl] {
		rerr = responses.ErrNotImplemented
	}
	return
}

func ParseObjectACL(r *http.Request) (acl string, rerr *responses.Error) {
	acl = r.Header.Get(consts.AmzACL)
	if len(acl) == 0 {
		acl = consts.DefaultObjectACL
	}
	if !consts.SupportedObjectACLs[acl] {
		rerr = responses.ErrNotImplemented
	}
	return
}

// unescapePath is similar to url.PathUnescape or url.QueryUnescape
// depending on input, additionally also handles situations such as
// `//` are normalized as `/`, also removes any `/` prefix before
// returning.
func unescapePath(p string) (string, error) {
	ep, err := url.PathUnescape(p)
	if err != nil {
		return "", err
	}
	return trimLeadingSlash(ep), nil
}

func trimLeadingSlash(ep string) string {
	if len(ep) > 0 && ep[0] == '/' {
		// Path ends with '/' preserve it
		if ep[len(ep)-1] == '/' && len(ep) > 1 {
			ep = path.Clean(ep)
			ep += "/"
		} else {
			ep = path.Clean(ep)
		}
		ep = ep[1:]
	}
	return ep
}
