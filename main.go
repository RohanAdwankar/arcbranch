// arcbranch: A CLI tool to create multiple git branches, worktrees, open VSCode windows, and tile them.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: arcbranch <count> [base-branch] OR arcbranch merge [base-branch] OR arcbranch pymerge")
		os.Exit(1)
	}
	if os.Args[1] == "merge" {
		arcMerge()
		return
	}
	if os.Args[1] == "pymerge" {
		fmt.Println("[arcbranch] Running pytest-based merge (pymerge)...")
		arcMergePytest()
		return
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil || count < 1 {
		fmt.Printf("Error: count must be a positive integer (got '%s')\n", os.Args[1])
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

	// determine next branch index
	out, err := exec.Command("git", "branch", "--list", "arcbranch-*").Output()
	if err != nil {
		fmt.Println("Error listing existing branches:", err)
		os.Exit(1)
	}
	maxIdx := 0
	lines := strings.Split(string(out), "\n")
	re := regexp.MustCompile(`arcbranch-(\d+)`)
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		if parts := re.FindStringSubmatch(line); parts != nil {
			if idx, err := strconv.Atoi(parts[1]); err == nil && idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	startIndex := maxIdx + 1
	var created []string

	for i := 0; i < count; i++ {
		idx := startIndex + i
		branchName := fmt.Sprintf("arcbranch-%d", idx)
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
	arcPath := filepath.Join(repoRoot, ".arcgit")
	var session ArcSession
	session.Base = baseBranch
	session.Parent = parentDir
	// if existing session file, merge previous branches
	if buf, err := os.ReadFile(arcPath); err == nil {
		var prev ArcSession
		if err := json.Unmarshal(buf, &prev); err == nil {
			session.Branches = append(prev.Branches, created...)
		} else {
			session.Branches = created
		}
	} else {
		session.Branches = created
	}
	data, _ := json.MarshalIndent(session, "", "  ")
	_ = os.WriteFile(arcPath, data, 0644)
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

// arcMergePytest attempts a careful merge using pytest results
func arcMergePytest() {
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

	fmt.Println("[arcbranch] Discovered worktrees for merge:")
	for _, b := range session.Branches {
		fmt.Printf("  - %s (worktree: %s)\n", b, filepath.Join(session.Parent, b))
	}

	type mergeResult struct {
		branch    string
		merged    bool
		reason    string
		testFile  string
		pytestOut string
	}
	results := make(chan mergeResult, len(session.Branches))

	for _, b := range session.Branches {
		b := b // capture for goroutine
		go func() {
			wtPath := filepath.Join(session.Parent, b)
			fmt.Printf("[arcbranch][%s] Checking for uncommitted changes...\n", b)
			cmd := exec.Command("git", "status", "--porcelain")
			cmd.Dir = wtPath
			out, err := cmd.Output()
			if err != nil {
				results <- mergeResult{branch: b, merged: false, reason: "git status failed: " + err.Error()}
				return
			}
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			var changedFiles []string
			for _, l := range lines {
				if len(l) > 3 {
					changedFiles = append(changedFiles, strings.TrimSpace(l[3:]))
				}
			}
			if len(changedFiles) == 0 {
				results <- mergeResult{branch: b, merged: false, reason: "no changes"}
				return
			}
			if len(changedFiles) > 1 {
				results <- mergeResult{branch: b, merged: false, reason: "multiple files changed"}
				return
			}
			changedFile := changedFiles[0]
			baseName := filepath.Base(changedFile)
			if !strings.HasSuffix(baseName, ".py") {
				results <- mergeResult{branch: b, merged: false, reason: "not a .py file"}
				return
			}
			testFile := "test_" + baseName
			// Find test file anywhere in repo
			var foundTest string
			_ = filepath.Walk(wtPath, func(path string, info os.FileInfo, err error) error {
				if err == nil && info != nil && info.Name() == testFile {
					foundTest = path
					return filepath.SkipDir
				}
				return nil
			})
			if foundTest == "" {
				results <- mergeResult{branch: b, merged: false, reason: "no test file found"}
				return
			}
			fmt.Printf("[arcbranch][%s] Running pytest on %s...\n", b, foundTest)
			pytestCmd := exec.Command("pytest", foundTest, "--maxfail=1", "--disable-warnings", "-q")
			pytestCmd.Dir = wtPath
			pytestOut, err := pytestCmd.CombinedOutput()
			fmt.Printf("[arcbranch][%s] Pytest output:\n%s\n", b, string(pytestOut))
			if err != nil {
				results <- mergeResult{branch: b, merged: false, reason: "pytest failed", testFile: foundTest, pytestOut: string(pytestOut)}
				return
			}
			addCmd := exec.Command("git", "add", changedFile)
			addCmd.Dir = wtPath
			if err := addCmd.Run(); err != nil {
				results <- mergeResult{branch: b, merged: false, reason: "git add failed: " + err.Error()}
				return
			}
			msg := fmt.Sprintf("[arcbranch merge pytest] Auto-commit %s\n\npytest output:\n%s", changedFile, string(pytestOut))
			commitCmd := exec.Command("git", "commit", "-m", msg)
			commitCmd.Dir = wtPath
			if err := commitCmd.Run(); err != nil {
				results <- mergeResult{branch: b, merged: false, reason: "git commit failed: " + err.Error()}
				return
			}
			results <- mergeResult{branch: b, merged: true, testFile: foundTest, pytestOut: string(pytestOut)}
		}()
	}

	mergedBranches := make(map[string]bool)
	for i := 0; i < len(session.Branches); i++ {
		res := <-results
		if res.merged {
			fmt.Printf("[arcbranch][%s] MERGED (test: %s)\n", res.branch, res.testFile)
		} else {
			fmt.Printf("[arcbranch][%s] NOT merged: %s\n", res.branch, res.reason)
		}
		if res.pytestOut != "" {
			fmt.Printf("[arcbranch][%s] Pytest output:\n%s\n", res.branch, res.pytestOut)
		}
		mergedBranches[res.branch] = res.merged
	}

	fmt.Println("[arcbranch] Merging successful branches into base and cleaning up...")
	exec.Command("git", "checkout", session.Base).Run()
	for _, b := range session.Branches {
		if mergedBranches[b] {
			fmt.Printf("[arcbranch][%s] Merging into base...\n", b)
			if err := exec.Command("git", "merge", "--no-ff", b).Run(); err != nil {
				fmt.Printf("[arcbranch][%s] Conflict merging: %v\n", b, err)
				continue
			}
			wtPath := filepath.Join(session.Parent, b)
			exec.Command("git", "worktree", "remove", "--force", wtPath).Run()
			exec.Command("git", "branch", "-d", b).Run()
			os.RemoveAll(wtPath)
		}
		// For NOT merged branches, do nothing: leave branch and worktree intact
	}
	fmt.Println("[arcbranch] Syncing unmerged branches with base...")
	for _, b := range session.Branches {
		if mergedBranches[b] {
			continue
		}
		wtPath := filepath.Join(session.Parent, b)
		fmt.Printf("[arcbranch][%s] Syncing with base %s...\n", b, session.Base)
		mergeCmd := exec.Command("git", "merge", session.Base)
		mergeCmd.Dir = wtPath
		if err := mergeCmd.Run(); err != nil {
			fmt.Printf("[arcbranch][%s] Error syncing: %v\n", b, err)
		}
	}

	// Only remove .arcgit if all worktrees and branches are gone
	allWorktreesGone := true
	for _, b := range session.Branches {
		wtPath := filepath.Join(session.Parent, b)
		if _, err := os.Stat(wtPath); err == nil {
			allWorktreesGone = false
			break
		}
	}
	allBranchesGone := true
	out, err := exec.Command("git", "branch", "--list", "arcbranch-*").Output()
	if err == nil {
		branchLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range branchLines {
			if strings.HasPrefix(strings.TrimSpace(strings.TrimPrefix(line, "*")), "arcbranch-") {
				allBranchesGone = false
				break
			}
		}
	}
	if allWorktreesGone && allBranchesGone {
		os.Remove(filepath.Join(repoRoot, ".arcgit"))
		fmt.Println("[arcbranch] Pytest merge complete and cleaned up.")
	} else {
		fmt.Println("[arcbranch] Some worktrees or branches remain; .arcgit session preserved.")
	}
}

// tileWindows attempts a best-effort layout based on OS. Users may need to install yabai (macOS) or wmctrl (Linux).
func tileWindows(n int) {
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("macOS tiling not implemented. Install yabai or use AppleScript to tile windows.")
	case "linux":
		fmt.Println("Linux tiling not implemented. Install wmctrl or xdotool to tile windows.")
	case "windows":
		fmt.Println("Windows tiling not implemented. Install NirCmd to tile windows.")
	default:
		fmt.Println("Window tiling not supported on this OS.")
	}
}
