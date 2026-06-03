package linkutil

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

const shortNameChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func ShortURL(baseURL, shortName string) string {
	return strings.TrimRight(baseURL, "/") + "/r/" + shortName
}

func ToResponse(linkID int64, originalURL, shortName, baseURL string) map[string]any {
	return map[string]any{
		"id":           linkID,
		"original_url": originalURL,
		"short_name":   shortName,
		"short_url":    ShortURL(baseURL, shortName),
	}
}

func RandomShortName(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = shortNameChars[int(b[i])%len(shortNameChars)]
	}
	return string(b), nil
}

func NormalizeBaseURL(baseURL string) string {
	return strings.TrimRight(strings.TrimSpace(baseURL), "/")
}

func ParseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}
