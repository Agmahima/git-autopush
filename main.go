// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// )

// type GitAutoPush struct {
// 	workingDir string
// 	scanner    *bufio.Scanner
// }

// func NewGitAutoPush() *GitAutoPush {
// 	return &GitAutoPush{
// 		workingDir: ".",
// 		scanner:    bufio.NewScanner(os.Stdin),
// 	}
// }

// func (g *GitAutoPush) isGitRepo() bool {
// 	_, err := os.Stat(filepath.Join(g.workingDir, ".git"))
// 	return err == nil
// }

// func (g *GitAutoPush) execGitCommand(args ...string) (string, error) {
// 	cmd := exec.Command("git", args...)
// 	cmd.Dir = g.workingDir
// 	output, err := cmd.CombinedOutput()
// 	return strings.TrimSpace(string(output)), err
// }

// func (g *GitAutoPush) hasRemoteOrigin() bool {
// 	output, err := g.execGitCommand("remote", "get-url", "origin")
// 	return err == nil && output != ""
// }

// func (g *GitAutoPush) getCurrentBranch() (string, error) {
// 	return g.execGitCommand("rev-parse", "--abrev-ref", "HEAD")
// }

// func (g *GitAutoPush) branchExistsOnRemote(branch string) bool {
// 	_, err := g.execGitCommand("ls-remote", "--heads", "origin", branch)
// 	return err == nil
// }

// func (g *GitAutoPush) getGitStatus() (string, error) {
// 	return g.execGitCommand("status", "--porcelain")
// }

// func (g *GitAutoPush) hasUncommitedChanges() bool {
// 	status, err := g.getGitStatus()
// 	return err == nil && status != ""
// }

// func (g *GitAutoPush) promptYesNo(question string) bool {
// 	fmt.Printf("%s (y/n): ", question)
// 	g.scanner.Scan()
// 	response := strings.ToLower(strings.TrimSpace(g.scanner.Text()))
// 	return response == "y" || response == "yes"
// }

// func (g *GitAutoPush) getCommitMessage() string {
// 	fmt.Print("Enter commit message: ")
// 	g.scanner.Scan()
// 	message := strings.TrimSpace(g.scanner.Text())
// 	if message == "" {
// 		return "Auto commit: Updated Code"
// 	}
// 	return message
// }

// func (g *GitAutoPush) initGitRepo() error {
// 	fmt.Println("Initializing new Git repository...")
// 	_, err := g.execGitCommand(("init"))
// 	return err
// }

// func (g *GitAutoPush) addRemoteOrigin() error {
// 	fmt.Print("Enter Github Repository URL: ")
// 	g.scanner.Scan()
// 	repoURL := strings.TrimSpace(g.scanner.Text())

// 	if repoURL == "" {
// 		return fmt.Errorf("repository URL cannot be empty")
// 	}

// 	fmt.Println("Adding remote origin...")
// 	_, err := g.execGitCommand("remote", "add", "origin", repoURL)
// 	return err
// }

// func (g *GitAutoPush) stageChanges() error {
// 	fmt.Println("Staging changes...")
// 	_, err := g.execGitCommand("add", ".")
// 	return err
// }

// func (g *GitAutoPush) commitChanges(message string) error {
// 	fmt.Printf("Committing changes with messages: '%s'\n", message)
// 	_, err := g.execGitCommand("commit", "-m", message)
// 	return err
// }

// func (g *GitAutoPush) pushToRemote(branch string, setUpstream bool) error {
// 	fmt.Printf("Pushing changes to remote branch '%s'...\n", branch)

// 	if setUpstream {
// 		_, err := g.execGitCommand("push", "-u", "origin", branch)
// 		return err
// 	}
// 	_, err := g.execGitCommand("push", "origin", branch)
// 	return err
// }

// func (g *GitAutoPush) handleFirstTimePush() error {
// 	fmt.Println(" first time setup detected!")

// 	if !g.isGitRepo() {
// 		if err := g.initGitRepo(); err != nil {
// 			return fmt.Errorf("failed to initialize git repository : %v", err)
// 		}
// 	}
// 	if !g.hasRemoteOrigin() {
// 		if err := g.addRemoteOrigin(); err != nil {
// 			return fmt.Errorf("failed to add remote origin: %v", err)
// 		}
// 	}

// 	if err := g.stageChanges(); err != nil {
// 		return fmt.Errorf("failed to stage changes: %v", err)
// 	}
// 	commitMessage := g.getCommitMessage()

// 	if err := g.commitChanges(commitMessage); err != nil {
// 		return fmt.Errorf("failed to commit changes: %v", err)
// 	}

// 	branch, err := g.getCurrentBranch()
// 	if err != nil {
// 		return fmt.Errorf("failed to get current branch: %v", err)
// 	}

// 	if err := g.pushToRemote(branch, true); err != nil {
// 		return fmt.Errorf("failed to push to remote: %v", err)
// 	}

// 	fmt.Println("First time push completed successfully!")
// 	return nil

// }

// func (g *GitAutoPush) handleSubsequentPush() error {
// 	fmt.Println("Subsequent push detected!")
// 	branch, err := g.getCurrentBranch()
// 	if err != nil {
// 		return fmt.Errorf("failed to get current branch: %v", err)
// 	}

// 	if !g.hasUncommitedChanges() {
// 		fmt.Println("No changes to commit. Exiting.")
// 		return nil
// 	}

// 	if !g.promptYesNo(fmt.Sprintf("Push changes to branch '%s'?", branch)) {
// 		fmt.Println("Push cancelled by user.")
// 		return nil
// 	}

// 	if err := g.stageChanges(); err != nil {
// 		return fmt.Errorf("failed to stage changes: %v", err)
// 	}

// 	commitMessage := g.getCommitMessage()

// 	if err := g.commitChanges(commitMessage); err != nil {
// 		return fmt.Errorf("failed to commit changes: %v", err)
// 	}

// 	if !g.branchExistsOnRemote(branch) {
// 		if g.promptYesNo(fmt.Sprintf("Branch '%s' does not exist on remote. Create it?", branch)) {
// 			if err := g.pushToRemote(branch, true); err != nil {
// 				return fmt.Errorf("failed to push new branch to remote: %v", err)
// 			}
// 		} else {
// 			fmt.Println("Push cancelled.")
// 			return nil
// 		}
// 	} else {
// 		if err := g.pushToRemote(branch, false); err != nil {
// 			return fmt.Errorf("failed to push changes to remote: %v", err)
// 		}

// 	}

// 	fmt.Println("Push completed successfully!")
// 	return nil
// }

// func (g *GitAutoPush) Run() error {
// 	fmt.Println("Welcome to Git Auto Push!")
// 	fmt.Println("==========================")

// 	isGitRepo := g.isGitRepo()
// 	hasRemote := g.hasRemoteOrigin()

// 	if !isGitRepo || !hasRemote {
// 		return g.handleFirstTimePush()
// 	} else {
// 		return g.handleSubsequentPush()
// 	}
// }

// func main() {
// 	autoPush := NewGitAutoPush()

// 	if err := autoPush.Run; err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		os.Exit(1)
// 	}
// }


package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitAutoPush struct {
	workingDir string
	scanner    *bufio.Scanner
}

func NewGitAutoPush() *GitAutoPush {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	return &GitAutoPush{
		workingDir: ".",
		scanner:    scanner,
	}
}

// Check if current directory is a git repository
func (g *GitAutoPush) isGitRepo() bool {
	_, err := os.Stat(filepath.Join(g.workingDir, ".git"))
	return err == nil
}

// Execute git command and return output
func (g *GitAutoPush) execGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.workingDir
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// Check if remote origin exists
func (g *GitAutoPush) hasRemoteOrigin() bool {
	output, err := g.execGitCommand("remote", "get-url", "origin")
	return err == nil && output != ""
}

// Get current branch name
func (g *GitAutoPush) getCurrentBranch() (string, error) {
	return g.execGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

// Check if branch exists on remote
func (g *GitAutoPush) branchExistsOnRemote(branch string) bool {
	_, err := g.execGitCommand("ls-remote", "--heads", "origin", branch)
	return err == nil
}

// Get git status
func (g *GitAutoPush) getGitStatus() (string, error) {
	return g.execGitCommand("status", "--porcelain")
}

// Check if there are uncommitted changes
func (g *GitAutoPush) hasUncommittedChanges() bool {
	status, err := g.getGitStatus()
	return err == nil && status != ""
}

// Get user input for yes/no questions
func (g *GitAutoPush) promptYesNo(question string) bool {
	fmt.Printf("%s (y/n): ", question)
	if !g.scanner.Scan() {
		if err := g.scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return false
		}
		return false
	}
	response := strings.ToLower(strings.TrimSpace(g.scanner.Text()))
	return response == "y" || response == "yes"
}

// Get user input for commit message
func (g *GitAutoPush) getCommitMessage() string {
	fmt.Print("Enter commit message: ")
	if !g.scanner.Scan() {
		if err := g.scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
		}
		return "Auto-commit: Updated code"
	}
	message := strings.TrimSpace(g.scanner.Text())
	if message == "" {
		return "Auto-commit: Updated code"
	}
	return message
}

// Initialize git repository
func (g *GitAutoPush) initGitRepo() error {
	fmt.Println("Initializing git repository...")
	_, err := g.execGitCommand("init")
	return err
}

// Add remote origin
func (g *GitAutoPush) addRemoteOrigin() error {
	fmt.Print("Enter GitHub repository URL: ")
	if !g.scanner.Scan() {
		if err := g.scanner.Err(); err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}
		return fmt.Errorf("failed to read repository URL")
	}
	
	repoURL := strings.TrimSpace(g.scanner.Text())
	if repoURL == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}
	
	fmt.Println("Adding remote origin...")
	_, err := g.execGitCommand("remote", "add", "origin", repoURL)
	return err
}

// Stage all changes
func (g *GitAutoPush) stageChanges() error {
	fmt.Println("Staging changes...")
	_, err := g.execGitCommand("add", ".")
	return err
}

// Commit changes
func (g *GitAutoPush) commitChanges(message string) error {
	fmt.Printf("Committing changes with message: '%s'\n", message)
	_, err := g.execGitCommand("commit", "-m", message)
	return err
}

// Push to remote branch
func (g *GitAutoPush) pushToRemote(branch string, setUpstream bool) error {
	fmt.Printf("Pushing to remote branch '%s'...\n", branch)
	
	if setUpstream {
		_, err := g.execGitCommand("push", "-u", "origin", branch)
		return err
	}
	
	_, err := g.execGitCommand("push", "origin", branch)
	return err
}

// Handle first-time push
func (g *GitAutoPush) handleFirstTimePush() error {
	fmt.Println("🚀 First time setup detected!")
	
	// Initialize git if not already initialized
	if !g.isGitRepo() {
		if err := g.initGitRepo(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %v", err)
		}
	}
	
	// Add remote origin if not exists
	if !g.hasRemoteOrigin() {
		if err := g.addRemoteOrigin(); err != nil {
			return fmt.Errorf("failed to add remote origin: %v", err)
		}
	}
	
	// Stage changes
	if err := g.stageChanges(); err != nil {
		return fmt.Errorf("failed to stage changes: %v", err)
	}
	
	// Get commit message
	commitMessage := g.getCommitMessage()
	
	// Commit changes
	if err := g.commitChanges(commitMessage); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}
	
	// Get current branch
	branch, err := g.getCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	
	// Push to remote (first time, so set upstream)
	if err := g.pushToRemote(branch, true); err != nil {
		return fmt.Errorf("failed to push to remote: %v", err)
	}
	
	fmt.Println("✅ First-time push completed successfully!")
	return nil
}

// Handle subsequent pushes
func (g *GitAutoPush) handleSubsequentPush() error {
	fmt.Println("🔄 Subsequent push detected!")
	
	// Get current branch
	branch, err := g.getCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	
	// Check if there are uncommitted changes
	if !g.hasUncommittedChanges() {
		fmt.Println("No changes to commit.")
		return nil
	}
	
	// Ask user if they want to proceed
	if !g.promptYesNo(fmt.Sprintf("Push changes to branch '%s'?", branch)) {
		fmt.Println("Push cancelled.")
		return nil
	}
	
	// Stage changes
	if err := g.stageChanges(); err != nil {
		return fmt.Errorf("failed to stage changes: %v", err)
	}
	
	// Get commit message
	commitMessage := g.getCommitMessage()
	
	// Commit changes
	if err := g.commitChanges(commitMessage); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}
	
	// Check if we need to create a new branch or use existing
	if !g.branchExistsOnRemote(branch) {
		if g.promptYesNo(fmt.Sprintf("Branch '%s' doesn't exist on remote. Create it?", branch)) {
			if err := g.pushToRemote(branch, true); err != nil {
				return fmt.Errorf("failed to push new branch: %v", err)
			}
		} else {
			fmt.Println("Push cancelled.")
			return nil
		}
	} else {
		// Push to existing branch
		if err := g.pushToRemote(branch, false); err != nil {
			return fmt.Errorf("failed to push to existing branch: %v", err)
		}
	}
	
	fmt.Println("✅ Push completed successfully!")
	return nil
}

// Main automation logic
func (g *GitAutoPush) Run() error {
	fmt.Println("🤖 Git Auto-Push Tool")
	fmt.Println("=====================")
	
	// Check if this is a git repository and has remote origin
	isGitRepo := g.isGitRepo()
	hasRemote := g.hasRemoteOrigin()
	
	// Determine if this is first-time setup or subsequent push
	if !isGitRepo || !hasRemote {
		return g.handleFirstTimePush()
	} else {
		return g.handleSubsequentPush()
	}
}

func main() {
	autoPush := NewGitAutoPush()
	
	if err := autoPush.Run(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}

// new code