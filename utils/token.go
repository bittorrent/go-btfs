package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Claims struct {
	PeerId string
	Expiry int64
}

func GenerateToken(peerId, secret string, expiryDuration time.Duration) (string, error) {
	expiryTime := time.Now().Add(expiryDuration).Unix()
	claims := Claims{
		PeerId: peerId,
		Expiry: expiryTime,
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := base64.URLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.URLEncoding.EncodeToString(claimsBytes)
	signature := computeHMACSHA256(header+"."+payload, []byte(secret))

	token := header + "." + payload + "." + signature
	return token, nil
}

func VerifyToken(token string, secret string) (*Claims, error) {
	parts := splitToken(token)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	header := parts[0]
	payload := parts[1]
	receivedSignature := parts[2]

	expectedSignature := computeHMACSHA256(header+"."+payload, []byte(secret))
	if receivedSignature != expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimsBytes, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	var claims Claims
	err = json.Unmarshal(claimsBytes, &claims)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() > claims.Expiry {
		return nil, fmt.Errorf("token has expired")
	}

	return &claims, nil
}

func computeHMACSHA256(message string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func splitToken(token string) []string {
	return splitString(token, '.')
}

func splitString(s string, sep rune) []string {
	var result []string
	var buffer []rune
	for _, char := range s {
		if char == sep {
			result = append(result, string(buffer))
			buffer = buffer[:0]
		} else {
			buffer = append(buffer, char)
		}
	}
	if len(buffer) > 0 {
		result = append(result, string(buffer))
	}
	return result
}
