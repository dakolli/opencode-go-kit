package cli

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Container represents standard metadata for an open-code/agent container in Docker.
type Container struct {
	ID     string
	Name   string
	Status string
	Port   string
}
type fetchedContainersMsg []Container
type errMsg struct{ err error }

func fetchContainersCmd() tea.Msg {
	// Call your existing Docker logic here
	containers, err := ListOpenCodeContainers(context.Background())
	if err != nil {
		return errMsg{err: err}
	}
	return fetchedContainersMsg(containers)
}

// ListOpenCodeContainers lists all running or stopped containers that contain "opencode" in their name.
func ListOpenCodeContainers(ctx context.Context) ([]Container, error) {
	// 1. Added |{{.Ports}} to the end of the format string
	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--filter", "name=opencode", "--format", "{{.ID}}|{{.Names}}|{{.Status}}|{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var containers []Container
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 4 { // Look for 4 parts now
			continue
		}

		containers = append(containers, Container{
			ID:     parts[0],
			Name:   parts[1],
			Status: parts[2],
			Port:   parseExternalPort(parts[3]), // Clean up the messy port string
		})
	}
	return containers, nil
}

// CopyToWorkspace copies a local host directory into the container's /workspace folder.
func CopyToWorkspace(
	ctx context.Context,
	containerName string,
	hostPath string,
) error {
	// Using the docker cp command: docker cp /host/path/. <container>:/workspace
	// The '/.' suffix is used to copy the contents of the directory, not the directory itself.
	cmd := exec.CommandContext(ctx, "docker", "cp", hostPath+"/.", containerName+":/workspace")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy files to workspace: %w", err)
	}
	return nil
}

// CopyFromWorkspace copies a folder or file out of the container's workspace onto the host.
func CopyFromWorkspace(
	ctx context.Context,
	containerName string,
	containerPath string,
	hostDestPath string,
) error {
	// Example: docker cp <container>:/workspace/some_results /host/dest
	cmd := exec.CommandContext(ctx, "docker", "cp", containerName+":"+containerPath, hostDestPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy files from workspace: %w", err)
	}
	return nil
}

// ClearWorkspace removes all files inside the container's /workspace directory while preserving the volume.
func ClearWorkspace(
	ctx context.Context,
	containerName string,
) error {
	// Example: docker exec <container> sh -c "rm -rf /workspace/*"
	cmd := exec.CommandContext(ctx, "docker", "exec", containerName, "sh", "-c", "rm -rf /workspace/*")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clear workspace: %w", err)
	}
	return nil
}

// DeleteProjectFromWorkspace deletes a specific directory /workspace/<projectName>
func DeleteProjectFromWorkspace(
	ctx context.Context,
	containerName string,
	projectName string,
) error {
	// Example: docker exec <container> rm -rf /workspace/<projectName>
	targetDir := fmt.Sprintf("/workspace/%s", projectName)
	cmd := exec.CommandContext(ctx, "docker", "exec", containerName, "rm", "-rf", targetDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete project %s: %w", projectName, err)
	}
	return nil
}

var portRegex = regexp.MustCompile(`:(\d+)->`)

// parseExternalPort turns "0.0.0.0:8080->80/tcp" into "8080"
func parseExternalPort(dockerPorts string) string {
	if dockerPorts == "" {
		return "none"
	}

	matches := portRegex.FindStringSubmatch(dockerPorts)
	if len(matches) > 1 {
		return matches[1] // Returns just the captured digits (e.g., "8080")
	}

	return dockerPorts // Fallback to raw string if it's formatted unexpectedly
}
