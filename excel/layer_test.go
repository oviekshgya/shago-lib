package excel

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestGenerateFile(t *testing.T) {
	result, err := GenerateFile(Request{
		FileName: "orders",
		Columns:  []string{"Name", "Age", "City"},
		Data: []map[string]any{
			{"Name": "Budi", "Age": 28, "City": "Jakarta", "Ignored": "X"},
			{"Name": "Sari", "Age": 30, "City": "Bandung"},
		},
	})
	if err != nil {
		t.Fatalf("GenerateFile error: %v", err)
	}

	if result.FileName != "orders.xlsx" {
		t.Fatalf("expected orders.xlsx, got %s", result.FileName)
	}

	if result.ContentType != ContentTypeXLSX {
		t.Fatalf("unexpected content type: %s", result.ContentType)
	}

	if result.OutputPath != "orders.xlsx" {
		t.Fatalf("expected output path orders.xlsx, got %s", result.OutputPath)
	}

	if len(result.Content) == 0 {
		t.Fatal("expected non-empty content")
	}

	workbook, err := excelize.OpenReader(bytes.NewReader(result.Content))
	if err != nil {
		t.Fatalf("OpenReader error: %v", err)
	}
	defer func() {
		_ = workbook.Close()
	}()

	sheet := workbook.GetSheetName(0)

	assertCellValue(t, workbook, sheet, "A1", "Name")
	assertCellValue(t, workbook, sheet, "B1", "Age")
	assertCellValue(t, workbook, sheet, "C1", "City")

	assertCellValue(t, workbook, sheet, "A2", "Budi")
	assertCellValue(t, workbook, sheet, "B2", "28")
	assertCellValue(t, workbook, sheet, "C2", "Jakarta")

	assertCellValue(t, workbook, sheet, "A3", "Sari")
	assertCellValue(t, workbook, sheet, "B3", "30")
	assertCellValue(t, workbook, sheet, "C3", "Bandung")
}

func TestGenerateValidateColumns(t *testing.T) {
	_, err := Generate(nil, nil)
	if !errors.Is(err, ErrEmptyColumns) {
		t.Fatalf("expected ErrEmptyColumns, got %v", err)
	}

	_, err = Generate([]string{"Name", " "}, nil)
	if !errors.Is(err, ErrEmptyColumnName) {
		t.Fatalf("expected ErrEmptyColumnName, got %v", err)
	}
}

func TestGenerateAndSave(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "employees-report")

	result, err := GenerateAndSave(Request{
		Columns: []string{"Name"},
		Data: []map[string]any{
			{"Name": "Budi"},
		},
	}, outputPath)
	if err != nil {
		t.Fatalf("GenerateAndSave error: %v", err)
	}

	if filepath.Ext(result.OutputPath) != ".xlsx" {
		t.Fatalf("expected .xlsx extension, got %s", result.OutputPath)
	}

	info, err := os.Stat(result.OutputPath)
	if err != nil {
		t.Fatalf("expected file exists: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected generated file not empty")
	}
}

func assertCellValue(t *testing.T, workbook *excelize.File, sheet, cell, expected string) {
	t.Helper()

	value, err := workbook.GetCellValue(sheet, cell)
	if err != nil {
		t.Fatalf("GetCellValue(%s) error: %v", cell, err)
	}

	if value != expected {
		t.Fatalf("cell %s expected %s, got %s", cell, expected, value)
	}
}
