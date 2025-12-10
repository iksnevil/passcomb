package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/iksnevil/passcomb/pkg/generator"
	"github.com/iksnevil/passcomb/pkg/tui"
)

type CLI struct {
	config generator.Config
	tui    bool
}

func NewCLI() *CLI {
	return &CLI{
		config: generator.Config{
			CombinationSize: 2,
			MaxFileSizeMB:   100,
		},
	}
}

func (c *CLI) ParseArgs(args []string) error {
	flags := flag.NewFlagSet("passcomb", flag.ExitOnError)

	var (
		inputFile       = flags.String("input", "", "Input file with passwords (one per line)")
		outputFile      = flags.String("output", "", "Output file for combinations")
		combinationSize = flags.Int("size", 2, "Combination size (2-4)")
		extraSymbols    = flags.String("symbols", "", "Extra symbols to use (e.g., '!@#$')")
		positions       = flags.String("positions", "", "Symbol positions: start,end,between")
		maxFileSize     = flags.Int("maxsize", 100, "Max file size in MB")
		showHelp        = flags.Bool("help", false, "Show help")
		useTUI          = flags.Bool("tui", false, "Use terminal UI interface")
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *showHelp {
		c.showHelp()
		os.Exit(0)
	}

	c.tui = *useTUI

	if !c.tui {
		// CLI mode - validate required parameters
		if *inputFile == "" {
			return fmt.Errorf("input file is required in CLI mode")
		}
		if *outputFile == "" {
			return fmt.Errorf("output file is required in CLI mode")
		}

		c.config.InputFile = *inputFile
		c.config.OutputFile = *outputFile
		c.config.CombinationSize = *combinationSize
		c.config.MaxFileSizeMB = *maxFileSize

		// Parse extra symbols
		if *extraSymbols != "" {
			c.config.ExtraSymbols = []rune(*extraSymbols)
		}

		// Parse symbol positions
		if *positions != "" {
			posList := strings.Split(*positions, ",")
			for _, pos := range posList {
				switch strings.TrimSpace(pos) {
				case "start":
					c.config.SymbolPositions = append(c.config.SymbolPositions, generator.PositionStart)
				case "end":
					c.config.SymbolPositions = append(c.config.SymbolPositions, generator.PositionEnd)
				case "between":
					c.config.SymbolPositions = append(c.config.SymbolPositions, generator.PositionBetween)
				default:
					return fmt.Errorf("invalid position: %s (valid: start, end, between)", pos)
				}
			}
		}

		// Validate combination size
		if c.config.CombinationSize < 2 || c.config.CombinationSize > 4 {
			return fmt.Errorf("combination size must be between 2 and 4")
		}
	}

	return nil
}

func (c *CLI) Run() error {
	if c.tui {
		// Run TUI mode
		program := tui.NewProgram()
		return program.Start()
	}

	// Run CLI mode
	return c.runCLI()
}

func (c *CLI) runCLI() error {
	fmt.Printf("Password Combination Generator\n")
	fmt.Printf("===============================\n\n")

	// Create generator
	gen := generator.NewGenerator(c.config)

	// Load passwords
	fmt.Printf("Loading passwords from: %s\n", c.config.InputFile)
	if err := gen.LoadPasswords(); err != nil {
		return fmt.Errorf("failed to load passwords: %w", err)
	}

	passwordCount := gen.GetPasswordCount()
	fmt.Printf("Loaded %d passwords\n", passwordCount)

	// Calculate combinations
	totalCombinations := gen.CalculateTotalCombinations()
	fmt.Printf("Total combinations to generate: %d\n", totalCombinations)

	// Show configuration
	fmt.Printf("\nConfiguration:\n")
	fmt.Printf("  Combination size: %d\n", c.config.CombinationSize)
	if len(c.config.ExtraSymbols) > 0 {
		fmt.Printf("  Extra symbols: %s\n", string(c.config.ExtraSymbols))
		var positions []string
		for _, pos := range c.config.SymbolPositions {
			switch pos {
			case generator.PositionStart:
				positions = append(positions, "start")
			case generator.PositionEnd:
				positions = append(positions, "end")
			case generator.PositionBetween:
				positions = append(positions, "between")
			}
		}
		fmt.Printf("  Symbol positions: %s\n", strings.Join(positions, ", "))
	}
	fmt.Printf("  Max file size: %d MB\n", c.config.MaxFileSizeMB)

	// Generate combinations
	fmt.Printf("\nGenerating combinations...\n")

	progressChan := make(chan generator.ProgressInfo)
	go func() {
		defer close(progressChan)
		if err := gen.GenerateCombinations(progressChan); err != nil {
			fmt.Printf("Error during generation: %v\n", err)
		}
	}()

	// Simple progress display
	for progress := range progressChan {
		percent := float64(progress.Generated) / float64(progress.TotalCombinations) * 100
		fmt.Printf("\rProgress: %.1f%% (%d/%d) - File: %s",
			percent, progress.Generated, progress.TotalCombinations, progress.CurrentFile)
	}

	fmt.Printf("\n\nGeneration complete!\n")
	return nil
}

func (c *CLI) showHelp() {
	fmt.Printf(`passcomb - Password Combination Generator

USAGE:
    passcomb [OPTIONS]

DESCRIPTION:
    Generates password combinations from a list of base passwords. Creates all possible
    combinations of specified length with optional extra symbols.

MODES:
    TUI Mode:    passcomb -tui
    CLI Mode:    passcomb -input <file> -output <file> [options]

TUI OPTIONS:
    -tui              Use terminal user interface (interactive mode)

CLI OPTIONS:
    -input string     Input file with passwords (one per line) [required in CLI mode]
    -output string    Output file for combinations [required in CLI mode]
    -size int         Combination size (2-4) [default: 2]
    -symbols string   Extra symbols to use (e.g., '!@#$') [default: none]
    -positions string Symbol positions: start,end,between [default: none]
    -maxsize int      Max file size in MB [default: 100]
    -help             Show this help message

EXAMPLES:
    # Interactive mode
    passcomb -tui

    # CLI mode - basic combinations
    passcomb -input passwords.txt -output combinations.txt -size 2

    # CLI mode - with extra symbols
    passcomb -input passwords.txt -output combos.txt -size 3 -symbols '!@#' -positions start,end

    # CLI mode - large output with file splitting
    passcomb -input passwords.txt -output combos.txt -size 4 -maxsize 50

SYMBOL POSITIONS:
    start     Add symbols at the beginning of combinations
    end       Add symbols at the end of combinations  
    between   Add symbols between password parts

FILE SIZE LIMITING:
    When the output exceeds the max file size, new files are created with numeric suffixes:
    combos.txt, combos_1.txt, combos_2.txt, etc.

EXIT CODES:
    0    Success
    1    Error
`)
}
