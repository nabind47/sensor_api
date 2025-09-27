package util

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

func GenerateHash(clientID, clientSecret string) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	data := clientID + clientSecret

	h := sha256.New()
	h.Write([]byte(data))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash + ":" + timestamp
}

func ValidateHash(clientID string, clientSecret string, clientExpiry time.Duration, token string) bool {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return false
	}

	providedHash := parts[0]
	timestampStr := parts[1]

	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	now := time.Now()
	tokenTime := time.Unix(timestampInt, 0)

	if now.Sub(tokenTime) > clientExpiry {
		return false
	}

	data := clientID + clientSecret
	h := sha256.New()
	h.Write([]byte(data))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	return expectedHash == providedHash
}
