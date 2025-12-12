package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/iksnevil/passcomb/pkg/generator"
)

type Model struct {
	config generator.Config
}

func NewModel() *Model {
	return &Model{
		config: generator.Config{
			CombinationSize: 2,
			MaxFileSizeMB:   100,
		},
	}
}

func (m *Model) Start() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter input file path: ")
	input, _ := reader.ReadString('\n')
	m.config.InputFile = strings.TrimSpace(input)

	fmt.Print("Enter output file path: ")
	output, _ := reader.ReadString('\n')
	m.config.OutputFile = strings.TrimSpace(output)

	fmt.Print("Enter combination size (2-4): ")
	sizeStr, _ := reader.ReadString('\n')
	sizeStr = strings.TrimSpace(sizeStr)

	return m.startGeneration()
}

func (m *Model) startGeneration() error {
	fmt.Println("\nStarting generation...")

	gen := generator.NewGenerator(m.config)

	if err := gen.LoadPasswords(); err != nil {
		return fmt.Errorf("failed to load passwords: %w", err)
	}

	passwordCount := gen.GetPasswordCount()
	totalCombinations := gen.CalculateTotalCombinations()

	fmt.Printf("Loaded %d passwords\n", passwordCount)
	fmt.Printf("Total combinations to generate: %d\n", totalCombinations)

	progressChan := make(chan generator.ProgressInfo)
	go func() {
		defer close(progressChan)
		if err := gen.GenerateCombinations(progressChan); err != nil {
			fmt.Printf("Generation error: %v\n", err)
		}
	}()

	for progress := range progressChan {
		percent := float64(progress.Generated) / float64(progress.TotalCombinations) * 100
		fmt.Printf("\rProgress: %.1f%% (%d/%d) - File: %s",
			percent, progress.Generated, progress.TotalCombinations, progress.CurrentFile)
	}

	fmt.Printf("\n\nGeneration complete!\n")
	return nil
}
