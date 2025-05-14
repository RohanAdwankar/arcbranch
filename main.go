// arcbranch: A CLI tool to create multiple git branches, worktrees, open VSCode windows, and tile them.

package main

import (
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: arcbranch <count> [base-branch]")
        os.Exit(1)
    }
    count, err := strconv.Atoi(os.Args[1])
    if err != nil || count < 1 {
        fmt.Println("Error: count must be a positive integer")
        os.Exit(1)
    }

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

    for i := 1; i <= count; i++ {
        branchName := fmt.Sprintf("arcbranch-%d", i)
        fmt.Printf("Creating branch %s from %s\n", branchName, baseBranch)
        if err := exec.Command("git", "branch", branchName, baseBranch).Run(); err != nil {
            fmt.Println("Error creating branch", branchName, ":", err)
            continue
        }

        worktreePath := branchName
        fmt.Printf("Adding worktree at %s\n", worktreePath)
        if err := exec.Command("git", "worktree", "add", worktreePath, branchName).Run(); err != nil {
            fmt.Println("Error adding worktree for", branchName, ":", err)
            continue
        }

        fmt.Printf("Opening VSCode for %s\n", worktreePath)
        if err := exec.Command("code", "-n", worktreePath).Start(); err != nil {
            fmt.Println("Error opening VSCode for", worktreePath, ":", err)
        }
    }

    fmt.Println("Tiling windows...")
    tileWindows(count)
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