package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WriteToJson writes structured data to a JSON file
func WriteToJson(filename string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}

// ReadFromJson reads structured data from a JSON file
func ReadFromJson(filename string, dest interface{}) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, dest)
}

func ParseDate(input string) (time.Time, error) {
    layout := "02/01/2006" // DD/MM/YYYY
    parsedDate, err := time.Parse(layout, input)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid date format. Please use DD/MM/YYYY")
    }
    return parsedDate, nil
}
