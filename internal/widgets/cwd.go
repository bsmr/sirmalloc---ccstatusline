package widgets

import (
	"os"
	"path/filepath"
	"strings"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("cwd", cwdWidget{})
}

type cwdWidget struct{}

func (cwdWidget) Render(si *input.StatusInput, item config.WidgetItem) string {
	if si == nil {
		return ""
	}
	path := ""
	if si.Workspace != nil && si.Workspace.CurrentDir != "" {
		path = si.Workspace.CurrentDir
	} else {
		path = si.CWD
	}
	if path == "" {
		return ""
	}
	return formatCWD(path, item.Segments, item.FishStyle)
}

func formatCWD(path string, segments int, fishStyle bool) string {
	home, _ := os.UserHomeDir()
	if home != "" && strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}
	parts := strings.Split(filepath.ToSlash(path), "/")
	// remove empty leading part for absolute paths and re-attach root
	if len(parts) > 1 && parts[0] == "" {
		parts = parts[1:]
		parts[0] = "/" + parts[0]
	}
	if segments > 0 && len(parts) > segments {
		parts = parts[len(parts)-segments:]
	}
	if fishStyle && len(parts) > 1 {
		for i := 0; i < len(parts)-1; i++ {
			if len(parts[i]) > 1 && parts[i] != "~" {
				r := []rune(parts[i])
				parts[i] = string(r[:1])
			}
		}
	}
	return strings.Join(parts, "/")
}
