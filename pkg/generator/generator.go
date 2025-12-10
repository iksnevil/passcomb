package generator

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	InputFile       string
	OutputFile      string
	CombinationSize int
	ExtraSymbols    []rune
	SymbolPositions []SymbolPosition
	MaxFileSizeMB   int
}

type SymbolPosition int

const (
	PositionNone SymbolPosition = iota
	PositionStart
	PositionEnd
	PositionBetween
)

type Generator struct {
	config    Config
	passwords []string
}

type ProgressInfo struct {
	TotalCombinations int64
	Generated         int64
	CurrentFile       string
	FileNumber        int
}

func NewGenerator(config Config) *Generator {
	return &Generator{config: config}
}

func (g *Generator) LoadPasswords() error {
	file, err := os.Open(g.config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	var passwords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if password != "" {
			passwords = append(passwords, password)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	g.passwords = passwords
	return nil
}

func (g *Generator) CalculateTotalCombinations() int64 {
	if len(g.passwords) == 0 {
		return 0
	}

	baseCombinations := int64(math.Pow(float64(len(g.passwords)), float64(g.config.CombinationSize)))

	symbolMultiplier := 1
	if len(g.config.ExtraSymbols) > 0 && len(g.config.SymbolPositions) > 0 {
		symbolMultiplier = 1 + len(g.config.ExtraSymbols)*len(g.config.SymbolPositions)
	}

	return baseCombinations * int64(symbolMultiplier)
}

func (g *Generator) GenerateCombinations(progressChan chan<- ProgressInfo) error {
	if len(g.passwords) == 0 {
		return fmt.Errorf("no passwords loaded")
	}

	totalCombinations := g.CalculateTotalCombinations()
	if totalCombinations == 0 {
		return fmt.Errorf("no combinations to generate")
	}

	baseDir := filepath.Dir(g.config.OutputFile)
	baseName := filepath.Base(g.config.OutputFile)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	maxFileSize := int64(g.config.MaxFileSizeMB) * 1024 * 1024
	if maxFileSize <= 0 {
		maxFileSize = 100 * 1024 * 1024 // Default 100MB
	}

	var currentFile *os.File
	var currentFileSize int64
	fileNumber := 1
	generated := int64(0)

	createNewFile := func(num int) (*os.File, error) {
		if currentFile != nil {
			currentFile.Close()
		}

		var fileName string
		if num == 1 {
			fileName = g.config.OutputFile
		} else {
			fileName = filepath.Join(baseDir, fmt.Sprintf("%s_%d%s", nameWithoutExt, num, ext))
		}

		file, err := os.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}

		return file, nil
	}

	var err error
	currentFile, err = createNewFile(fileNumber)
	if err != nil {
		return err
	}
	defer currentFile.Close()

	writeCombination := func(combination string) error {
		combinationBytes := []byte(combination + "\n")

		if currentFileSize+int64(len(combinationBytes)) > maxFileSize {
			fileNumber++
			currentFile, err = createNewFile(fileNumber)
			if err != nil {
				return err
			}
			currentFileSize = 0
		}

		n, err := currentFile.Write(combinationBytes)
		if err != nil {
			return fmt.Errorf("failed to write combination: %w", err)
		}

		currentFileSize += int64(n)
		generated++

		progressChan <- ProgressInfo{
			TotalCombinations: totalCombinations,
			Generated:         generated,
			CurrentFile:       currentFile.Name(),
			FileNumber:        fileNumber,
		}

		return nil
	}

	// Generate base combinations
	g.generateBaseCombinations(writeCombination)

	// Generate combinations with extra symbols
	if len(g.config.ExtraSymbols) > 0 && len(g.config.SymbolPositions) > 0 {
		g.generateSymbolCombinations(writeCombination)
	}

	return nil
}

func (g *Generator) generateBaseCombinations(writeFunc func(string) error) {
	indices := make([]int, g.config.CombinationSize)
	for i := range indices {
		indices[i] = 0
	}

	for {
		// Build current combination
		var combination strings.Builder
		for i := 0; i < g.config.CombinationSize; i++ {
			combination.WriteString(g.passwords[indices[i]])
		}

		if err := writeFunc(combination.String()); err != nil {
			return
		}

		// Move to next combination
		carry := 1
		for i := g.config.CombinationSize - 1; i >= 0 && carry > 0; i-- {
			indices[i] += carry
			if indices[i] >= len(g.passwords) {
				indices[i] = 0
				carry = 1
			} else {
				carry = 0
			}
		}

		if carry > 0 {
			break // All combinations generated
		}
	}
}

func (g *Generator) generateSymbolCombinations(writeFunc func(string) error) {
	indices := make([]int, g.config.CombinationSize)
	for i := range indices {
		indices[i] = 0
	}

	for {
		// Build base combination
		var baseCombination strings.Builder
		for i := 0; i < g.config.CombinationSize; i++ {
			baseCombination.WriteString(g.passwords[indices[i]])
		}
		base := baseCombination.String()

		// Generate combinations with symbols
		for _, symbol := range g.config.ExtraSymbols {
			for _, position := range g.config.SymbolPositions {
				var combination string
				switch position {
				case PositionStart:
					combination = string(symbol) + base
				case PositionEnd:
					combination = base + string(symbol)
				case PositionBetween:
					if g.config.CombinationSize > 1 {
						parts := make([]string, g.config.CombinationSize)
						for i := 0; i < g.config.CombinationSize; i++ {
							parts[i] = g.passwords[indices[i]]
						}
						combination = strings.Join(parts[:len(parts)-1], "") + string(symbol) + parts[len(parts)-1]
					} else {
						combination = base + string(symbol) // For single password, treat as end
					}
				}

				if err := writeFunc(combination); err != nil {
					return
				}
			}
		}

		// Move to next combination
		carry := 1
		for i := g.config.CombinationSize - 1; i >= 0 && carry > 0; i-- {
			indices[i] += carry
			if indices[i] >= len(g.passwords) {
				indices[i] = 0
				carry = 1
			} else {
				carry = 0
			}
		}

		if carry > 0 {
			break // All combinations generated
		}
	}
}

func (g *Generator) GetPasswordCount() int {
	return len(g.passwords)
}
