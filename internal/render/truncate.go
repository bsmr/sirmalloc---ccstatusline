package render

import "unicode/utf8"

// VisibleLen returns the visible (printable) character count of s,
// stripping ANSI CSI sequences (\033[...m) and OSC sequences
// (\033]...\007 or \033]...\033\\).
func VisibleLen(s string) int {
	count := 0
	i := 0
	for i < len(s) {
		if s[i] == '\033' && i+1 < len(s) {
			switch s[i+1] {
			case '[': // CSI sequence: \033[ ... <final byte in @-~>
				i += 2
				for i < len(s) {
					b := s[i]
					i++
					if b >= 0x40 && b <= 0x7e { // final byte
						break
					}
				}
				continue
			case ']': // OSC sequence: \033] ... \007  or  \033] ... \033\\
				i += 2
				for i < len(s) {
					if s[i] == '\007' {
						i++
						break
					}
					if s[i] == '\033' && i+1 < len(s) && s[i+1] == '\\' {
						i += 2
						break
					}
					i++
				}
				continue
			}
		}
		// Ordinary rune — count it once regardless of byte width.
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		count++
	}
	return count
}

// Truncate shortens s to at most maxVisible visible characters,
// appending "..." (preceded by a Reset) if truncation occurred.
// ANSI escape sequences are copied verbatim and do not count toward the limit.
func Truncate(s string, maxVisible int) string {
	if maxVisible <= 0 {
		return Reset() + "..."
	}
	if VisibleLen(s) <= maxVisible {
		return s
	}

	// We need to truncate. Walk rune-by-rune collecting bytes until we hit
	// maxVisible - 3 visible chars (to leave room for "...").
	const ellipsis = "..."
	limit := maxVisible - len(ellipsis)
	if limit < 0 {
		limit = 0
	}

	var buf []byte
	visible := 0
	i := 0
	for i < len(s) && visible < limit {
		if s[i] == '\033' && i+1 < len(s) {
			switch s[i+1] {
			case '[': // CSI
				start := i
				i += 2
				for i < len(s) {
					b := s[i]
					i++
					if b >= 0x40 && b <= 0x7e {
						break
					}
				}
				buf = append(buf, s[start:i]...)
				continue
			case ']': // OSC
				start := i
				i += 2
				for i < len(s) {
					if s[i] == '\007' {
						i++
						break
					}
					if s[i] == '\033' && i+1 < len(s) && s[i+1] == '\\' {
						i += 2
						break
					}
					i++
				}
				buf = append(buf, s[start:i]...)
				continue
			}
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		buf = append(buf, s[i:i+size]...)
		i += size
		visible++
	}

	buf = append(buf, Reset()...)
	buf = append(buf, ellipsis...)
	return string(buf)
}
