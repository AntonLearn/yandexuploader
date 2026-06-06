// Package progress provides tools to monitor and visualize data transfer operations.
package progress

import (
	"fmt"
	"os"
	"time"
)

// Reader wraps an *os.File to intercept read operations,
// allowing real-time tracking of upload speeds and completion percentages.
type Reader struct {
	File        *os.File
	Total       int64
	ReadBytes   int64
	StartTime   time.Time
	LastPrinted time.Time
}

// NewReader initializes and returns a pointer to a new Reader instance.
func NewReader(file *os.File, total int64) *Reader {
	return &Reader{
		File:  file,
		Total: total,
	}
}

// Read implements the io.Reader interface. It intercepts data chunks read from the file,
// calculates performance metrics, and updates the CLI progress bar throttling output to 200ms.
func (pr *Reader) Read(p []byte) (int, error) {
	if pr.ReadBytes == 0 {
		pr.StartTime = time.Now()
		pr.LastPrinted = time.Now()
	}

	n, err := pr.File.Read(p)
	if n > 0 {
		pr.ReadBytes += int64(n)

		// Throttle terminal output updates to optimize I/O overhead.
		if time.Since(pr.LastPrinted) >= 200*time.Millisecond || pr.ReadBytes >= pr.Total {
			percentage := float64(pr.ReadBytes) / float64(pr.Total) * 100
			duration := time.Since(pr.StartTime).Seconds()

			var speed float64
			if duration > 0 {
				speed = float64(pr.ReadBytes) / 1024 / 1024 / duration
			}

			if pr.ReadBytes >= pr.Total {
				fmt.Printf("\r   Uploading: 100.00%% | [Waiting for Yandex response...]")
			} else {
				fmt.Printf("\r   Uploading: %.2f%% | Speed: %.2f MB/sec", percentage, speed)
			}
			pr.LastPrinted = time.Now()
		}
	}
	return n, err
}
