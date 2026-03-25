package widgets_test

import (
	"os"
	"path/filepath"
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func cwdSI(dir string) *input.StatusInput {
	return &input.StatusInput{
		Workspace: &input.WorkspaceInfo{CurrentDir: dir},
	}
}

func TestCWDWidget(t *testing.T) {
	w := mustGetWidget(t, "cwd")

	home, _ := os.UserHomeDir()

	tests := []struct {
		name string
		si   *input.StatusInput
		item config.WidgetItem
		want string
	}{
		{
			name: "nil StatusInput",
			si:   nil,
			item: config.WidgetItem{Type: "cwd"},
			want: "",
		},
		{
			name: "empty path",
			si:   &input.StatusInput{},
			item: config.WidgetItem{Type: "cwd"},
			want: "",
		},
		{
			name: "full absolute path no truncation",
			si:   cwdSI("/usr/local/bin"),
			item: config.WidgetItem{Type: "cwd"},
			want: "/usr/local/bin",
		},
		{
			// Split: ["","usr","local","bin"] → ["/usr","local","bin"]
			// segments=2 → take last 2: ["local","bin"]
			name: "segments=2 returns last 2 components",
			si:   cwdSI("/home/user/projects/myapp"),
			item: config.WidgetItem{Type: "cwd", Segments: 2},
			want: "projects/myapp",
		},
		{
			name: "segments=1 returns last component only",
			si:   cwdSI("/home/user/projects/myapp"),
			item: config.WidgetItem{Type: "cwd", Segments: 1},
			want: "myapp",
		},
		{
			// Split /home/user/projects/myapp → ["/home","user","projects","myapp"]
			// fishStyle: abbreviate i=0..2:
			//   "/home" → len>1, != "~" → runes[:1] = "/"
			//   "user"  → "u"
			//   "projects" → "p"
			// Join → "//u/p/myapp"
			name: "fishStyle=true abbreviates parent dirs",
			si:   cwdSI("/home/user/projects/myapp"),
			item: config.WidgetItem{Type: "cwd", FishStyle: true},
			want: "//u/p/myapp",
		},
		{
			// Single-component absolute path: ["/myapp"], len=1 → fishStyle loop skipped
			name: "fishStyle=true with single component",
			si:   cwdSI("/myapp"),
			item: config.WidgetItem{Type: "cwd", FishStyle: true},
			want: "/myapp",
		},
		{
			name: "home dir abbreviated to tilde",
			si:   cwdSI(filepath.Join(home, "projects")),
			item: config.WidgetItem{Type: "cwd"},
			want: "~/projects",
		},
		{
			name: "home dir itself abbreviated to tilde",
			si:   cwdSI(home),
			item: config.WidgetItem{Type: "cwd"},
			want: "~",
		},
		{
			name: "uses CWD fallback when Workspace is nil",
			si:   &input.StatusInput{CWD: "/tmp/fallback"},
			item: config.WidgetItem{Type: "cwd"},
			want: "/tmp/fallback",
		},
		{
			// segments=2 from /home/user/projects/myapp → ["projects","myapp"]
			// fishStyle: abbreviate i=0: "projects" → "p"
			// Join → "p/myapp"
			name: "fishStyle with segments=2",
			si:   cwdSI("/home/user/projects/myapp"),
			item: config.WidgetItem{Type: "cwd", Segments: 2, FishStyle: true},
			want: "p/myapp",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if home == "" && (tc.name == "home dir abbreviated to tilde" || tc.name == "home dir itself abbreviated to tilde") {
				t.Skip("UserHomeDir not available")
			}
			got := w.Render(tc.si, tc.item)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
