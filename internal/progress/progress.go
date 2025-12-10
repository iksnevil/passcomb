package progress

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Total         int64
	Current       int64
	CurrentFile   string
	FileNumber    int
	Width         int
	Active        bool
	StartTime     time.Time
	EstimatedTime time.Duration
}

type ProgressMsg struct {
	Total       int64
	Current     int64
	CurrentFile string
	FileNumber  int
}

func NewModel() Model {
	return Model{
		Active:    false,
		StartTime: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		m.Total = msg.Total
		m.Current = msg.Current
		m.CurrentFile = msg.CurrentFile
		m.FileNumber = msg.FileNumber
		m.Active = true

		if m.Current > 0 {
			elapsed := time.Since(m.StartTime)
			rate := float64(m.Current) / elapsed.Seconds()
			remaining := float64(m.Total-m.Current) / rate
			m.EstimatedTime = time.Duration(remaining) * time.Second
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}

	return m, nil
}

func (m Model) View() string {
	if !m.Active || m.Total == 0 {
		return ""
	}

	var content strings.Builder

	// Progress bar
	barWidth := 40
	if m.Width > 0 && m.Width < 80 {
		barWidth = m.Width - 40
		if barWidth < 20 {
			barWidth = 20
		}
	}

	percent := float64(m.Current) / float64(m.Total) * 100
	filled := int(float64(barWidth) * percent / 100)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	content.WriteString(fmt.Sprintf("[%s] %.1f%%\n", bar, percent))

	// Stats
	content.WriteString(fmt.Sprintf("Generated: %d / %d\n", m.Current, m.Total))
	content.WriteString(fmt.Sprintf("Current file: %s\n", m.CurrentFile))
	content.WriteString(fmt.Sprintf("File number: %d\n", m.FileNumber))

	// Time estimation
	if m.EstimatedTime > 0 {
		content.WriteString(fmt.Sprintf("ETA: %v\n", m.EstimatedTime.Round(time.Second)))
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render(content.String())
}

func (m Model) IsComplete() bool {
	return m.Active && m.Current >= m.Total
}
