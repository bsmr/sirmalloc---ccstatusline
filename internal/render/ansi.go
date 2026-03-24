package render

import "fmt"

// Reset returns the ANSI reset escape sequence.
func Reset() string {
	return "\033[0m"
}

// Bold returns the ANSI bold escape sequence.
func Bold() string {
	return "\033[1m"
}

// FG returns an ANSI foreground color escape sequence.
// Accepted formats:
//   - ""              → empty string (no-op)
//   - "red", "green"  → basic named color (30-37 range)
//   - "38;5;N"        → 256-color
//   - "38;2;R;G;B"    → truecolor
func FG(color string) string {
	if color == "" {
		return ""
	}
	switch color {
	case "black":
		return "\033[30m"
	case "red":
		return "\033[31m"
	case "green":
		return "\033[32m"
	case "yellow":
		return "\033[33m"
	case "blue":
		return "\033[34m"
	case "magenta":
		return "\033[35m"
	case "cyan":
		return "\033[36m"
	case "white":
		return "\033[37m"
	}
	// 256-color or truecolor passthrough
	return fmt.Sprintf("\033[%sm", color)
}

// BG returns an ANSI background color escape sequence.
// Accepted formats mirror FG but for background (40-47, 48;5;N, 48;2;R;G;B).
func BG(color string) string {
	if color == "" {
		return ""
	}
	switch color {
	case "black":
		return "\033[40m"
	case "red":
		return "\033[41m"
	case "green":
		return "\033[42m"
	case "yellow":
		return "\033[43m"
	case "blue":
		return "\033[44m"
	case "magenta":
		return "\033[45m"
	case "cyan":
		return "\033[46m"
	case "white":
		return "\033[47m"
	}
	return fmt.Sprintf("\033[%sm", color)
}
