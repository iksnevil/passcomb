package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iksnevil/passcomb/pkg/generator"
)

type Program struct {
	model Model
}

func NewProgram() *Program {
	return &Program{
		model: InitialModel(),
	}
}

func (p *Program) Start() error {
	program := tea.NewProgram(p.model, tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}

	if model, ok := finalModel.(Model); ok {
		return p.handleCompletion(model)
	}

	return nil
}

func (p *Program) handleCompletion(model Model) error {
	if model.state == StateComplete {
		return nil
	}
	return nil
}

func StartGeneration(model Model) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// Build generator config
		var extraSymbols []rune
		symbols := []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')'}
		for i := range symbols {
			if model.selectedSymbols[i] {
				extraSymbols = append(extraSymbols, symbols[i])
			}
		}

		var symbolPositions []generator.SymbolPosition
		positions := []generator.SymbolPosition{
			generator.PositionStart,
			generator.PositionEnd,
			generator.PositionBetween,
		}
		for i := range positions {
			if model.selectedPositions[i] {
				symbolPositions = append(symbolPositions, positions[i])
			}
		}

		config := generator.Config{
			InputFile:       model.inputFile,
			OutputFile:      model.outputFile,
			CombinationSize: model.combinationSize,
			ExtraSymbols:    extraSymbols,
			SymbolPositions: symbolPositions,
			MaxFileSizeMB:   model.maxFileSizeMB,
		}

		gen := generator.NewGenerator(config)

		// Load passwords
		if err := gen.LoadPasswords(); err != nil {
			return errorMsg{err: err}
		}

		// Start generation with progress channel
		progressChan := make(chan generator.ProgressInfo)
		go func() {
			defer close(progressChan)
			if err := gen.GenerateCombinations(progressChan); err != nil {
				// Send error through channel or handle differently
			}
		}()

		// Return first progress message
		select {
		case progress := <-progressChan:
			return progressMsg{
				progress: generator.ProgressInfo{
					TotalCombinations: progress.TotalCombinations,
					Generated:         progress.Generated,
					CurrentFile:       progress.CurrentFile,
					FileNumber:        progress.FileNumber,
				},
			}
		default:
			return progressMsg{
				progress: generator.ProgressInfo{
					TotalCombinations: gen.CalculateTotalCombinations(),
					Generated:         0,
					CurrentFile:       "",
					FileNumber:        0,
				},
			}
		}
	})
}

type errorMsg struct {
	err error
}

type progressMsg struct {
	progress generator.ProgressInfo
}

func UpdateModelWithProgress(model Model, msg progressMsg) Model {
	model.progress = ProgressInfo{
		TotalCombinations: msg.progress.TotalCombinations,
		Generated:         msg.progress.Generated,
		CurrentFile:       msg.progress.CurrentFile,
		FileNumber:        msg.progress.FileNumber,
		Percent:           float64(msg.progress.Generated) / float64(msg.progress.TotalCombinations) * 100,
	}

	if model.progress.Generated >= model.progress.TotalCombinations {
		model.state = StateComplete
	}

	return model
}

func CompleteFilePath(path string) string {
	if path == "" {
		return ""
	}

	// If path ends with separator, complete directory
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return completeDirectory(path)
	}

	// Try to complete file/directory name
	matches, err := filepath.Glob(path + "*")
	if err != nil || len(matches) == 0 {
		return path
	}

	if len(matches) == 1 {
		if isDir(matches[0]) {
			return matches[0] + string(filepath.Separator)
		}
		return matches[0]
	}

	// Find common prefix
	common := findCommonPrefix(matches)
	if common != "" {
		return common
	}

	return path
}

func completeDirectory(path string) string {
	matches, err := filepath.Glob(path + "*")
	if err != nil || len(matches) == 0 {
		return path
	}

	if len(matches) == 1 {
		if isDir(matches[0]) {
			return matches[0] + string(filepath.Separator)
		}
		return matches[0]
	}

	return path
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func findCommonPrefix(matches []string) string {
	if len(matches) == 0 {
		return ""
	}

	prefix := matches[0]
	for _, match := range matches[1:] {
		prefix = commonPrefix(prefix, match)
		if prefix == "" {
			break
		}
	}

	return prefix
}

func commonPrefix(a, b string) string {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	i := 0
	for i < minLen && a[i] == b[i] {
		i++
	}

	return a[:i]
}
