package main

import (
	"fmt"
	"testing"

	"gotest.tools/assert"

	cmds "github.com/bittorrent/go-btfs-cmds"
	Cfg "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
)

func TestGetChainID(t *testing.T) {
	//chain init
	statestore, err := chain.InitStateStore("~/.btfs.tmp.not.stored")
	if err != nil {
		t.Errorf("init InitStateStore err: %v", err)
	}

	// 1.not store chainid, first input chainid, second default chainid
	req := &cmds.Request{}
	cfg := &Cfg.Config{}
	cfg.ChainInfo.ChainId = int64(199) // default chainid
	chainid, stored, err := getChainID(req, cfg, statestore)
	if err != nil {
		t.Errorf("1 init getChainID err: %v", err)
	}
	assert.Equal(t, chainid, cfg.ChainInfo.ChainId, "not stored default chainid")
	assert.Equal(t, stored, false, "not stored")

	req2 := &cmds.Request{
		Options: make(cmds.OptMap),
	}
	req2.Options[chainID] = "199" // input chainid
	chainid, stored, err = getChainID(req2, cfg, statestore)
	if err != nil {
		t.Errorf("2 init getChainID err: %v", err)
	}
	assert.Equal(t, chainid, int64(199), "not stored input chainid")
	assert.Equal(t, stored, false, "not stored")

	// 2.write to leveldb chainid, getChainID must be storedChainId
	storedChainId := int64(199) // stored chainid
	statestore2, err := chain.InitStateStore("~/.btfs.tmp.stored")
	if err != nil {
		t.Errorf("init InitStateStore err: %v", err)
	}
	err = chain.StoreChainIdToDisk(storedChainId, statestore2)
	if err != nil {
		t.Errorf("init StoreChainId err: %v", err)
	}

	req3 := &cmds.Request{
		Options: make(cmds.OptMap),
	}
	req3.Options[chainID] = "1029" // input chainid
	cfg3 := &Cfg.Config{}
	cfg3.ChainInfo.ChainId = int64(199)
	chainid, stored, err = getChainID(req3, cfg3, statestore2)
	if err != nil {
		fmt.Println("3 init getChainID warn: ", err)
	} else {
		assert.Equal(t, stored, true, "stored config wrong")
	}

	req4 := &cmds.Request{
		Options: make(cmds.OptMap),
	}
	req4.Options[chainID] = "199" // input chainid
	cfg4 := &Cfg.Config{}
	cfg4.ChainInfo.ChainId = int64(1029) // default chainid
	chainid, stored, err = getChainID(req4, cfg4, statestore2)
	if err != nil {
		fmt.Println("4 init getChainID warn: ", err)
	} else {
		assert.Equal(t, stored, true, "stored input wrong")
	}

	req5 := &cmds.Request{
		Options: make(cmds.OptMap),
	}
	req5.Options[chainID] = "199" // input chainid
	cfg5 := &Cfg.Config{}
	cfg5.ChainInfo.ChainId = int64(199) // default chainid
	chainid, stored, err = getChainID(req5, cfg5, statestore2)
	if err != nil {
		t.Errorf("5 init getChainID err: %v", err)
	}
	assert.Equal(t, chainid, int64(199), "stored all chainid")
	assert.Equal(t, stored, true, "stored all stored")
}
