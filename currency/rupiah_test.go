package currency

import (
	"testing"
)

func TestFormatRupiah(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1500000, "Rp1.500.000"},
		{1000, "Rp1.000"},
		{0, "Rp0"},
	}

	for _, tt := range tests {
		result := FormatRupiah(tt.input)
		if result != tt.expected {
			t.Errorf("FormatRupiah(%f): expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestParseRupiah(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Rp1.500.000", 1500000},
		{"1.500.000", 1500000},
		{"Rp 1.500.000,50", 1500000.50},
		{"1500", 1500},
	}

	for _, tt := range tests {
		result, err := ParseRupiah(tt.input)
		if err != nil {
			t.Errorf("ParseRupiah(%s) error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("ParseRupiah(%s): expected %f, got %f", tt.input, tt.expected, result)
		}
	}
}
