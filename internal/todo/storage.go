package todo

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	
)

// SaveToCSV writes the current tasks to a CSV file	
func (s *Service) SaveToCSV(filename string) error {
	// 1. Create or overwrite the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 2. Initialize the CSV writer
	writer := csv.NewWriter(file)

	defer writer.Flush()

	headers := []string{"ID", "Description", "CreatedAt", "IsComplete"}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, task := range s.tasks {
		row := []string{
			strconv.Itoa(task.ID),
			task.Description,
			task.CreatedAt.Format(time.RFC3339),
			strconv.FormatBool(task.IsComplete),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row for task ID %d: %w", task.ID, err)
		}
	}
	return nil
}

// LoadFromCSV reads a CSV file and populates the service's tasks slice
func (s *Service) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to find with filename %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// 3. Read the header row first to skip it
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to load csv headers %w", err)
	}
	// 4. Read all the remaining rows
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read the data rows %w", err)
	}

	var loadedTasks []Task

	for i, row := range records {
		// Validation check: ensure row has the expected number of columns
		if len(row) < 4 {
			return fmt.Errorf("malformed row at line %d: expected 4 got %d ", i+2, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("error converting id type from string to integer %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, row[2])
		if err != nil {
			return fmt.Errorf("line %d: invalid time stamp '%s':%w", i+2, row[2], err)
		}

		isComplete, err := strconv.ParseBool(row[3])
		if err != nil {
			return fmt.Errorf("line %d: invalid boolean '%s':%w", i+2, row[3], err)
		}

		task := Task{
			ID:          id,
			Description: row[1],
			CreatedAt:   createdAt,
			IsComplete:  isComplete,
		}
		loadedTasks = append(loadedTasks, task)
	}
	// 6. Update the service state only after successful parsing
	s.tasks = loadedTasks
	return nil
}