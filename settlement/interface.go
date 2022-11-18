package settlement

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	ErrPeerNoSettlements = errors.New("no settlements for peer")
)

// Interface is the interface used by Accounting to trigger settlement
type Interface interface {
	// TotalSent returns the total amount sent to a peer
	TotalSent(peer string) (totalSent *big.Int, err error)
	// TotalReceived returns the total amount received from a peer
	TotalReceived(peer string) (totalSent *big.Int, err error)
	// SettlementsSent returns sent settlements for each individual known peer
	SettlementsSent() (map[string]*big.Int, error)
	// SettlementsReceived returns received settlements for each individual known peer
	SettlementsReceived() (map[string]*big.Int, error)
}

type Accounting interface {
	Settle(peer string, amount *big.Int, contractId string, token common.Address) error
	NotifyPaymentReceived(peer string, amount *big.Int, token common.Address) error
	NotifyPaymentSent(peer string, amount *big.Int, receivedError error, token common.Address)
}
