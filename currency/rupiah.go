package currency

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Option defines formatting options
type Option struct {
	Prefix  string
	Decimal int
}

var defaultOption = Option{
	Prefix:  "Rp",
	Decimal: 0,
}

// FormatRupiah converts float64 to Rupiah string format
// Example: 1500000 -> "Rp1.500.000"
func FormatRupiah(amount float64) string {
	return FormatRupiahWithOption(amount, defaultOption)
}

// FormatRupiahWithOption converts with custom options
func FormatRupiahWithOption(amount float64, opt Option) string {
	p := message.NewPrinter(language.Indonesian)

	// Create format string based on decimal
	formatStr := "%." + fmt.Sprintf("%d", opt.Decimal) + "f"

	formattedNumber := p.Sprintf(formatStr, amount)

	// Default 'text/message' uses comma for decimal in Indonesian, which is correct (1.500,00)
	// If the user wants standard "Rp" prefix without space:
	return opt.Prefix + formattedNumber
}

// ParseRupiah converts Rupiah string to float64
// Handles "Rp1.500.000", "1.500.000", "Rp 1.500.000,00"
func ParseRupiah(s string) (float64, error) {
	// Remove Prefix "Rp" (case insensitive) and whitespaces
	re := regexp.MustCompile(`(?i)rp\s*`)
	cleanStr := re.ReplaceAllString(s, "")
	cleanStr = strings.TrimSpace(cleanStr)

	// Remove thousands separator (.)
	cleanStr = strings.ReplaceAll(cleanStr, ".", "")

	// Replace decimal separator (,) with (.) for float parsing
	cleanStr = strings.ReplaceAll(cleanStr, ",", ".")

	return strconv.ParseFloat(cleanStr, 64)
}
