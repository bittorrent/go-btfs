package sessions

import (
	"context"
	"strings"

	"github.com/bittorrent/protobuf/proto"

	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

func Batch(d ds.Datastore, keys []string, vals []proto.Message) error {
	ctx := context.TODO()
	batch := ds.NewBasicBatch(d)
	for i, k := range keys {
		if vals[i] == nil {
			batch.Delete(ctx, ds.NewKey(k))
			continue
		}
		bytes, err := proto.Marshal(vals[i])
		if err != nil {
			return err
		}
		batch.Put(ctx, ds.NewKey(k), bytes)
	}
	return batch.Commit(ctx)
}

func Save(d ds.Datastore, key string, val proto.Message) error {
	ctx := context.TODO()
	bytes, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	return d.Put(ctx, ds.NewKey(key), bytes)
}

func Get(d ds.Datastore, key string, m proto.Message) error {
	ctx := context.TODO()
	bytes, err := d.Get(ctx, ds.NewKey(key))
	if err != nil {
		return err
	}
	return proto.Unmarshal(bytes, m)
}

func Remove(d ds.Datastore, key string) error {
	ctx := context.TODO()
	return d.Delete(ctx, ds.NewKey(key))
}

func List(d ds.Datastore, prefix string, substrInKey ...string) ([][]byte, error) {
	vs := make([][]byte, 0)
	ctx := context.TODO()
	results, err := d.Query(ctx, query.Query{
		Prefix:  prefix,
		Filters: []query.Filter{},
	})
	if err != nil {
		return nil, err
	}
	for entry := range results.Next() {
		contains := true
		for _, substr := range substrInKey {
			contains = contains && strings.Contains(entry.Key, substr)
		}
		if contains {
			value := entry.Value
			vs = append(vs, value)
		}
	}
	return vs, nil
}

func ListKeys(d ds.Datastore, prefix string, substrInKey ...string) ([]string, error) {
	ks := make([]string, 0)
	ctx := context.TODO()
	results, err := d.Query(ctx, query.Query{
		Prefix:   prefix,
		Filters:  []query.Filter{},
		KeysOnly: true,
	})
	if err != nil {
		return nil, err
	}
	for entry := range results.Next() {
		contains := true
		for _, substr := range substrInKey {
			contains = contains && strings.Contains(entry.Key, substr)
		}
		if contains {
			ks = append(ks, entry.Key)
		}
	}
	return ks, nil
}
