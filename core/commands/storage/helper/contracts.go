package helper

import guardpb "github.com/bittorrent/go-btfs-common/protos/guard"

var ContractFilterMap = map[string]map[guardpb.Contract_ContractState]bool{
	"active": {
		guardpb.Contract_UPLOADED: true,
		guardpb.Contract_RENEWED:  true,
		guardpb.Contract_WARN:     true,
	},
	"finished": {
		guardpb.Contract_CLOSED: true,
	},
	"invalid": {
		guardpb.Contract_DRAFT:             true,
		guardpb.Contract_SIGNED:            true,
		guardpb.Contract_LOST:              true,
		guardpb.Contract_CANCELED:          true,
		guardpb.Contract_OBSOLETE:          true,
		guardpb.Contract_DRAFT_WARN:        true,
		guardpb.Contract_READY_CHALLENGE:   true,
		guardpb.Contract_RECREATED:         true,
		guardpb.Contract_REQUEST_CHALLENGE: true,
	},
	"all": {
		guardpb.Contract_DRAFT:             true,
		guardpb.Contract_SIGNED:            true,
		guardpb.Contract_UPLOADED:          true,
		guardpb.Contract_LOST:              true,
		guardpb.Contract_CANCELED:          true,
		guardpb.Contract_CLOSED:            true,
		guardpb.Contract_RENEWED:           true,
		guardpb.Contract_OBSOLETE:          true,
		guardpb.Contract_WARN:              true,
		guardpb.Contract_DRAFT_WARN:        true,
		guardpb.Contract_READY_CHALLENGE:   true,
		guardpb.Contract_RECREATED:         true,
		guardpb.Contract_REQUEST_CHALLENGE: true,
	},
}
