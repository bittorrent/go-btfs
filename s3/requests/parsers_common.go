package requests

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"path"
)

func ParseBucketAndObject(r *http.Request) (bucket, object string, err error) {
	vars := mux.Vars(r)
	bucket = vars["bucket"]
	object, err = unescapePath(vars["object"])
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
