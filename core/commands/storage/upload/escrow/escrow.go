package escrow

import (
	"fmt"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"

	"github.com/bittorrent/go-btfs-common/crypto"
	escrowpb "github.com/bittorrent/go-btfs-common/protos/escrow"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/protobuf/proto"
)

func NewSignedContract(contract *escrowpb.EscrowContract) *escrowpb.SignedEscrowContract {
	return &escrowpb.SignedEscrowContract{
		Contract: contract,
	}
}

func VerifyEscrowRes(configuration *config.Config, message proto.Message, sig []byte) error {
	escrowPubkey, err := helper.ConvertPubKeyFromString(configuration.Services.EscrowPubKeys[0])
	if err != nil {
		return err
	}
	ok, err := crypto.Verify(escrowPubkey, message, sig)
	if err != nil || !ok {
		return fmt.Errorf("verify escrow failed %v", err)
	}
	return nil
}
