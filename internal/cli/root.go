package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kitty [path]",
	Short: "kitty is a developer utility for OpenCode workspace isolation",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		var targetDir string
		if len(args) > 0 {
			targetDir = args[0] // e.g., "."
		}

		// 1. Fetch available containers first
		// containers, err := ListOpenCodeContainers(ctx)
		// if err != nil {
		// 	return err
		// }

		// 2. Start the TUI, passing along the target directory context if it exists
		m := initialModel(targetDir)
		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("failed to run TUI: %w", err)
		}

		resModel, ok := finalModel.(tuiModel)
		if !ok || resModel.selectedContainer == "" || resModel.selectedAction == "" {
			fmt.Println("Session picker cancelled or no selection made.")
			return nil
		}

		// 3. Handle execution based on what the TUI resolved
		switch resModel.selectedAction {
		case "init":
			// If targetDir was passed to the TUI, use it; otherwise use cwd
			activePath := cwd
			if resModel.presetTargetDir != "" {
				activePath = resModel.presetTargetDir
			}

			fmt.Printf("Initializing workspace of '%s' into container '%s'...\n", activePath, resModel.selectedContainer)
			session, err := InitProjectSession(ctx, activePath, resModel.selectedContainer)
			if err != nil {
				return fmt.Errorf("initialization failed: %w", err)
			}
			fmt.Printf("Session initialized successfully! Session ID: %s\n", session.ID)

		case "attach":
			fmt.Printf("Attaching to container %s...\n", resModel.selectedContainer)
			instruction := "opencode; sh"
			execCmd := exec.Command("docker", "exec", "-it", resModel.selectedContainer, "sh", "-c", instruction)
			execCmd.Stdin = os.Stdin
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			if err := execCmd.Run(); err != nil {
				return fmt.Errorf("failed to attach to container terminal: %w", err)
			}
		}

		return nil
	},
}

// Global flags variables
var Verbose bool

func init() {
	// Define persistent flags (available to this command and all subcommands)
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose logging output")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
