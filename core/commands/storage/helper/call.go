package helper

import (
	"context"
	"fmt"
	"strings"

	shell "github.com/bittorrent/go-btfs-api"
	"github.com/bittorrent/go-btfs-config"
)

func Call(ctx context.Context, cfg *config.Config, sub string) error {
	url := fmt.Sprint(strings.Split(cfg.Addresses.API[0], "/")[2], ":", strings.Split(cfg.Addresses.API[0], "/")[4])
	resp, err := shell.NewShell(url).Request(sub).Send(ctx)
	defer resp.Close()
	if err != nil {
		return err
	}
	return resp.Error
}
