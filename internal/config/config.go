package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultMaxFileSizeMB   = 100
	DefaultCombinationSize = 2
)

type AppConfig struct {
	InputFile       string
	OutputFile      string
	CombinationSize int
	ExtraSymbols    []rune
	SymbolPositions []string
	MaxFileSizeMB   int
}

func LoadDefaults() AppConfig {
	return AppConfig{
		CombinationSize: DefaultCombinationSize,
		MaxFileSizeMB:   DefaultMaxFileSizeMB,
	}
}

func ValidatePath(path string) error {
	if path == "" {
		return nil
	}

	// Check if directory exists for output files
	dir := filepath.Dir(path)
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		}
	}

	return nil
}

func GetAbsolutePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, path), nil
}
