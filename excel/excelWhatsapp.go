package excel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type WorkbookJSON struct {
	File   string               `json:"file"`
	Sheets map[string]SheetJSON `json:"sheets"`
}

type SheetJSON struct {
	Type  string           `json:"type"`
	Title string           `json:"title,omitempty"`
	Data  map[string]any   `json:"data,omitempty"`
	Items []map[string]any `json:"items,omitempty"`
	Text  []string         `json:"text,omitempty"`
	Rows  [][]any          `json:"rows,omitempty"`
}

func ExcelToJSON(filePath string) (*WorkbookJSON, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := &WorkbookJSON{
		File:   filePath,
		Sheets: make(map[string]SheetJSON),
	}

	for _, sheetName := range f.GetSheetList() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("read sheet %s: %w", sheetName, err)
		}

		cleaned := cleanRows(rows)
		result.Sheets[sheetName] = parseSheet(cleaned)
	}

	return result, nil
}

func parseSheet(rows [][]string) SheetJSON {
	if len(rows) == 0 {
		return SheetJSON{
			Type: "empty",
		}
	}

	title := ""
	if nonEmptyCount(rows[0]) == 1 {
		title = firstNonEmpty(rows[0])
	}

	// Pattern: title + key-value rows
	if start, ok := detectKeyValue(rows); ok {
		data := map[string]any{}

		for i := start; i < len(rows); i++ {
			row := rows[i]
			if len(row) < 2 {
				continue
			}

			key := strings.TrimSpace(row[0])
			value := strings.TrimSpace(row[1])

			if key == "" || value == "" {
				continue
			}

			jsonKey := toSnakeCase(key)
			data[jsonKey] = normalizeValue(key, value)
		}

		return SheetJSON{
			Type:  "key_value",
			Title: title,
			Data:  data,
		}
	}

	// Pattern: title + table header + data rows
	if headerIdx := findTableHeader(rows); headerIdx >= 0 {
		headers := normalizeHeaders(rows[headerIdx])
		items := make([]map[string]any, 0)

		for i := headerIdx + 1; i < len(rows); i++ {
			row := rows[i]
			if nonEmptyCount(row) == 0 {
				continue
			}

			item := map[string]any{}

			for colIdx, header := range headers {
				if header == "" {
					continue
				}

				value := ""
				if colIdx < len(row) {
					value = strings.TrimSpace(row[colIdx])
				}

				if value == "" {
					item[header] = nil
				} else {
					item[header] = normalizeValue(header, value)
				}
			}

			if len(item) > 0 {
				items = append(items, item)
			}
		}

		return SheetJSON{
			Type:  "table",
			Title: title,
			Items: items,
		}
	}

	// Pattern: title + text rows
	text := make([]string, 0)
	start := 0
	if title != "" {
		start = 1
	}

	for i := start; i < len(rows); i++ {
		for _, cell := range rows[i] {
			cell = strings.TrimSpace(cell)
			if cell != "" {
				text = append(text, cell)
			}
		}
	}

	if len(text) > 0 {
		return SheetJSON{
			Type:  "text",
			Title: title,
			Text:  text,
		}
	}

	// Fallback: raw rows
	rawRows := make([][]any, 0)
	for _, row := range rows {
		rawRow := make([]any, 0)
		for _, cell := range row {
			rawRow = append(rawRow, strings.TrimSpace(cell))
		}
		rawRows = append(rawRows, rawRow)
	}

	return SheetJSON{
		Type:  "raw",
		Title: title,
		Rows:  rawRows,
	}
}

func cleanRows(rows [][]string) [][]string {
	cleaned := make([][]string, 0)

	for _, row := range rows {
		newRow := make([]string, len(row))
		for i, cell := range row {
			newRow[i] = strings.TrimSpace(cell)
		}

		newRow = trimTrailingEmptyCells(newRow)

		if nonEmptyCount(newRow) > 0 {
			cleaned = append(cleaned, newRow)
		}
	}

	return cleaned
}

func trimTrailingEmptyCells(row []string) []string {
	last := -1

	for i := len(row) - 1; i >= 0; i-- {
		if strings.TrimSpace(row[i]) != "" {
			last = i
			break
		}
	}

	if last == -1 {
		return []string{}
	}

	return row[:last+1]
}

func detectKeyValue(rows [][]string) (int, bool) {
	start := 0

	if len(rows) > 1 && nonEmptyCount(rows[0]) == 1 {
		start = 1
	}

	total := 0
	valid := 0

	for i := start; i < len(rows); i++ {
		row := rows[i]
		count := nonEmptyCount(row)

		if count == 0 {
			continue
		}

		total++

		if len(row) >= 2 &&
			strings.TrimSpace(row[0]) != "" &&
			strings.TrimSpace(row[1]) != "" &&
			count <= 2 {
			valid++
		}
	}

	return start, total >= 2 && total == valid
}

func findTableHeader(rows [][]string) int {
	for i := 0; i < len(rows)-1; i++ {
		currentCount := nonEmptyCount(rows[i])
		nextCount := nonEmptyCount(rows[i+1])

		if currentCount >= 3 && nextCount >= 2 {
			return i
		}
	}

	return -1
}

func normalizeHeaders(row []string) []string {
	headers := make([]string, len(row))

	for i, header := range row {
		headers[i] = toSnakeCase(header)
	}

	return headers
}

func normalizeValue(key string, value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	lowerKey := strings.ToLower(key)

	// Field seperti kode, nomor, contact jangan dipaksa jadi number.
	// Contoh: nomor HP, kode merchant, ID, nomor urut.
	keepStringKeys := []string{
		"kode",
		"code",
		"contact",
		"kontak",
		"phone",
		"telp",
		"wa",
		"whatsapp",
		"id",
		"nomor",
		"no",
	}

	for _, marker := range keepStringKeys {
		if strings.Contains(lowerKey, marker) {
			return value
		}
	}

	if strings.EqualFold(value, "true") {
		return true
	}

	if strings.EqualFold(value, "false") {
		return false
	}

	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	return value
}

func nonEmptyCount(row []string) int {
	count := 0

	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			count++
		}
	}

	return count
}

func firstNonEmpty(row []string) string {
	for _, cell := range row {
		cell = strings.TrimSpace(cell)
		if cell != "" {
			return cell
		}
	}

	return ""
}

func toSnakeCase(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	replacer := strings.NewReplacer(
		"/", " ",
		"-", " ",
		".", " ",
		":", " ",
		"(", " ",
		")", " ",
	)

	input = replacer.Replace(input)

	reSpace := regexp.MustCompile(`\s+`)
	input = reSpace.ReplaceAllString(input, "_")

	reInvalid := regexp.MustCompile(`[^a-z0-9_]+`)
	input = reInvalid.ReplaceAllString(input, "")

	reUnderscore := regexp.MustCompile(`_+`)
	input = reUnderscore.ReplaceAllString(input, "_")

	return strings.Trim(input, "_")
}
