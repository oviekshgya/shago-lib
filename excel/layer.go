package excel

import (
	"os"
	"path/filepath"
	"strings"
)

const ContentTypeXLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// Request is the dynamic input for GenerateFile.
type Request struct {
	Columns  []string
	Data     []map[string]any
	FileName string
}

// Result contains generated file metadata and bytes.
type Result struct {
	FileName    string
	ContentType string
	Content     []byte
	OutputPath  string
}

// GenerateFile is a higher-level layer that returns a ready-to-send file payload.
// Caller only needs to provide column names and row data.
func GenerateFile(req Request) (Result, error) {
	content, err := Generate(req.Columns, req.Data)
	if err != nil {
		return Result{}, err
	}

	fileName := normalizeFileName(req.FileName)

	return Result{
		FileName:    fileName,
		ContentType: ContentTypeXLSX,
		Content:     content,
		OutputPath:  fileName,
	}, nil
}

// GenerateAndSave generates an Excel file and writes it to disk.
// If outputPath is empty, the default output file name from request is used.
func GenerateAndSave(req Request, outputPath string) (Result, error) {
	result, err := GenerateFile(req)
	if err != nil {
		return Result{}, err
	}

	path := strings.TrimSpace(outputPath)
	if path == "" {
		path = result.FileName
	}

	if strings.TrimSpace(filepath.Ext(path)) == "" {
		path += ".xlsx"
	}

	if err := os.WriteFile(path, result.Content, 0644); err != nil {
		return Result{}, err
	}

	result.OutputPath = path
	return result, nil
}

func normalizeFileName(fileName string) string {
	name := strings.TrimSpace(fileName)
	if name == "" {
		return "report.xlsx"
	}

	if strings.HasSuffix(strings.ToLower(name), ".xlsx") {
		return name
	}

	return name + ".xlsx"
}
