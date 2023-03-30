package helper

import (
	"encoding/base64"
	"fmt"

	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func ConvertPubKeyFromString(pubKeyStr string) (ic.PubKey, error) {
	raw, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	return ic.UnmarshalPublicKey(raw)
}

func ConvertToPubKey(pubKeyStr string) (ic.PubKey, error) {
	pubKey, err := ConvertPubKeyFromString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func PidFromString(key string) (peer.ID, error) {
	pubKey, err := ConvertPubKeyFromString(key)
	if err != nil {
		return "", err
	}
	return peer.IDFromPublicKey(pubKey)
}

const (
	Text = iota + 1
	Base64
)

func BytesToString(data []byte, encoding int) (string, error) {
	switch encoding {
	case Text:
		return string(data), nil
	case Base64:
		return base64.StdEncoding.EncodeToString(data), nil
	default:
		return "", fmt.Errorf(`unexpected parameter [%d] is given, either "text" or "base64" should be used`, encoding)
	}
}

func StringToBytes(str string, encoding int) ([]byte, error) {
	switch encoding {
	case Text:
		return []byte(str), nil
	case Base64:
		by, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, err
		}
		return by, nil
	default:
		return nil, fmt.Errorf(`unexpected encoding [%d], expected 1(Text) or 2(Base64)`, encoding)
	}
}
