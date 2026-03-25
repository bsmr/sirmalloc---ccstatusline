package widgets_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
	"go.a8l.eu/ccstatusline/internal/widgets"
)

// initGitRepo creates a temporary git repository with one initial commit.
// The test is skipped if git is not available.
func initGitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available:", err)
	}
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")
	return dir
}

// ---------------------------------------------------------------------------
// git-branch
// ---------------------------------------------------------------------------

func TestGitBranchWidget_NilInput(t *testing.T) {
	w, _ := widgets.Get("git-branch")
	if got := w.Render(nil, config.WidgetItem{Type: "git-branch"}); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestGitBranchWidget_EmptyCWD(t *testing.T) {
	w, _ := widgets.Get("git-branch")
	if got := w.Render(&input.StatusInput{}, config.WidgetItem{Type: "git-branch"}); got != "" {
		t.Errorf("expected empty for empty CWD, got %q", got)
	}
}

func TestGitBranchWidget_NonGitDir(t *testing.T) {
	w, _ := widgets.Get("git-branch")
	si := &input.StatusInput{CWD: t.TempDir()}
	if got := w.Render(si, config.WidgetItem{Type: "git-branch"}); got != "" {
		t.Errorf("expected empty for non-git dir, got %q", got)
	}
}

func TestGitBranchWidget_GitRepo(t *testing.T) {
	dir := initGitRepo(t)
	w, _ := widgets.Get("git-branch")
	si := &input.StatusInput{CWD: dir}
	got := w.Render(si, config.WidgetItem{Type: "git-branch"})
	if got == "" {
		t.Error("expected non-empty branch name in git repo")
	}
}

func TestGitBranchWidget_HideNoGit(t *testing.T) {
	w, _ := widgets.Get("git-branch")
	si := &input.StatusInput{CWD: t.TempDir()}
	item := config.WidgetItem{Type: "git-branch", HideNoGit: true}
	if got := w.Render(si, item); got != "" {
		t.Errorf("expected empty with HideNoGit in non-git dir, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// git-worktree
// ---------------------------------------------------------------------------

func TestGitWorktreeWidget_NilInput(t *testing.T) {
	w, _ := widgets.Get("git-worktree")
	if got := w.Render(nil, config.WidgetItem{Type: "git-worktree"}); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestGitWorktreeWidget_NonGitDir(t *testing.T) {
	w, _ := widgets.Get("git-worktree")
	si := &input.StatusInput{CWD: t.TempDir()}
	if got := w.Render(si, config.WidgetItem{Type: "git-worktree"}); got != "" {
		t.Errorf("expected empty for non-git dir, got %q", got)
	}
}

func TestGitWorktreeWidget_GitRepo(t *testing.T) {
	dir := initGitRepo(t)
	w, _ := widgets.Get("git-worktree")
	si := &input.StatusInput{CWD: dir}
	got := w.Render(si, config.WidgetItem{Type: "git-worktree"})
	if got == "" {
		t.Error("expected non-empty worktree name in git repo")
	}
	if want := filepath.Base(dir); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

// ---------------------------------------------------------------------------
// git-changes
// ---------------------------------------------------------------------------

func TestGitChangesWidget_NilInput(t *testing.T) {
	w, _ := widgets.Get("git-changes")
	if got := w.Render(nil, config.WidgetItem{Type: "git-changes"}); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestGitChangesWidget_NonGitDir(t *testing.T) {
	w, _ := widgets.Get("git-changes")
	si := &input.StatusInput{CWD: t.TempDir()}
	if got := w.Render(si, config.WidgetItem{Type: "git-changes"}); got != "" {
		t.Errorf("expected empty for non-git dir, got %q", got)
	}
}

func TestGitChangesWidget_CleanRepo(t *testing.T) {
	dir := initGitRepo(t)
	w, _ := widgets.Get("git-changes")
	si := &input.StatusInput{CWD: dir}
	if got := w.Render(si, config.WidgetItem{Type: "git-changes"}); got != "" {
		t.Errorf("expected empty for clean repo, got %q", got)
	}
}

func TestGitChangesWidget_WithOneChange(t *testing.T) {
	dir := initGitRepo(t)
	// Modify the tracked file without staging.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	w, _ := widgets.Get("git-changes")
	si := &input.StatusInput{CWD: dir}
	got := w.Render(si, config.WidgetItem{Type: "git-changes"})
	if got != "1 file" {
		t.Errorf(`expected "1 file", got %q`, got)
	}
}

func TestGitChangesWidget_WithMultipleChanges(t *testing.T) {
	dir := initGitRepo(t)
	// Add a second tracked file.
	if err := os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Run() //nolint:errcheck
	}
	run("add", "file2.txt")
	run("commit", "-m", "add file2")

	// Modify both files without staging.
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("changed"), 0o644)  //nolint:errcheck
	os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("changed"), 0o644) //nolint:errcheck

	w, _ := widgets.Get("git-changes")
	si := &input.StatusInput{CWD: dir}
	got := w.Render(si, config.WidgetItem{Type: "git-changes"})
	if got != "2 files" {
		t.Errorf(`expected "2 files", got %q`, got)
	}
}
