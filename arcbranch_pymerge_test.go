package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to run a shell command in a given directory
func runCmd(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestArcbranchPymergeSessionPreserved(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "arcbranch-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy example files into tempDir
	exampleDir := filepath.Join("examples")
	exampleAbs, err := filepath.Abs(exampleDir)
	if err != nil {
		t.Fatalf("failed to get abs path for examples: %v", err)
	}
	out, err := runCmd(tempDir, "cp", "-r", exampleAbs+"/.", ".")
	if err != nil {
		t.Fatalf("failed to copy examples: %v, %s", err, out)
	}

	// Init git repo
	_, err = runCmd(tempDir, "git", "init")
	if err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	_, err = runCmd(tempDir, "git", "add", ".")
	if err != nil {
		t.Fatalf("git add failed: %v", err)
	}
	_, err = runCmd(tempDir, "git", "commit", "-m", "Initial commit")
	if err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	// Run arcbranch 2 (create 2 branches)
	_, err = runCmd(tempDir, "arcbranch", "2")
	if err != nil {
		t.Fatalf("arcbranch 2 failed: %v", err)
	}

	// Run arcbranch pymerge twice
	for i := 0; i < 2; i++ {
		_, err = runCmd(tempDir, "arcbranch", "pymerge")
		if err != nil && !strings.Contains(err.Error(), "No .arcgit session found") {
			t.Fatalf("arcbranch pymerge run %d failed: %v", i+1, err)
		}
		arcgitPath := filepath.Join(tempDir, ".arcgit")
		if _, err := os.Stat(arcgitPath); os.IsNotExist(err) {
			t.Fatalf(".arcgit was deleted after pymerge run %d, but should be preserved", i+1)
		}
	}
}
