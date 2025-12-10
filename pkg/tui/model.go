package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateInputFile State = iota
	StateOutputFile
	StateCombinationSize
	StateExtraSymbols
	StateSymbolPositions
	StateMaxFileSize
	StateGenerating
	StateComplete
)

type Model struct {
	state             State
	inputFile         string
	outputFile        string
	combinationSize   int
	extraSymbols      []rune
	symbolPositions   []int
	maxFileSizeMB     int
	cursor            int
	selectedSymbols   map[int]bool
	selectedPositions map[int]bool
	inputBuffer       string
	progress          ProgressInfo
	error             string
	width             int
	height            int
}

type ProgressInfo struct {
	TotalCombinations int64
	Generated         int64
	CurrentFile       string
	FileNumber        int
	Percent           float64
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C"))

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4"))
)

func InitialModel() Model {
	return Model{
		state:             StateInputFile,
		combinationSize:   2,
		maxFileSizeMB:     100,
		selectedSymbols:   make(map[int]bool),
		selectedPositions: make(map[int]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateInputFile, StateOutputFile:
			return m.handleFileInput(msg)
		case StateCombinationSize:
			return m.handleCombinationSizeInput(msg)
		case StateExtraSymbols:
			return m.handleExtraSymbolsInput(msg)
		case StateSymbolPositions:
			return m.handleSymbolPositionsInput(msg)
		case StateMaxFileSize:
			return m.handleMaxFileSizeInput(msg)
		case StateGenerating:
			return m.handleGeneratingInput(msg)
		case StateComplete:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ProgressInfo:
		m.progress = msg
		m.progress.Percent = float64(msg.Generated) / float64(msg.TotalCombinations) * 100
	}

	return m, nil
}

func (m Model) handleFileInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.state == StateInputFile {
			if m.inputFile != "" {
				m.state = StateOutputFile
				m.inputBuffer = ""
			} else {
				m.error = "Input file cannot be empty"
			}
		} else if m.state == StateOutputFile {
			if m.outputFile != "" {
				m.state = StateCombinationSize
				m.inputBuffer = ""
			} else {
				m.error = "Output file cannot be empty"
			}
		}
	case tea.KeyBackspace:
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
	case tea.KeyTab:
		// TODO: Implement tab completion
	default:
		if len(msg.String()) == 1 {
			m.inputBuffer += msg.String()
			if m.state == StateInputFile {
				m.inputFile = m.inputBuffer
			} else {
				m.outputFile = m.inputBuffer
			}
		}
	}
	return m, nil
}

func (m Model) handleCombinationSizeInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.combinationSize >= 2 && m.combinationSize <= 4 {
			m.state = StateExtraSymbols
		} else {
			m.error = "Combination size must be between 2 and 4"
		}
	case tea.KeyUp:
		if m.combinationSize < 4 {
			m.combinationSize++
		}
	case tea.KeyDown:
		if m.combinationSize > 2 {
			m.combinationSize--
		}
	}
	return m, nil
}

func (m Model) handleExtraSymbolsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	symbols := []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')'}

	switch msg.Type {
	case tea.KeyEnter:
		m.state = StateSymbolPositions
		m.cursor = 0
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(symbols)-1 {
			m.cursor++
		}
	case tea.KeySpace:
		idx := m.cursor
		if m.selectedSymbols[idx] {
			delete(m.selectedSymbols, idx)
		} else {
			m.selectedSymbols[idx] = true
		}
	}
	return m, nil
}

func (m Model) handleSymbolPositionsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	positions := []string{"Start", "End", "Between"}

	switch msg.Type {
	case tea.KeyEnter:
		m.state = StateMaxFileSize
		m.cursor = 0
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(positions)-1 {
			m.cursor++
		}
	case tea.KeySpace:
		idx := m.cursor
		if m.selectedPositions[idx] {
			delete(m.selectedPositions, idx)
		} else {
			m.selectedPositions[idx] = true
		}
	}
	return m, nil
}

func (m Model) handleMaxFileSizeInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.maxFileSizeMB > 0 {
			m.state = StateGenerating
			return m, nil
		} else {
			m.error = "Max file size must be greater than 0"
		}
	case tea.KeyUp:
		m.maxFileSizeMB += 10
	case tea.KeyDown:
		if m.maxFileSizeMB > 10 {
			m.maxFileSizeMB -= 10
		}
	}
	return m, nil
}

func (m Model) handleGeneratingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		m.width = 80
	}

	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render(" Password Combination Generator "))
	content.WriteString("\n\n")

	// Error message
	if m.error != "" {
		content.WriteString(errorStyle.Render(" Error: " + m.error))
		content.WriteString("\n\n")
	}

	switch m.state {
	case StateInputFile:
		content.WriteString(m.renderFileInput("Input File", m.inputFile, m.inputBuffer))
	case StateOutputFile:
		content.WriteString(m.renderFileInput("Output File", m.outputFile, m.inputBuffer))
	case StateCombinationSize:
		content.WriteString(m.renderCombinationSize())
	case StateExtraSymbols:
		content.WriteString(m.renderExtraSymbols())
	case StateSymbolPositions:
		content.WriteString(m.renderSymbolPositions())
	case StateMaxFileSize:
		content.WriteString(m.renderMaxFileSize())
	case StateGenerating:
		content.WriteString(m.renderGenerating())
	case StateComplete:
		content.WriteString(m.renderComplete())
	}

	return content.String()
}

func (m Model) renderFileInput(title, value, buffer string) string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("%s:\n", title))
	content.WriteString(borderStyle.Render(fmt.Sprintf("📁 %s", value)))
	content.WriteString(fmt.Sprintf("\n\nCurrent input: %s", buffer))
	content.WriteString("\n\nEnter file path and press Enter")
	return content.String()
}

func (m Model) renderCombinationSize() string {
	var content strings.Builder
	content.WriteString("Combination Size:\n\n")

	for i := 2; i <= 4; i++ {
		prefix := " "
		if i == m.combinationSize {
			prefix = cursorStyle.Render("▶")
		}
		content.WriteString(fmt.Sprintf("%s %d passwords\n", prefix, i))
	}

	content.WriteString("\nUse ↑/↓ to change, Enter to continue")
	return content.String()
}

func (m Model) renderExtraSymbols() string {
	var content strings.Builder
	content.WriteString("Select Extra Symbols (Space to toggle):\n\n")

	symbols := []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')'}
	for i, symbol := range symbols {
		prefix := " "
		if i == m.cursor {
			prefix = cursorStyle.Render("▶")
		}

		checked := " "
		if m.selectedSymbols[i] {
			checked = "✓"
		}

		content.WriteString(fmt.Sprintf("%s [%s] %c\n", prefix, checked, symbol))
	}

	content.WriteString("\nPress Enter to continue")
	return content.String()
}

func (m Model) renderSymbolPositions() string {
	var content strings.Builder
	content.WriteString("Select Symbol Positions (Space to toggle):\n\n")

	positions := []string{"Start", "End", "Between"}
	for i, position := range positions {
		prefix := " "
		if i == m.cursor {
			prefix = cursorStyle.Render("▶")
		}

		checked := " "
		if m.selectedPositions[i] {
			checked = "✓"
		}

		content.WriteString(fmt.Sprintf("%s [%s] %s\n", prefix, checked, position))
	}

	content.WriteString("\nPress Enter to continue")
	return content.String()
}

func (m Model) renderMaxFileSize() string {
	var content strings.Builder
	content.WriteString("Max File Size (MB):\n\n")
	content.WriteString(fmt.Sprintf("%s %d MB", cursorStyle.Render("▶"), m.maxFileSizeMB))
	content.WriteString("\n\nUse ↑/↓ to change, Enter to continue")
	return content.String()
}

func (m Model) renderGenerating() string {
	var content strings.Builder
	content.WriteString("Generating combinations...\n\n")

	if m.progress.TotalCombinations > 0 {
		// Progress bar
		barWidth := 40
		filled := int(float64(barWidth) * m.progress.Percent / 100)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		content.WriteString(fmt.Sprintf("[%s] %.1f%%\n", bar, m.progress.Percent))

		content.WriteString(fmt.Sprintf("Generated: %d / %d\n", m.progress.Generated, m.progress.TotalCombinations))
		content.WriteString(fmt.Sprintf("Current file: %s\n", m.progress.CurrentFile))
		content.WriteString(fmt.Sprintf("File number: %d\n", m.progress.FileNumber))
	}

	content.WriteString("\nPress Ctrl+C to cancel")
	return content.String()
}

func (m Model) renderComplete() string {
	var content strings.Builder
	content.WriteString(successStyle.Render("✓ Generation Complete!"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("Total combinations generated: %d\n", m.progress.Generated))
	content.WriteString(fmt.Sprintf("Files created: %d\n", m.progress.FileNumber))
	content.WriteString("\nPress Enter to exit")
	return content.String()
}
