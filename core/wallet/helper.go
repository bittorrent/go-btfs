package wallet

import (
	"strings"

	config "github.com/TRON-US/go-btfs-config"
)

func getTokenId(cfg *config.Config) string {
	tokenId := TokenId
	if strings.Contains(cfg.Services.EscrowDomain, "dev") ||
		strings.Contains(cfg.Services.EscrowDomain, "staging") {
		tokenId = TokenIdDev
	}
	return tokenId
}
