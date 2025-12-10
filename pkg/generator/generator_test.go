package generator

import (
	"testing"
)

func TestCalculateTotalCombinations(t *testing.T) {
	tests := []struct {
		name            string
		passwords       []string
		combinationSize int
		extraSymbols    []rune
		symbolPositions []SymbolPosition
		expected        int64
	}{
		{
			name:            "basic 2 from 3 passwords",
			passwords:       []string{"a", "b", "c"},
			combinationSize: 2,
			extraSymbols:    []rune{},
			symbolPositions: []SymbolPosition{},
			expected:        9, // 3^2
		},
		{
			name:            "basic 3 from 2 passwords",
			passwords:       []string{"a", "b"},
			combinationSize: 3,
			extraSymbols:    []rune{},
			symbolPositions: []SymbolPosition{},
			expected:        8, // 2^3
		},
		{
			name:            "with symbols",
			passwords:       []string{"a", "b"},
			combinationSize: 2,
			extraSymbols:    []rune{'!', '@'},
			symbolPositions: []SymbolPosition{PositionStart, PositionEnd},
			expected:        20, // 2^2 * (1 + 2*2) = 4 * 5 = 20
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				config: Config{
					CombinationSize: tt.combinationSize,
					ExtraSymbols:    tt.extraSymbols,
					SymbolPositions: tt.symbolPositions,
				},
				passwords: tt.passwords,
			}

			result := g.CalculateTotalCombinations()
			if result != tt.expected {
				t.Errorf("CalculateTotalCombinations() = %v, want %v", result, tt.expected)
			}
		})
	}
}
