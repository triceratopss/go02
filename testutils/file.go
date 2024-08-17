package testutils

import (
	"fmt"
	"os"
	"testing"
)

func ReadJSONFile(t *testing.T, filePath string) (string, error) {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}
