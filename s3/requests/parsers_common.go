package requests

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/utils"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"path"
)

func parseBucket(r *http.Request) (bucket string, rerr *responses.Error) {
	bucket = mux.Vars(r)["bucket"]
	err := s3utils.CheckValidBucketNameStrict(bucket)
	if err != nil {
		rerr = responses.ErrInvalidBucketName
	}
	return
}

func parseObject(r *http.Request) (object string, rerr *responses.Error) {
	object, err := unescapePath(mux.Vars(r)["object"])
	if err != nil {
		rerr = responses.ErrInvalidRequestParameter
	}
	return
}

// Parses location constraint from the incoming reader.
func parseLocationConstraint(r *http.Request) (location string, rerr *responses.Error) {
	// If the request has no body with content-length set to 0,
	// we do not have to validate location constraint. Bucket will
	// be created at default region.
	locationConstraint := createBucketLocationConfiguration{}
	err := utils.XmlDecoder(r.Body, &locationConstraint, r.ContentLength)
	if err != nil && r.ContentLength != 0 {
		rerr = responses.ErrMalformedXML
		return
	} // else for both err as nil or io.EOF

	location = locationConstraint.Location
	if location == "" {
		location = consts.DefaultRegion
	}

	return
}

var supportAcls = map[string]struct{}{
	policy.Private:         {},
	policy.PublicRead:      {},
	policy.PublicReadWrite: {},
}

func parseAcl(r *http.Request) (acl string, rerr *responses.Error) {
	acl = r.Header.Get(consts.AmzACL)
	if acl == "" {
		acl = consts.DefaultAcl
	}
	_, ok := supportAcls[acl]
	if !ok {
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
