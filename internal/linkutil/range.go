package linkutil

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRange parses range=[start,end] as a half-open interval [start, end).
func ParseRange(raw string) (start, end int, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, 0, fmt.Errorf("empty range")
	}
	if !strings.HasPrefix(raw, "[") || !strings.HasSuffix(raw, "]") {
		return 0, 0, fmt.Errorf("invalid range format")
	}
	inner := strings.TrimSpace(raw[1 : len(raw)-1])
	parts := strings.Split(inner, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}
	start, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || start < 0 {
		return 0, 0, fmt.Errorf("invalid range start")
	}
	end, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || end < start {
		return 0, 0, fmt.Errorf("invalid range end")
	}
	return start, end, nil
}

func ContentRange(start, end int, total int64) string {
	return fmt.Sprintf("links %d-%d/%d", start, end, total)
}
