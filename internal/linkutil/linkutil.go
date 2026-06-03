package linkutil

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zept888/go-project-278/internal/store"
)

const shortNameChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func ShortURL(baseURL, shortName string) string {
	return strings.TrimRight(baseURL, "/") + "/r/" + shortName
}

func ToResponse(link store.Link, baseURL string) map[string]any {
	out := map[string]any{
		"id":           link.ID,
		"original_url": link.OriginalURL,
		"short_name":   link.ShortName,
		"short_url":    ShortURL(baseURL, link.ShortName),
	}
	if !link.CreatedAt.IsZero() {
		out["created_at"] = link.CreatedAt.UTC().Format(time.RFC3339)
	}
	return out
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
