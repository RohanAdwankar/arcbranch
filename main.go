// arcbranch: A CLI tool to create multiple git branches, worktrees, open VSCode windows, and tile them.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: arcbranch <count> [base-branch] OR arcbranch merge")
		os.Exit(1)
	}
	if os.Args[1] == "merge" {
		arcMerge()
		return
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil || count < 1 {
		fmt.Println("Error: count must be a positive integer")
		os.Exit(1)
	}

	// determine repo root and parent dir for worktrees
	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}
	parentDir := filepath.Dir(repoRoot)

	var baseBranch string
	if len(os.Args) >= 3 {
		baseBranch = os.Args[2]
	} else {
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			fmt.Println("Error getting current branch:", err)
			os.Exit(1)
		}
		baseBranch = strings.TrimSpace(string(out))
	}

	var created []string
	for i := 1; i <= count; i++ {
		branchName := fmt.Sprintf("arcbranch-%d", i)
		fmt.Printf("Creating branch %s from %s\n", branchName, baseBranch)
		if err := exec.Command("git", "branch", branchName, baseBranch).Run(); err != nil {
			fmt.Println("Error creating branch", branchName, ":", err)
			continue
		}

		worktreePath := filepath.Join(parentDir, branchName)
		fmt.Printf("Adding worktree at %s\n", worktreePath)
		if err := exec.Command("git", "worktree", "add", worktreePath, branchName).Run(); err != nil {
			fmt.Println("Error adding worktree for", branchName, ":", err)
			continue
		}
		created = append(created, branchName)

		fmt.Printf("Opening VSCode for %s\n", worktreePath)
		if err := exec.Command("code", "-n", worktreePath).Start(); err != nil {
			fmt.Println("Error opening VSCode for", worktreePath, ":", err)
		}
	}

	fmt.Println("Tiling windows...")
	tileWindows(count)

	// record session for merging
	session := ArcSession{Base: baseBranch, Branches: created, Parent: parentDir}
	data, _ := json.MarshalIndent(session, "", "  ")
	_ = os.WriteFile(filepath.Join(repoRoot, ".arcgit"), data, 0644)
}

// ArcSession tracks branches for merge
type ArcSession struct {
	Base     string   `json:"base_branch"`
	Branches []string `json:"branches"`
	Parent   string   `json:"parent_dir"`
}

// arcMerge reads the last session from .arcgit and merges branches back, then cleans up
func arcMerge() {
	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}
	buf, err := os.ReadFile(filepath.Join(repoRoot, ".arcgit"))
	if err != nil {
		fmt.Println("No .arcgit session found or error reading it:", err)
		os.Exit(1)
	}
	var session ArcSession
	if err := json.Unmarshal(buf, &session); err != nil {
		fmt.Println("Error parsing .arcgit:", err)
		os.Exit(1)
	}
	// checkout base
	fmt.Println("Checking out base branch", session.Base)
	exec.Command("git", "checkout", session.Base).Run()
	for _, b := range session.Branches {
		fmt.Println("Merging branch", b)
		if err := exec.Command("git", "merge", "--no-ff", b).Run(); err != nil {
			fmt.Println("Conflict merging", b, ":", err)
			os.Exit(1)
		}
		// remove worktree and branch
		wtPath := filepath.Join(session.Parent, b)
		exec.Command("git", "worktree", "remove", "--force", wtPath).Run()
		exec.Command("git", "branch", "-d", b).Run()
		os.RemoveAll(wtPath)
	}
	os.Remove(filepath.Join(repoRoot, ".arcgit"))
	fmt.Println("Merge complete and cleaned up.")
}

// tileWindows attempts a best-effort layout based on OS. Users may need to install yabai (macOS) or wmctrl (Linux).
func tileWindows(n int) {
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("macOS tiling not implemented. Install yabai or use AppleScript to tile windows.")
	case "linux":
		fmt.Println("Linux tiling not implemented. Install wmctrl or xdotool to tile windows.")
	default:
		fmt.Println("Window tiling not supported on this OS.")
	}
}
