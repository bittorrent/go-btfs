package object

import (
	"context"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
	"strings"
	"time"
)

const (
	objectKeyFormat = "obj/%s/%s"
)

var _ Service = (*service)(nil)

// service captures all bucket metadata for a given cluster.
type service struct {
	providers providers.Providerser
}

// NewService - creates new policy system.
func NewService(providers providers.Providerser, options ...Option) Service {
	s := &service{
		providers: providers,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s *service) getObjectKey(buc, obj string) string {
	return fmt.Sprintf(objectKeyFormat, buc, obj)
}

func (s *service) StoreObject(ctx context.Context, bucname, objname string, reader *hash.Reader, size int64, meta map[string]string) (obj Object, err error) {
	cid, err := s.providers.GetFileStore().AddWithOpts(reader, true, true)
	if err != nil {
		return
	}

	obj = Object{
		Bucket:           bucname,
		Name:             objname,
		ModTime:          time.Now().UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             reader.ETag().String(),
		Cid:              cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		Acl:              meta[consts.AmzACL],
		ContentType:      meta[strings.ToLower(consts.ContentType)],
		ContentEncoding:  meta[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: time.Now().UTC(),
	}

	// Update expires
	if exp, ok := meta[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			obj.Expires = t.UTC()
		}
	}

	err = s.providers.GetStateStore().Put(s.getObjectKey(bucname, objname), obj)
	if err != nil {
		return
	}

	return
}
