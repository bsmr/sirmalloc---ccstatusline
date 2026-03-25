package render

import "go.a8l.eu/ccstatusline/internal/config"

// segment holds the rendered text and color attributes for a single widget slot.
type segment struct {
	text string
	fg   string
	bg   string
	bold bool
}

// renderPowerline builds a single output line using Powerline-style arrow separators.
//
// Layout:
//   - Optional start cap in the color of the first segment's BG.
//   - Each segment rendered as BG(seg.bg) + FG(seg.fg) + text.
//   - Between adjacent segments: Reset + FG(prevBG) + BG(nextBG) + separator.
//     If pw.SeparatorInvertBackground[i] is true for that separator index,
//     the FG/BG of the separator itself are swapped.
//   - Optional end cap in the color of the last segment's BG.
//   - Reset appended at the very end.
func renderPowerline(segments []segment, pw *config.PowerlineConfig) string {
	if len(segments) == 0 {
		return ""
	}

	sep := "\ue0b0" // nerd-font right arrow default
	if len(pw.Separators) > 0 && pw.Separators[0] != "" {
		sep = pw.Separators[0]
	}

	var out []byte

	// Start cap: drawn in the foreground color of the first segment's BG.
	if len(pw.StartCaps) > 0 && pw.StartCaps[0] != "" {
		out = append(out, FG(segments[0].bg)...)
		out = append(out, pw.StartCaps[0]...)
	}

	for i, seg := range segments {
		// Segment body.
		out = append(out, BG(seg.bg)...)
		out = append(out, FG(seg.fg)...)
		if seg.bold {
			out = append(out, Bold()...)
		}
		out = append(out, seg.text...)

		// Separator between this segment and the next.
		if i < len(segments)-1 {
			next := segments[i+1]

			// Determine whether to invert FG/BG for this separator.
			invert := false
			if i < len(pw.SeparatorInvertBackground) {
				invert = pw.SeparatorInvertBackground[i]
			}

			out = append(out, Reset()...)
			if invert {
				out = append(out, FG(next.bg)...)
				out = append(out, BG(seg.bg)...)
			} else {
				out = append(out, FG(seg.bg)...)
				out = append(out, BG(next.bg)...)
			}
			out = append(out, sep...)
		}
	}

	// End cap: drawn in the foreground color of the last segment's BG,
	// after resetting all attributes.
	if len(pw.EndCaps) > 0 && pw.EndCaps[0] != "" {
		last := segments[len(segments)-1]
		out = append(out, Reset()...)
		out = append(out, FG(last.bg)...)
		out = append(out, pw.EndCaps[0]...)
	}

	out = append(out, Reset()...)
	return string(out)
}
