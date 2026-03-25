package render

import (
	"fmt"
	"strconv"
	"strings"
)

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
//   - ""                   → empty string (no-op)
//   - "red", "green" …     → basic named color (30–37)
//   - "brightRed" …        → bright named color (90–97)
//   - "hex:RRGGBB"         → truecolor \033[38;2;R;G;Bm
//   - "ansi256:N"          → 256-color \033[38;5;Nm
//   - "38;5;N" / "38;2;…"  → raw passthrough (legacy)
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
	// Bright foreground colors (90–97).
	case "brightBlack":
		return "\033[90m"
	case "brightRed":
		return "\033[91m"
	case "brightGreen":
		return "\033[92m"
	case "brightYellow":
		return "\033[93m"
	case "brightBlue":
		return "\033[94m"
	case "brightMagenta":
		return "\033[95m"
	case "brightCyan":
		return "\033[96m"
	case "brightWhite":
		return "\033[97m"
	}
	// hex:RRGGBB → truecolor foreground.
	if strings.HasPrefix(color, "hex:") {
		hex := color[4:]
		if len(hex) == 6 {
			r, _ := strconv.ParseUint(hex[0:2], 16, 8)
			g, _ := strconv.ParseUint(hex[2:4], 16, 8)
			b, _ := strconv.ParseUint(hex[4:6], 16, 8)
			return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
		}
	}
	// ansi256:N → 256-color foreground.
	if strings.HasPrefix(color, "ansi256:") {
		n, err := strconv.Atoi(color[8:])
		if err == nil && n >= 0 && n <= 255 {
			return fmt.Sprintf("\033[38;5;%dm", n)
		}
	}
	// Legacy raw passthrough (e.g. "38;5;N").
	return fmt.Sprintf("\033[%sm", color)
}

// BG returns an ANSI background color escape sequence.
// Accepted formats:
//   - ""                    → empty string (no-op)
//   - "red", "green" …      → basic named color (40–47)
//   - "bgRed", "bgGreen" …  → TypeScript-style bg prefix (40–47)
//   - "bgBrightRed" …       → bright bg (100–107)
//   - "hex:RRGGBB"          → truecolor \033[48;2;R;G;Bm
//   - "ansi256:N"           → 256-color \033[48;5;Nm
//   - "48;5;N" / "48;2;…"   → raw passthrough (legacy)
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
	// TypeScript-style "bg"-prefixed named colors (40–47).
	case "bgBlack":
		return "\033[40m"
	case "bgRed":
		return "\033[41m"
	case "bgGreen":
		return "\033[42m"
	case "bgYellow":
		return "\033[43m"
	case "bgBlue":
		return "\033[44m"
	case "bgMagenta":
		return "\033[45m"
	case "bgCyan":
		return "\033[46m"
	case "bgWhite":
		return "\033[47m"
	// Bright background colors (100–107).
	case "bgBrightBlack":
		return "\033[100m"
	case "bgBrightRed":
		return "\033[101m"
	case "bgBrightGreen":
		return "\033[102m"
	case "bgBrightYellow":
		return "\033[103m"
	case "bgBrightBlue":
		return "\033[104m"
	case "bgBrightMagenta":
		return "\033[105m"
	case "bgBrightCyan":
		return "\033[106m"
	case "bgBrightWhite":
		return "\033[107m"
	}
	// hex:RRGGBB → truecolor background.
	if strings.HasPrefix(color, "hex:") {
		hex := color[4:]
		if len(hex) == 6 {
			r, _ := strconv.ParseUint(hex[0:2], 16, 8)
			g, _ := strconv.ParseUint(hex[2:4], 16, 8)
			b, _ := strconv.ParseUint(hex[4:6], 16, 8)
			return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
		}
	}
	// ansi256:N → 256-color background.
	if strings.HasPrefix(color, "ansi256:") {
		n, err := strconv.Atoi(color[8:])
		if err == nil && n >= 0 && n <= 255 {
			return fmt.Sprintf("\033[48;5;%dm", n)
		}
	}
	// Legacy raw passthrough (e.g. "48;5;N").
	return fmt.Sprintf("\033[%sm", color)
}

// BGColorToFG converts a background color string to its foreground equivalent.
// This is used in powerline separators where the previous segment's BG color
// becomes the separator's FG color.
//
//   - "bgRed"         → "red"
//   - "bgBrightCyan"  → "brightCyan"
//   - "hex:..."       → "hex:..." (pass through unchanged)
//   - "ansi256:..."   → "ansi256:..." (pass through unchanged)
func BGColorToFG(bg string) string {
	if strings.HasPrefix(bg, "bgBright") {
		name := bg[8:] // e.g. "Red"
		return "bright" + name // e.g. "brightRed"
	}
	if strings.HasPrefix(bg, "bg") {
		name := bg[2:] // e.g. "Red"
		return strings.ToLower(name[:1]) + name[1:] // e.g. "red"
	}
	// hex:, ansi256:, plain named colors — pass through unchanged.
	return bg
}
