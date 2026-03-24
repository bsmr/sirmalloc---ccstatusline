package util

import (
	"fmt"
)

func FormatTokens(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

func FormatDuration(ms int64) string {
	totalMin := ms / 60_000
	if totalMin >= 60 {
		h := totalMin / 60
		m := totalMin % 60
		return fmt.Sprintf("%dhr %dm", h, m)
	}
	return fmt.Sprintf("%dm", totalMin)
}
