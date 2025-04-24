package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJsonArrayConfig[T any](filePath string) ([]T, error) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// Unmarshal directly into slice
	result, err := ParseJsonArrayConfig[T](content)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config from %s: %w", filePath, err)
	}
	return result, nil
}

func ReadJsonConfig[T any](filePath string) (*T, error) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// Unmarshal directly into slice
	result, err := ParseJsonConfig[T](content)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config from %s: %w", filePath, err)
	}

	return result, nil
}

func ParseJsonArrayConfig[T any](content []byte) ([]T, error) {
	var result []T
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func ParseJsonConfig[T any](content []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
