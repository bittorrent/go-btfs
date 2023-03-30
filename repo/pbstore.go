package repo

import (
	"context"

	"github.com/tron-us/protobuf/proto"

	"github.com/ipfs/go-datastore"
)

func Get(d datastore.Datastore, k string, m proto.Message) (proto.Message, error) {
	ctx := context.TODO()
	v, err := d.Get(ctx, datastore.NewKey(k))
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(v, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Put(d datastore.Datastore, k string, v proto.Message) error {
	bytes, err := proto.Marshal(v)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	return d.Put(ctx, datastore.NewKey(k), bytes)
}
