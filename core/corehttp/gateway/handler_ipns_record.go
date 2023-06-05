package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	ipns_pb "github.com/bittorrent/go-btns/pb"
	ipath "github.com/bittorrent/interface-go-btfs-core/path"
	"github.com/cespare/xxhash"
	"github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-cid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (i *handler) serveIpnsRecord(ctx context.Context, w http.ResponseWriter, r *http.Request, contentPath ipath.Path, begin time.Time, logger *zap.SugaredLogger) bool {
	ctx, span := spanTrace(ctx, "Handler.ServeIPNSRecord", trace.WithAttributes(attribute.String("path", contentPath.String())))
	defer span.End()

	if contentPath.Namespace() != "btns" {
		err := fmt.Errorf("%s is not an BTNS link", contentPath.String())
		webError(w, err, http.StatusBadRequest)
		return false
	}

	key := contentPath.String()
	key = strings.TrimSuffix(key, "/")
	key = strings.TrimPrefix(key, "/btns/")
	if strings.Count(key, "/") != 0 {
		err := errors.New("cannot find btns key for subpath")
		webError(w, err, http.StatusBadRequest)
		return false
	}

	c, err := cid.Decode(key)
	if err != nil {
		webError(w, err, http.StatusBadRequest)
		return false
	}

	rawRecord, err := i.api.GetIPNSRecord(ctx, c)
	if err != nil {
		webError(w, err, http.StatusInternalServerError)
		return false
	}

	var record ipns_pb.IpnsEntry
	err = proto.Unmarshal(rawRecord, &record)
	if err != nil {
		webError(w, err, http.StatusInternalServerError)
		return false
	}

	// Set cache control headers based on the TTL set in the IPNS record. If the
	// TTL is not present, we use the Last-Modified tag. We are tracking IPNS
	// caching on: https://github.com/ipfs/kubo/issues/1818.
	// TODO: use addCacheControlHeaders once #1818 is fixed.
	recordEtag := strconv.FormatUint(xxhash.Sum64(rawRecord), 32)
	w.Header().Set("Etag", recordEtag)
	if record.Ttl != nil {
		seconds := int(time.Duration(*record.Ttl).Seconds())
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", seconds))
	} else {
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	}

	// Set Content-Disposition
	var name string
	if urlFilename := r.URL.Query().Get("filename"); urlFilename != "" {
		name = urlFilename
	} else {
		name = key + ".btns-record"
	}
	setContentDispositionHeader(w, name, "attachment")

	w.Header().Set("Content-Type", "application/vnd.ipfs.btns-record")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	_, err = w.Write(rawRecord)
	if err == nil {
		// Update metrics
		i.ipnsRecordGetMetric.WithLabelValues(contentPath.Namespace()).Observe(time.Since(begin).Seconds())
		return true
	}

	return false
}
