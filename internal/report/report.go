// Package report handles formatting, structural alignment, and writing
// execution summaries that remain compatible across diverse OS environments.
package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

// Item stores structural components of a log row to isolate formatting from data collection.
type Item struct {
	Left  string // Object type prefix along with its relative layout alignment structure
	Right string // Generated public download link or captured operation error message
}

// Generate writes a visually aligned report file to the designated path.
// It explicitly utilizes Windows-style CRLF (\r\n) termination for backwards compatibility with legacy text editors.
func Generate(targetPath string, lines []Item) error {
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header containing explicit CRLF line break for target environment compatibility.
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("=== DOWNLOAD LINKS (%s) ===\r\n", timestamp)); err != nil {
		return err
	}

	// Calculate maximum layout width using UTF-8 rune sequence length instead of standard byte counters.
	maxRunes := 0
	for _, item := range lines {
		currentRunes := utf8.RuneCountInString(item.Left)
		if currentRunes > maxRunes {
			maxRunes = currentRunes
		}
	}

	// Write padded items sequentially ensuring strict vertical column alignments.
	for _, item := range lines {
		currentRunes := utf8.RuneCountInString(item.Left)
		padding := strings.Repeat(" ", maxRunes-currentRunes)

		row := fmt.Sprintf("%s%s -> %s\r\n", item.Left, padding, item.Right)
		if _, err := f.WriteString(row); err != nil {
			return err
		}
	}

	return nil
}
