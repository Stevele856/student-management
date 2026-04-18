package utils

import (
	"encoding/csv"
	"fmt"
	"io"
)

// ReadCSV reads all rows from a CSV reader, skipping the header row.
// Returns the header separately and all remaining rows.
func ReadCSV(r io.Reader) (header []string, rows [][]string, err error){
	reader := csv.NewReader(r)

	// Read header row first
	header, err = reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// reader remaining rows one by one
	for {
		row, err := reader.Read()
		if err == io.EOF{
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read CSV row: %w", err)
		}
		rows = append(rows, row)
	}
	return header, rows, nil
}

// WriteCSV writes a header and rows to a CSV writer
func WriteCSV(w io.Writer, header []string, rows [][]string) error{
	writer := csv.NewWriter(w)

	// Write header row
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write each data row
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	// Flush buffered data to the underlying writer
	writer.Flush()
	if err := writer.Error(); err != nil{
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return nil
}