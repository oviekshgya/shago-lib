package excel

import (
	"errors"
	"strings"

	"github.com/xuri/excelize/v2"
)

var (
	ErrEmptyColumns    = errors.New("columns must not be empty")
	ErrEmptyColumnName = errors.New("column name must not be empty")
)

// Generate builds an XLSX file dynamically from columns and row data.
// Data is mapped by column name so field ordering always follows columns.
func Generate(columns []string, data []map[string]any) ([]byte, error) {
	if err := validateColumns(columns); err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	defer func() {
		_ = file.Close()
	}()

	sheetName := file.GetSheetName(0)

	for i, column := range columns {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, err
		}
		if err := file.SetCellValue(sheetName, cell, column); err != nil {
			return nil, err
		}
	}

	for rowIndex, row := range data {
		for colIndex, column := range columns {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return nil, err
			}

			if row == nil {
				continue
			}

			if err := file.SetCellValue(sheetName, cell, row[column]); err != nil {
				return nil, err
			}
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func validateColumns(columns []string) error {
	if len(columns) == 0 {
		return ErrEmptyColumns
	}

	for _, column := range columns {
		if strings.TrimSpace(column) == "" {
			return ErrEmptyColumnName
		}
	}

	return nil
}
