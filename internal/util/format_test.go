package util_test

import (
	"testing"

	"go.a8l.eu/ccstatusline/internal/util"
)

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1.0k"},
		{1500, "1.5k"},
		{15200, "15.2k"},
		{999999, "1000.0k"},
		{1000000, "1.0M"},
		{1500000, "1.5M"},
	}
	for _, tc := range tests {
		got := util.FormatTokens(tc.n)
		if got != tc.want {
			t.Errorf("FormatTokens(%d) = %q, want %q", tc.n, got, tc.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms   int64
		want string
	}{
		{0, "0m"},
		{59999, "0m"},
		{60000, "1m"},
		{45000, "0m"},
		{2700000, "45m"},
		{3600000, "1hr 0m"},
		{3660000, "1hr 1m"},
		{8100000, "2hr 15m"},
	}
	for _, tc := range tests {
		got := util.FormatDuration(tc.ms)
		if got != tc.want {
			t.Errorf("FormatDuration(%d) = %q, want %q", tc.ms, got, tc.want)
		}
	}
}
