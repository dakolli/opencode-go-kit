package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tuiState int

const (
	stateSelectContainer tuiState = iota
	stateSelectAction
	stateSelectTemplate
	stateDone
)

type tuiModel struct {
	state             tuiState
	containers        []Container
	templates         []Template
	selectedContainer string
	selectedTemplate  string
	selectedAction    string // "init" or "attach"
	cursor            int
	actionCursor      int
	templateCursor    int
	err               error

	// --- NEW FIELD ---
	isLoading       bool
	presetTargetDir string // Stores the path context if 'kitty .' was used
}

// Updated constructor to accept the preset string
func initialModel(presetTargetDir string) tuiModel {
	return tuiModel{
		state: stateSelectContainer,
		// containers:      containers,
		presetTargetDir: presetTargetDir,
		isLoading:       true,
	}
}

func (m tuiModel) Init() tea.Cmd {
	// tea.Batch takes multiple commands and runs them concurrently in the background
	return tea.Batch(
		fetchContainersCmd,
		fetchTemplatesCmd,
	)
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchedContainersMsg:
		m.containers = msg
		return m, nil

	// 2. Templates finished fetching!
	case fetchedTemplatesMsg:
		// Convert []Template into []string if your existing tuiModel tracking relies on strings,
		// or just update your model to store []Template directly.
		var names []Template
		for _, t := range msg {
			names = append(names, newTemplate(t.Name))
		}
		m.templates = names
		return m, nil

	// 3. Something went wrong in either of them
	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			switch m.state {
			case stateSelectContainer:
				if m.cursor > 0 {
					m.cursor--
				}
			case stateSelectAction:
				if m.actionCursor > 0 {
					m.actionCursor--
				}
			}

		case "down", "j":
			switch m.state {
			case stateSelectContainer:
				if m.cursor < len(m.containers)-1 {
					m.cursor++
				}
			case stateSelectAction:
				if m.actionCursor < 1 {
					m.actionCursor++
				}
			}

		case "enter":
			switch m.state {
			case stateSelectContainer:
				if len(m.containers) > 0 {
					m.selectedContainer = m.containers[m.cursor].Name

					// --- FAST-TRACK BORDER INTERCEPT ---
					if m.presetTargetDir != "" {
						m.selectedAction = "init"
						// Instead of quitting, skip right over 'stateSelectAction'
						// and head straight into the template workspace room
						m.state = stateSelectTemplate
						return m, nil
					}

					// Standard route: go to action menu picker
					m.state = stateSelectAction
				}

			case stateSelectAction:
				if m.actionCursor == 0 {
					m.selectedAction = "init"
					// User chose "init", so route them down the template choosing corridor
					m.state = stateSelectTemplate
					return m, nil
				} else {
					m.selectedAction = "attach"
					m.state = stateDone
					return m, tea.Quit // "attach" is finished immediately, exit to Cobra
				}

			case stateSelectTemplate:
				if len(m.templates) > 0 {
					m.selectedTemplate = m.templates[m.templateCursor].Name
				}
				// All wizard choices are fully locked in! Exit back to Cobra
				m.state = stateDone
				return m, tea.Quit
			}

		case "esc":
			// Only allow backing up if they weren't locked into a fast-track selection
			if m.state == stateSelectAction && m.presetTargetDir == "" {
				m.state = stateSelectContainer
			}
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	kStyled := pinkStyle.Render(letterK)
	iStyled := blueStyle.Render(letterI)
	t1Styled := pinkStyle.Render(letterT)
	t2Styled := blueStyle.Render(letterT)
	yStyled := pinkStyle.Render(letterY)

	// Join with slight horizontal offset for readability
	kittyTitle := lipgloss.JoinHorizontal(lipgloss.Bottom, kStyled, "  ", iStyled, "  ", t1Styled, "  ", t2Styled, "  ", yStyled)

	// 1. Build the menu content starting with the Kitty title
	menuContent := kittyTitle + "\n\n"

	switch m.state {
	case stateSelectContainer:
		if m.presetTargetDir != "" {
			menuContent += fmt.Sprintf("=== Target Workspace: %s ===\n", m.presetTargetDir)
			menuContent += "Select a destination OpenCode Container:\n\n"
		} else {
			menuContent += "=== Select Container ===\n\n"
		}

		if len(m.containers) == 0 {
			menuContent += "No matching containers found.\nPress q to exit.\n"
		} else {
			for i, c := range m.containers {
				cursor := " "
				if m.cursor == i {
					cursor = ">"
				}
				menuContent += fmt.Sprintf("%s [%s] %s\n", cursor, c.Status, c.Name)
			}
			menuContent += "\nUse arrows / jk to move, enter to select.\nPress q to quit.\n"
		}

	case stateSelectAction:
		menuContent += fmt.Sprintf("=== Select Action for [%s] ===\n\n", m.selectedContainer)
		actions := []string{
			"Initialize Session (Copy local project into container)",
			"Attach to Shell (Run opencode TUI inside container)",
		}
		for i, act := range actions {
			cursor := " "
			if m.actionCursor == i {
				cursor = ">"
			}
			menuContent += fmt.Sprintf("%s %s\n", cursor, act)
		}
		menuContent += "\nEsc to go back. Press q to quit.\n"

	case stateSelectTemplate:
		menuContent += "=== Select a Workspace Template ===\n\n"
		for i, t := range m.templates {
			cursor := " "
			if m.templateCursor == i {
				cursor = ">"
			}
			menuContent += fmt.Sprintf("%s %s\n", cursor, t)
		}
		menuContent += "\nEnter to finalize workspace setup.\n"

	case stateDone:
		menuContent += "Launching Container...\n"
	}

	// 2. Apply Lip Gloss styles to your blocks to turn them into official renderable columns
	leftColumn := asciiColumnStyle.Render(kittyASCII)
	rightColumn := menuColumnStyle.Render(menuContent)

	// 3. Join them horizontally!
	// lipgloss.Center aligns the centers of the columns vertically.
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Center, leftColumn, rightColumn) + "\n"
}

var (
	asciiColumnStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				PaddingRight(4).
				Foreground(lipgloss.Color("#F780E2")) // Giving the kitty a nice pink tint

	menuColumnStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	pinkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F780E2"))
	blueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#5EA5FF"))
)

const (
	letterK = ` _  __
| |/ /
| ' / 
| . \ 
|_|\_\`

	letterI = ` _ 
 (_)
 | |
 | |
 |_|`

	letterT = ` _  
| |_
| __|
| |_
 \__|`

	letterY = ` _   _
| | / /
| |/ / 
\_  /  
 /_/   `
)

// Your ASCII art string (Make sure to escape backslashes if needed!)
const kittyASCII = `
⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢀⣾⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⣸⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣠⣤⣤⣼⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀
⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀
⠀⠀⠀⠘⣿⣿⣿⣿⠟⠁⠀⠀⠀⠹⣿⣿⣿⣿⣿⠟⠁⠀⠀⠹⣿⣿⡿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣿⡇⠀⠀⠀⢼⣿⠀⢿⣿⣿⣿⣿⠀⣾⣷⠀⠀⢿⣿⣷⠀⠀⠀⠀⠀
⠀⠀⠀⢠⣿⣿⣿⣷⡀⠀⠀⠈⠋⢀⣿⣿⣿⣿⣿⡀⠙⠋⠀⢀⣾⣿⣿⠀⠀⠀⠀⠀
⢀⣀⣀⣀⣿⣿⣿⣿⣿⣶⣶⣶⣶⣿⣿⣿⣿⣾⣿⣷⣦⣤⣴⣿⣿⣿⣿⣤⠤⢤⣤⡄
⠈⠉⠉⢉⣙⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⣀⣀⣀⡀⠀
⠐⠚⠋⠉⢀⣬⡿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⣥⣀⡀⠈⠀⠈⠛
⠀⠀⠴⠚⠉⠀⠀⠀⠉⠛⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠋⠁⠀⠀⠀⠉⠛⠢⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`

// ⣿⣿⣿⢾⡙⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⠟⣫⠟⣙⠛⢛
// ⣿⣿⣿⣾⣿⣄⣠⣼⠟⠙⠷⢶⣄⠈⠐⢾⡷⣙⣿⡀⣸
// ⣿⣿⣿⣿⣯⢟⣀⢌⠀⠰⡀⣄⠘⢿⣦⠀⠹⣾⣿⣢⡇
// ⣿⣿⣿⡿⠁⣸⣸⢘⠀⢦⢷⣘⣦⣞⢿⣞⡆⢹⢿⣾⣷
// ⣿⣿⣻⠃⠀⣿⢿⣼⣀⢹⣾⡿⣯⡿⣿⣻⢹⢹⣿⣋⣎⣢
// ⣯⣾⡟⠃⣶⣿⢉⣻⣿⣿⣻⣾⢾⢿⣿⢻⡞⢀⣿⣿⣷⣿
// ⣿⣿⣷⣧⢿⣿⡿⢿⡇⠙⠛⠊⠀⣛⣛⣈⡃⢸⣿⣿⢿⢻
// ⣿⣿⣿⣾⣼⣿⣗⣼⠁⠄⠀⠀⠀⠘⠉⣹⠁⣬⣿⡿⡞⡆
// ⣿⣿⣿⣿⣿⣿⢿⠉⠀⠀⠀⠀⠀⠀⣨⣿⣐⣿⣿⣿⢹⣃
// ⣿⣿⣿⣿⣿⣿⢸⢹⣲⣦⣤⣠⣴⣾⡿⢺⢿⣿⣿⢽⢿⣿

// ⣿⣿⣿⣿⣿⣿⣿⣤⡘⡍⠛⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠛⠋⡡⠊⠉⢀⠞⡻⠊⡐⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⡟⢿⠻⡀⠀⠀⠉⠛⠿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠉⠀⠀⣠⠞⠃⠀⢀⡀⠀⠀⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⡄⠀⢀⠳⡀⠀⠀⠀⠀⠀⠈⠙⢿⣿⣿⡿⠟⣻⠟⠉⠀⠀⠀⠀⠀⠀⠠⠟⠁⠀⠀⠀⣠⡞⠀⠀⡀⠠⢜⣡⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⡇⠂⣰⣧⠐⡄⠀⠀⠀⠀⠀⣐⣮⣤⠤⠀⠨⠴⣶⡶⡶⠤⣀⠀⠀⠀⠀⠀⠀⠑⠲⣿⣋⣭⣛⠁⠠⠒⢵⣏⡀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣹⣄⢘⣧⠀⢺⣆⠀⣠⠴⠛⠉⡱⠁⠀⠀⠀⠀⠀⠀⠀⠑⠢⡹⣦⡀⠀⠀⠀⠀⠀⠈⠻⡗⠀⠀⣀⣀⣒⡓⠧⡀⠀⢠⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⡈⣙⡿⠂⢈⡿⠋⠀⠀⠀⠈⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠙⢝⣆⠑⡄⠀⠀⠀⠀⠘⢶⣀⣈⠉⣻⢿⡧⠀⠀⣸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣿⣀⠞⠀⠀⡀⠀⠠⠂⠀⠀⠀⠀⠀⠱⠀⠀⠠⠀⠀⠀⠀⠀⠻⡱⣼⡄⠀⠀⠀⠀⠈⢗⠤⣅⣴⣟⣇⣠⠀⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⡿⠁⠀⠀⡰⠀⠀⡆⢰⠀⠀⠀⠀⠀⠀⢇⠀⠀⠑⡀⠀⠀⠀⠠⡘⢜⢿⡄⠀⢀⠀⠀⠈⢇⢨⣿⣿⠃⢀⣴⡅⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀⠀⠀⢠⢃⠀⣰⠀⠘⠀⠀⠀⠀⣇⠀⢸⣧⡀⠀⠘⢦⠀⠠⠀⠉⢌⠺⣷⠀⠈⡆⠀⠀⠘⣿⢟⡿⢶⣿⡟⠘⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⡿⠃⠀⠀⠀⠀⣤⠇⢠⡟⡆⠀⡄⠀⠀⠀⠘⡄⠀⢏⠓⢄⠀⠣⠳⣄⠑⢄⠀⠣⡹⣇⠀⢰⠀⠀⠀⢳⠈⣰⣿⣿⢅⠘⣱⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⠏⢠⠃⠀⠀⠀⢠⡟⠀⡾⠀⡇⠀⠂⠀⠀⠀⢳⠸⡄⠘⣄⣧⡦⡔⠒⠘⢯⡙⢍⠒⠚⢾⡀⠐⡆⠀⢡⠸⡮⠿⢫⡇⠘⡄⠀⢣⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⠟⠁⢠⠃⠀⠀⠀⠀⢸⠇⣸⠇⢠⢡⢀⢸⡀⠀⠀⠸⣧⢳⡠⢻⠀⠑⢜⣆⠘⠄⠑⢌⡳⡀⠀⠇⠀⣧⠀⢸⠀⣇⡰⠋⠀⠀⣷⠀⡀⠡⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⠟⠁⢀⣴⡏⠀⡄⠀⠀⠀⢸⢀⣿⡀⠤⢼⡞⡄⣷⡘⣆⠀⠹⣍⣷⡀⢣⠀⠀⣙⣻⣦⣤⣤⣭⣛⣲⣤⠀⢹⠀⠀⠀⣍⣠⣤⠁⡰⢿⠀⠰⡀⠱⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠋⢀⣠⣴⣿⣿⠁⠀⠁⠀⣸⠀⣾⠜⡏⠀⠀⠈⢿⡰⡘⣷⠘⣧⡀⠻⣿⣧⣀⣶⡿⢿⣿⣿⣿⣿⣿⡝⠻⣿⡶⠏⠀⠀⠀⢹⣿⣁⠜⡀⡿⣧⠀⡙⡄⠈⡢⡀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⢿⠀⣿⠀⠃⠀⢀⣀⣨⠷⡑⢼⣧⡽⣷⣄⡀⢻⡠⡏⠀⠸⠏⢹⣿⣿⠙⡇⠀⣼⡇⠀⠀⠀⡆⣎⡙⢇⡰⢁⡇⣿⣧⣿⢮⠢⡀⠀⠑⠢⡀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⡇⠀⠀⡇⠀⢸⡇⢻⠀⢐⣾⠟⣿⣿⣿⣗⠀⠈⠻⣆⠝⠣⢄⠈⠀⠀⠀⠰⣟⠹⢻⣿⠁⠀⠃⠁⠀⠀⢸⡇⣃⡁⢸⠁⣼⢧⢸⠙⡟⢷⣕⡌⠒⠤⣀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⡇⢸⠀⣧⠀⠸⣷⠘⣄⣾⠇⠀⠿⠛⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣈⣁⠉⣁⣀⣀⢀⡁⠀⠀⢸⡄⡏⣡⢟⣿⣿⠸⠀⠄⠸⠀⠀⠀⠉⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⢻⣿⣿⡇⢸⡆⢹⡄⠀⢹⣧⣻⡻⣇⠀⠀⢰⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⠟⢁⠠⠀⢸⠇⠀⠀⠸⠃⡗⣡⣾⡏⢿⠀⡇⡇⠀⠑⢄⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣸⣿⣿⡇⢸⣷⡘⣷⡀⠈⣿⣿⣧⠈⠂⠀⠀⠉⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠁⠀⠀⠀⠀⣾⠀⠀⠀⣀⠀⣿⣿⢿⡇⢸⡄⠀⢩⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⢸⣿⣿⣿⠘⣿⣷⣹⣷⣄⠘⣎⢻⣇⠀⠔⢫⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡔⢸⠀⠀⠀⣿⢠⣿⡇⢸⡇⠀⢷⠀⠸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣇⢿⣿⣿⣿⣿⡿⣏⠙⢿⠘⠒⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠊⢀⣼⠀⡄⢠⡏⢸⣿⣧⠘⣿⠀⠈⣇⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣾⣿⣿⣿⡿⢁⣿⠀⠸⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣴⠟⣿⠀⠁⢸⡁⣼⣿⣿⠀⡟⡇⠀⠼⡀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣵⣿⢟⡄⠀⢛⠢⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⣾⠟⠁⢀⡏⠸⠀⣇⡆⣿⣿⣿⣇⣧⣿⡤⣶⣷⣸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢸⣇⠀⠸⡀⢸⡏⠓⣦⣤⣀⣀⠀⠀⠀⠀⠀⣀⣤⣾⣿⠟⠁⣠⠐⠁⠃⡆⢠⣾⢱⣿⡏⢉⣀⣠⣴⢤⣤⣬⣁⡒⠢⠤⣄⣀⡀⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡏⣼⣿⠀⠀⠃⠘⡇⢀⣿⢸⣿⣿⣿⣿⣷⣶⣿⣿⣿⠟⠁⡠⠊⠀⠀⠀⣷⠁⢸⣟⣸⠿⠿⢧⡅⢱⣧⠀⡅⡇⠇⢳⠈⠁⠢⡀⠁⠀⠀⠀⠀⠀
// ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢃⣿⢸⡇⠀⠰⠀⣿⢸⣿⣿⣿⣿⣿⠟⢛⣩⣿⡏⠁⠔⠈⢀⣠⠄⠒⢸⡟⠀⣿⡇⡇⠀⠀⢸⠃⠘⣿⠀⡁⢁⢸⠈⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⢹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⣼⣿⢸⣿⡀⠀⠀⢸⣿⣿⣿⡿⠋⠵⡽⠿⠛⠛⠋⠉⠉⠉⠁⡁⠀⠀⣿⡇⢠⣿⣹⠁⠀⠀⣿⠀⢰⢿⠀⠙⠘⢸⡀⣷⠀⠀⠀⠀
