// Main entry point for the Yandex.Disk bulk storage automated upload engine.
// Located in the standard cmd/ idiomatic directory structure.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antonlearn/yandex-uploader/internal/report"
	"github.com/antonlearn/yandex-uploader/internal/yandex"
)

func main() {
	tokenFlag := flag.String("token", "", "Yandex.Disk OAuth API validation token credentials")
	pathFlag := flag.String("path", "", "Target storage filesystem location identifier pointing to a file/folder")
	reportFlag := flag.String("report", "", "Absolute path where compilation execution details logs will be saved")
	flag.Parse()

	token := *tokenFlag
	if token == "" {
		token = os.Getenv("YANDEX_TOKEN")
	}

	// Enhanced token validation routing execution errors directly to the standard error stream (os.Stderr).
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error execution failure: OAuth token credentials identifier parameter is missing.")
		fmt.Fprintln(os.Stderr, "Please provide a valid token via the -token flag or set the YANDEX_TOKEN environment variable.")
		os.Exit(1)
	}

	var inputPaths []string
	if flag.NArg() > 0 {
		inputPaths = flag.Args()
	} else if *pathFlag != "" {
		inputPaths = []string{*pathFlag}
	}

	// Maintain logging consistency by redirecting missing target path exceptions to os.Stderr.
	if len(inputPaths) == 0 {
		fmt.Fprintln(os.Stderr, "Error execution failure: No valid local target paths were specified for ingestion pipelines processing.")
		os.Exit(1)
	}

	// Initialize modularized infrastructure subsystems client components.
	diskClient := yandex.NewClient(token)
	var logRegistry []report.Item

	fmt.Printf("[Initialization] Processing queue contains total of %d object tracking items\n", len(inputPaths))

	for _, rawPath := range inputPaths {
		rawPath = strings.TrimSpace(rawPath)
		if rawPath == "" {
			continue
		}

		targetPath, _ := filepath.Abs(rawPath)
		fileInfo, err := os.Stat(targetPath)
		if os.IsNotExist(err) {
			fmt.Printf("\n[Bypassing Resource] Specified target platform path does not exist: %s\n", targetPath)
			continue
		}

		fmt.Printf("\n=== Operational Pipeline Active: %s ===", filepath.Base(targetPath))

		if fileInfo.IsDir() {
			err := diskClient.UploadFolder(targetPath, &logRegistry)
			if err != nil {
				fmt.Printf("\n[Error Event] Directory ingestion operation failure sequence logged: %v\n", err)
			}
		} else {
			remotePath := "/" + filepath.Base(targetPath)
			err := diskClient.UploadSingleFile(targetPath, remotePath)
			leftText := fmt.Sprintf("[File]  %s", filepath.Base(targetPath))

			if err == nil {
				// Trigger the updated public asset publishing routine to obtain a delivery URL.
				pubLink, errPub := diskClient.PublishAndGetLink(remotePath)
				if errPub != nil {
					logRegistry = append(logRegistry, report.Item{Left: leftText, Right: fmt.Sprintf("(Publish failed: %v)", errPub)})
				} else {
					logRegistry = append(logRegistry, report.Item{Left: leftText, Right: pubLink})
				}
			} else {
				logRegistry = append(logRegistry, report.Item{Left: leftText, Right: fmt.Sprintf("(Transfer failed: %v)", err)})
			}
		}
	}

	// Delegate report pipeline execution sequences directly to the reporting package layout system.
	if len(logRegistry) > 0 {
		var reportFile string
		if *reportFlag != "" {
			reportFile = *reportFlag
		} else {
			// os.Executable() successfully evaluates binary location regardless of source layouts.
			if exePath, err := os.Executable(); err == nil {
				reportFile = filepath.Join(filepath.Dir(exePath), "upload_report.txt")
			} else {
				reportFile = "upload_report.txt"
			}
		}

		if err := report.Generate(reportFile, logRegistry); err != nil {
			fmt.Printf("\n[Warning Event] Report writer pipeline encountered failure constraints: %v\n", err)
		} else {
			fmt.Printf("\n\n[Success Process Completed] Objects completely dispatched.\n[Log Output Targets] Filepath: %s\n", reportFile)
		}
	}
}
