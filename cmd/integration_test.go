package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// These tests run the compiled binary as a subprocess.
// They require `make build` to have been run first.

func binaryPath() string {
	// Walk up to find the project root and the binary
	wd, _ := os.Getwd()
	for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
		bin := filepath.Join(dir, "bin", "prm")
		if _, err := os.Stat(bin); err == nil {
			return bin
		}
	}
	return ""
}

func runPRM(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	bin := binaryPath()
	if bin == "" {
		t.Skip("binary not found; run 'make build' first")
	}
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func setupIntegrationDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	out, err := runPRM(t, dir, "init", "--name", "integration-test")
	if err != nil {
		t.Fatalf("init failed: %s\n%v", out, err)
	}
	return dir
}

func TestCLI_Init(t *testing.T) {
	dir := t.TempDir()

	out, err := runPRM(t, dir, "init", "--name", "my-project")
	if err != nil {
		t.Fatalf("init failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Initialized") {
		t.Errorf("init output = %q, want to contain 'Initialized'", out)
	}

	// .prm directory should exist
	if _, err := os.Stat(filepath.Join(dir, ".prm")); err != nil {
		t.Errorf(".prm dir not created: %v", err)
	}
}

func TestCLI_Init_Double(t *testing.T) {
	dir := t.TempDir()
	runPRM(t, dir, "init")

	_, err := runPRM(t, dir, "init")
	if err == nil {
		t.Error("double init should fail")
	}
}

func TestCLI_EpicLifecycle(t *testing.T) {
	dir := setupIntegrationDir(t)

	// Create
	out, err := runPRM(t, dir, "epic", "create", "--title", "Test Epic", "--priority", "high", "--tags", "backend,security")
	if err != nil {
		t.Fatalf("epic create failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Test Epic") {
		t.Errorf("create output missing title: %s", out)
	}

	// List
	out, _ = runPRM(t, dir, "epic", "list")
	if !strings.Contains(out, "Test Epic") {
		t.Errorf("list output missing epic: %s", out)
	}

	// Show
	out, err = runPRM(t, dir, "epic", "show", "test-epic")
	if err != nil {
		t.Fatalf("epic show failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Test Epic") {
		t.Errorf("show output missing title: %s", out)
	}

	// Update
	out, err = runPRM(t, dir, "epic", "update", "test-epic", "--title", "Updated Epic", "--priority", "critical")
	if err != nil {
		t.Fatalf("epic update failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Updated Epic") {
		t.Errorf("update output missing new title: %s", out)
	}

	// Delete
	out, err = runPRM(t, dir, "epic", "delete", "test-epic")
	if err != nil {
		t.Fatalf("epic delete failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Deleted") {
		t.Errorf("delete output missing confirmation: %s", out)
	}

	// Verify deleted
	out, _ = runPRM(t, dir, "epic", "list")
	if strings.Contains(out, "Updated Epic") {
		t.Error("deleted epic still appears in list")
	}
}

func TestCLI_Hierarchy(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "epic", "create", "--title", "E1", "--priority", "high")
	runPRM(t, dir, "story", "create", "--epic", "e1", "--title", "S1")
	runPRM(t, dir, "subtask", "create", "--story", "s1", "--title", "ST1")

	// Tree should show hierarchy
	out, err := runPRM(t, dir, "tree")
	if err != nil {
		t.Fatalf("tree failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "E1") || !strings.Contains(out, "S1") || !strings.Contains(out, "ST1") {
		t.Errorf("tree missing entities: %s", out)
	}
}

func TestCLI_Status(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "task", "create", "--title", "My Task")

	out, err := runPRM(t, dir, "status", "my-task", "in-progress")
	if err != nil {
		t.Fatalf("status failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "in-progress") {
		t.Errorf("status output = %s, want 'in-progress'", out)
	}
}

func TestCLI_Comment(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "task", "create", "--title", "My Task")

	out, err := runPRM(t, dir, "comment", "my-task", "--text", "Hello world")
	if err != nil {
		t.Fatalf("comment failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Comment added") {
		t.Errorf("comment output = %s", out)
	}

	// Verify comment appears in show
	out, _ = runPRM(t, dir, "task", "show", "my-task")
	if !strings.Contains(out, "Hello world") {
		t.Errorf("show doesn't contain comment: %s", out)
	}
}

func TestCLI_Search(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "task", "create", "--title", "Auth Module")
	runPRM(t, dir, "task", "create", "--title", "Database Setup")

	out, err := runPRM(t, dir, "search", "auth")
	if err != nil {
		t.Fatalf("search failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "Auth Module") {
		t.Errorf("search missing result: %s", out)
	}
	if strings.Contains(out, "Database Setup") {
		t.Errorf("search should not include non-matching: %s", out)
	}
}

func TestCLI_Dashboard(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "task", "create", "--title", "T1", "--priority", "high")
	runPRM(t, dir, "bug", "create", "--title", "B1", "--severity", "major")

	out, err := runPRM(t, dir, "dashboard")
	if err != nil {
		t.Fatalf("dashboard failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "integration-test") {
		t.Errorf("dashboard missing project name: %s", out)
	}
	if !strings.Contains(out, "Total:") {
		t.Errorf("dashboard missing total: %s", out)
	}
}

func TestCLI_List_WithFilters(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "task", "create", "--title", "High Task", "--priority", "high")
	runPRM(t, dir, "task", "create", "--title", "Low Task", "--priority", "low")

	out, _ := runPRM(t, dir, "list", "--priority", "high")
	if !strings.Contains(out, "High Task") {
		t.Errorf("filtered list missing high task: %s", out)
	}
	if strings.Contains(out, "Low Task") {
		t.Errorf("filtered list should not contain low task: %s", out)
	}
}

func TestCLI_Export_CSV(t *testing.T) {
	dir := setupIntegrationDir(t)
	runPRM(t, dir, "task", "create", "--title", "Export Task")

	outFile := filepath.Join(dir, "export.csv")
	out, err := runPRM(t, dir, "export", "--format", "csv", "--output", outFile)
	if err != nil {
		t.Fatalf("export failed: %s\n%v", out, err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading export: %v", err)
	}
	if !strings.Contains(string(data), "Export Task") {
		t.Errorf("CSV missing entity: %s", string(data))
	}
}

func TestCLI_Export_JSON(t *testing.T) {
	dir := setupIntegrationDir(t)
	runPRM(t, dir, "task", "create", "--title", "JSON Task")

	outFile := filepath.Join(dir, "export.json")
	_, err := runPRM(t, dir, "export", "--format", "json", "--output", outFile)
	if err != nil {
		t.Fatalf("export json failed: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "JSON Task") {
		t.Errorf("JSON missing entity: %s", string(data))
	}
}

func TestCLI_Reindex(t *testing.T) {
	dir := setupIntegrationDir(t)

	runPRM(t, dir, "task", "create", "--title", "T1")
	runPRM(t, dir, "task", "create", "--title", "T2")

	out, err := runPRM(t, dir, "reindex")
	if err != nil {
		t.Fatalf("reindex failed: %s\n%v", out, err)
	}
	if !strings.Contains(out, "2 entities") {
		t.Errorf("reindex output = %s, want '2 entities'", out)
	}
}

func TestCLI_Doc(t *testing.T) {
	dir := setupIntegrationDir(t)

	out, err := runPRM(t, dir, "doc", "create", "--title", "Design Notes")
	if err != nil {
		t.Fatalf("doc create failed: %s\n%v", out, err)
	}

	out, _ = runPRM(t, dir, "doc", "list")
	if !strings.Contains(out, "design-notes.md") {
		t.Errorf("doc list missing file: %s", out)
	}

	out, _ = runPRM(t, dir, "doc", "show", "design-notes")
	if !strings.Contains(out, "Design Notes") {
		t.Errorf("doc show missing title: %s", out)
	}
}

func TestCLI_Bug_WithSeverity(t *testing.T) {
	dir := setupIntegrationDir(t)

	out, err := runPRM(t, dir, "bug", "create", "--title", "Crash Bug", "--severity", "blocker", "--priority", "critical")
	if err != nil {
		t.Fatalf("bug create failed: %s\n%v", out, err)
	}

	out, _ = runPRM(t, dir, "bug", "show", "crash-bug")
	if !strings.Contains(out, "blocker") {
		t.Errorf("show missing severity: %s", out)
	}
}

func TestCLI_MissingTitle_Fails(t *testing.T) {
	dir := setupIntegrationDir(t)

	_, err := runPRM(t, dir, "task", "create")
	if err == nil {
		t.Error("task create without --title should fail")
	}
}

func TestCLI_InvalidStatus_Fails(t *testing.T) {
	dir := setupIntegrationDir(t)
	runPRM(t, dir, "task", "create", "--title", "T1")

	_, err := runPRM(t, dir, "status", "t1", "invalid")
	if err == nil {
		t.Error("invalid status should fail")
	}
}

func TestCLI_StoryWithoutEpic_Fails(t *testing.T) {
	dir := setupIntegrationDir(t)

	_, err := runPRM(t, dir, "story", "create", "--title", "Orphan")
	if err == nil {
		t.Error("story create without --epic should fail")
	}
}

func TestCLI_NotInitialized_Fails(t *testing.T) {
	dir := t.TempDir()

	_, err := runPRM(t, dir, "list")
	if err == nil {
		t.Error("commands on uninitialized project should fail")
	}
}
