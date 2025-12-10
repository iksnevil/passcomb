package main

import (
	"fmt"
	"os"

	"github.com/iksnevil/passcomb/pkg/cli"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	app := cli.NewCLI()

	// If no arguments provided, default to TUI mode
	if len(os.Args) == 1 {
		return app.Run()
	}

	// Parse command line arguments
	if err := app.ParseArgs(os.Args[1:]); err != nil {
		return err
	}

	// Run the application
	return app.Run()
}
