package config

// PowerlineThemeColors holds fg and bg color arrays for one color level.
type PowerlineThemeColors struct {
	FG []string
	BG []string
}

// PowerlineTheme defines a named powerline color theme.
type PowerlineTheme struct {
	Name        string
	Description string
	Basic       *PowerlineThemeColors // color level 0 (basic 16)
	Colors256   *PowerlineThemeColors // color level 1 (256-color)
	Truecolor   *PowerlineThemeColors // color level 2 (truecolor)
}

// DefaultPowerlineTheme is the theme used when no theme is specified.
const DefaultPowerlineTheme = "nord-aurora"

// PowerlineThemes contains all built-in themes keyed by their ID.
var PowerlineThemes = map[string]PowerlineTheme{
	"custom": {Name: "Custom", Description: "Uses individual widget background colors"},

	"nord": {
		Name: "Nord", Description: "Arctic, north-bluish color palette",
		Basic: &PowerlineThemeColors{
			FG: []string{"black", "brightWhite", "brightWhite", "black", "black"},
			BG: []string{"bgBrightCyan", "bgBrightBlack", "bgBlue", "bgBrightYellow", "bgBrightGreen"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:16", "ansi256:254", "ansi256:231", "ansi256:231", "ansi256:16"},
			BG: []string{"ansi256:73", "ansi256:239", "ansi256:25", "ansi256:96", "ansi256:152"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:2E3440", "hex:D8DEE9", "hex:FDF6E3", "hex:2E3440", "hex:2E3440"},
			BG: []string{"hex:88C0D0", "hex:4C566A", "hex:5E81AC", "hex:B48EAD", "hex:A3BE8C"},
		},
	},

	"nord-aurora": {
		Name: "Nord Aurora", Description: "Nord theme with aurora colors",
		Basic: &PowerlineThemeColors{
			FG: []string{"brightWhite", "black", "black", "black", "black"},
			BG: []string{"bgRed", "bgBrightYellow", "bgBrightBlue", "bgGreen", "bgBrightMagenta"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:231", "ansi256:16", "ansi256:231", "ansi256:16", "ansi256:16"},
			BG: []string{"ansi256:131", "ansi256:220", "ansi256:68", "ansi256:108", "ansi256:176"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:ECEFF4", "hex:2E3440", "hex:FDF6E3", "hex:2E3440", "hex:2E3440"},
			BG: []string{"hex:BF616A", "hex:EBCB8B", "hex:5E81AC", "hex:A3BE8C", "hex:B48EAD"},
		},
	},

	"monokai": {
		Name: "Monokai", Description: "Dark background with vibrant colors",
		Basic: &PowerlineThemeColors{
			FG: []string{"black", "brightWhite", "black", "white", "black"},
			BG: []string{"bgBrightGreen", "bgBrightBlack", "bgBrightYellow", "bgMagenta", "bgBrightCyan"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:235", "ansi256:255", "ansi256:235", "ansi256:16", "ansi256:235"},
			BG: []string{"ansi256:148", "ansi256:238", "ansi256:186", "ansi256:141", "ansi256:81"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:272822", "hex:F8F8F2", "hex:272822", "hex:272822", "hex:272822"},
			BG: []string{"hex:A6E22E", "hex:49483E", "hex:E6DB74", "hex:AE81FF", "hex:66D9EF"},
		},
	},

	"solarized": {
		Name: "Solarized", Description: "Precision colors for readability",
		Basic: &PowerlineThemeColors{
			FG: []string{"brightWhite", "black", "brightWhite", "black", "black"},
			BG: []string{"bgBlue", "bgBrightYellow", "bgBrightBlack", "bgCyan", "bgBrightWhite"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:231", "ansi256:234", "ansi256:254", "ansi256:16", "ansi256:234"},
			BG: []string{"ansi256:33", "ansi256:136", "ansi256:240", "ansi256:37", "ansi256:254"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:073642", "hex:073642", "hex:FDF6E3", "hex:073642", "hex:073642"},
			BG: []string{"hex:268BD2", "hex:B58900", "hex:586E75", "hex:2AA198", "hex:EEE8D5"},
		},
	},

	"minimal": {
		Name: "Minimal", Description: "Clean monochrome theme",
		Basic: &PowerlineThemeColors{
			FG: []string{"brightWhite", "black", "white", "black", "black"},
			BG: []string{"bgBrightBlack", "bgBrightWhite", "bgBlack", "bgWhite", "bgBrightWhite"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:255", "ansi256:232", "ansi256:255", "ansi256:232", "ansi256:252"},
			BG: []string{"ansi256:240", "ansi256:251", "ansi256:233", "ansi256:248", "ansi256:236"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:FFFFFF", "hex:1C1C1C", "hex:FFFFFF", "hex:1C1C1C", "hex:E4E4E4"},
			BG: []string{"hex:585858", "hex:D0D0D0", "hex:1A1A1A", "hex:A8A8A8", "hex:303030"},
		},
	},

	"dracula": {
		Name: "Dracula", Description: "Dark theme with purple accents",
		Basic: &PowerlineThemeColors{
			FG: []string{"brightWhite", "black", "brightWhite", "black", "white"},
			BG: []string{"bgMagenta", "bgBrightWhite", "bgRed", "bgBrightCyan", "bgBrightBlack"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:235", "ansi256:235", "ansi256:235", "ansi256:235", "ansi256:231"},
			BG: []string{"ansi256:141", "ansi256:253", "ansi256:204", "ansi256:117", "ansi256:236"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:282A36", "hex:282A36", "hex:282A36", "hex:282A36", "hex:F8F8F2"},
			BG: []string{"hex:BD93F9", "hex:F8F8F2", "hex:FF5555", "hex:8BE9FD", "hex:44475A"},
		},
	},

	"catppuccin": {
		Name: "Catppuccin", Description: "Soothing pastel theme",
		Basic: &PowerlineThemeColors{
			FG: []string{"black", "brightWhite", "black", "brightWhite", "black"},
			BG: []string{"bgBrightMagenta", "bgBrightBlack", "bgBrightGreen", "bgBlue", "bgBrightYellow"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:235", "ansi256:255", "ansi256:235", "ansi256:235", "ansi256:235"},
			BG: []string{"ansi256:176", "ansi256:238", "ansi256:150", "ansi256:210", "ansi256:111"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:1E1E2E", "hex:CDD6F4", "hex:1E1E2E", "hex:1E1E2E", "hex:CDD6F4"},
			BG: []string{"hex:CBA6F7", "hex:45475A", "hex:A6E3A1", "hex:F38BA8", "hex:585B70"},
		},
	},

	"gruvbox": {
		Name: "Gruvbox", Description: "Retro groove color scheme",
		Basic: &PowerlineThemeColors{
			FG: []string{"brightWhite", "black", "black", "brightWhite", "black"},
			BG: []string{"bgRed", "bgBrightYellow", "bgBrightWhite", "bgBlue", "bgBrightGreen"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:16", "ansi256:235", "ansi256:235", "ansi256:16", "ansi256:235"},
			BG: []string{"ansi256:167", "ansi256:214", "ansi256:246", "ansi256:109", "ansi256:142"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:EBDBB2", "hex:282828", "hex:282828", "hex:FDF6E3", "hex:282828"},
			BG: []string{"hex:CC241D", "hex:FABD2F", "hex:A89984", "hex:458588", "hex:98971A"},
		},
	},

	"onedark": {
		Name: "One Dark", Description: "Atom-inspired dark theme",
		Basic: &PowerlineThemeColors{
			FG: []string{"black", "brightWhite", "black", "brightWhite", "black"},
			BG: []string{"bgBrightBlue", "bgBrightBlack", "bgBrightGreen", "bgRed", "bgBrightYellow"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:235", "ansi256:251", "ansi256:235", "ansi256:16", "ansi256:235"},
			BG: []string{"ansi256:75", "ansi256:237", "ansi256:114", "ansi256:204", "ansi256:180"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:282C34", "hex:ABB2BF", "hex:282C34", "hex:282C34", "hex:282C34"},
			BG: []string{"hex:61AFEF", "hex:3E4452", "hex:98C379", "hex:E06C75", "hex:E5C07B"},
		},
	},

	"tokyonight": {
		Name: "Tokyo Night", Description: "Clean, dark theme inspired by Tokyo nightlife",
		Basic: &PowerlineThemeColors{
			FG: []string{"brightWhite", "black", "brightWhite", "black", "black"},
			BG: []string{"bgBlue", "bgBrightWhite", "bgMagenta", "bgBrightYellow", "bgBrightCyan"},
		},
		Colors256: &PowerlineThemeColors{
			FG: []string{"ansi256:16", "ansi256:234", "ansi256:16", "ansi256:234", "ansi256:234"},
			BG: []string{"ansi256:111", "ansi256:248", "ansi256:176", "ansi256:221", "ansi256:80"},
		},
		Truecolor: &PowerlineThemeColors{
			FG: []string{"hex:1A1B26", "hex:1A1B26", "hex:1A1B26", "hex:1A1B26", "hex:1A1B26"},
			BG: []string{"hex:7AA2F7", "hex:D5D6DB", "hex:BB9AF7", "hex:E0AF68", "hex:7DCFFF"},
		},
	},
}

// GetThemeColors returns the appropriate color set for a theme and color level.
// colorLevel: 0=basic, 1=256, 2=truecolor
// Returns nil if theme not found or has no colors for the requested level (falls back).
func GetThemeColors(themeName string, colorLevel int) *PowerlineThemeColors {
	theme, ok := PowerlineThemes[themeName]
	if !ok {
		return nil
	}
	// Try exact level first, fall back to lower levels.
	switch colorLevel {
	case 2:
		if theme.Truecolor != nil {
			return theme.Truecolor
		}
		if theme.Colors256 != nil {
			return theme.Colors256
		}
		return theme.Basic
	case 1:
		if theme.Colors256 != nil {
			return theme.Colors256
		}
		return theme.Basic
	default:
		return theme.Basic
	}
}
