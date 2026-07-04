package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Session represents a registered session linking a local workspace project to an active agent container.
type Session struct {
	ID            string   `json:"id"`
	ContainerName string   `json:"container_name"`
	LocalPath     string   `json:"local_path"`
	LoadedSkills  []string `json:"loaded_skills"`
}

// FindProjectRoot traverses upward from the startPath looking for a .git directory or a go.mod file.
func FindProjectRoot(startPath string) (string, error) {
	currentDir, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		// Check for .git or go.mod
		gitPath := filepath.Join(currentDir, ".git")
		goModPath := filepath.Join(currentDir, "go.mod")

		if _, err := os.Stat(gitPath); err == nil {
			return currentDir, nil
		}
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root without finding project files
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("no project root (.git or go.mod) found starting from: %s", startPath)
}

// BubbleTeaModel represents the TUI state for selection processes
type BubbleTeaModel struct {
	Containers []Container
	Cursor     int
	Selected   string
	State      string // e.g. "selecting_container", "copying", "done"
	Err        error
}

// LoadSessionConfig registers or configures the environment (skills, prompts, instructions)
// that will be injected alongside the project workspace inside the agent's volume.
func LoadSessionConfig(session *Session, skillsDir string) error {
	// TODO:
	// 1. Scan the skills/instructions directory
	// 2. Identify skills (e.g. SKILL.md, files in .agents/)
	// 3. Keep track of them in the session metadata
	// 4. Copy them to '/workspace/.agents/' inside the target named volume
	return nil
}

// InitializeOpencodeSession sets up the entire workspace, registers it, and initiates working context.
func InitializeOpencodeSession(
	ctx context.Context,
	hostPath string,
	containerName string,
) (*Session, error) {
	rootPath, err := FindProjectRoot(hostPath)
	if err != nil {
		return nil, fmt.Errorf("could not resolve project root: %w", err)
	}

	session := &Session{
		ID:            filepath.Base(rootPath),
		ContainerName: containerName,
		LocalPath:     rootPath,
	}

	// TODO:
	// 1. Clear container workspace if needed using ClearWorkspace(ctx, containerName)
	// 2. Copy project root into container volume using CopyToWorkspace(ctx, containerName, rootPath)

	return session, nil
}

func InitProjectSession(
	ctx context.Context,
	hostPath string, // e.g., working directory "."
	containerName string, // Picked by user via TUI
) (*Session, error) {
	rootPath, err := FindProjectRoot(hostPath)
	if err != nil {
		return nil, fmt.Errorf("could not resolve project root: %w", err)
	}

	session := &Session{
		ID:            filepath.Base(rootPath),
		ContainerName: containerName,
		LocalPath:     rootPath,
	}

	// 1. Clear container workspace if needed using ClearWorkspace(ctx, containerName)
	// 2. Copy project root into container volume using CopyToWorkspace(ctx, containerName, rootPath)

	return session, nil
}
