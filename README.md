## Yandex.Disk Bulk Uploader

---

## Features

* **Recursive Folder Ingestion:** Replicates complex local directory trees onto the remote storage seamlessly.
* **Streaming Progress Indicator:** Real-time throughput visualization tracking percentages and transfer speeds (MB/s) with optimized I/O overhead (200ms throttle updates).
* **Automated Asset Publishing:** Autonomously executes remote publishing pipelines to fetch public delivery URLs right after each transmission.
* **Smart Layout Reporting:** Generates clean, padded text reports. Column positions are dynamically calculated using UTF-8 runes to match padding lengths, with native Windows line endings (CRLF).
* **Secure Architecture:** Pre-configured environment decoupling (.gitignore) keeps your sensitive OAuth API credentials safe from accidental upstream leakage.
* **Legacy System Back-Compatibility:** Fully optimized for the 386 (32-bit) architecture to run reliably on legacy server environments and older OS builds (e.g., Windows 7).
* **Cross-Platform Native Code:** Written in pure Go with no OS-specific bindings, enabling native execution across Windows, Linux, and macOS.

### 📁 Project Structure

The codebase strictly follows the idiomatic Go Standard Project Layout:

    yandex-uploader/
    ├── cmd/
    │   └── uploader/
    │       └── main.go       # Application entry point (CLI parsing & pipeline orchestration)
    ├── internal/
    │   ├── progress/
    │   │   └── progress.go   # io.Reader proxy tracking data transfer metrics and telemetry
    │   ├── report/
    │   │   └── report.go     # Dynamic padding layout engine generating aligned reports
    │   └── yandex/
    │       └── client.go     # Low-level high-throughput REST API client for Yandex.Disk
    ├── .gitignore            # Git exclusion rules (isolates tokens, binaries, and local logs)
    ├── go.mod                # Go module manifest dependencies
    └── upload.bat               # Cross-architecture Windows automation wrapper script

---

## 🚀 Quick Start & Local Setup

### 1. Clone the repository
    git clone https://github.com/antonlearn/yandex-uploader.git
    cd yandex-uploader

### 2. Configure Local Secrets (Security Isolation)
The repository is designed to prevent your private OAuth credentials from being committed to source control. For Windows users, create a file named `config.local.bat` in the root directory:

    set YANDEX_TOKEN=your_actual_oauth_token_here

### 3. Compile the Binary
You can cross-compile the program for any specific target system directly from your terminal:

* **For 32-bit Windows Target (Windows 7/8/10 x86):**
    GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o yandex-uploader32.exe ./cmd/uploader

* **For Native 64-bit Windows:**
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader.exe ./cmd/uploader

* **For Linux (64-bit Server/Desktop):**
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader-linux ./cmd/uploader

* **For macOS (Apple Silicon M1/M2/M3):**
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o yandex-uploader-mac ./cmd/uploader

---

## 💻 Usage & Operational Examples

### Available CLI Flag Parameters

| Flag | Data Type | Description |
| :--- | :--- | :--- |
| <code>&#8209;token</code> | `string` | Yandex.Disk OAuth validation token. Falls back to the `YANDEX_TOKEN` environment variable if not specified. **This flag takes absolute precedence over environment variables.** |
| <code>&#8209;path</code> | `string` | Local filesystem identifier path targeting a file or a folder for processing. |
| <code>&#8209;report</code> | `string` | Explicit destination output path for compilation metrics. If omitted, defaults to an `upload_report.txt` file generated next to the application binary execution context. | phot |

### Windows Environments (CMD / PowerShell)

* **Standard Directory Backup (Implicit Token from config.local.bat):**
```cmd
    upload.bat -path="C:\Users\User\Documents\Reports"
```

* **Passing the OAuth Token Directly via CLI Flag (Bypassing config files / overriding env):**
```cmd
    yandex-uploader.exe -token="AgAAAAA..." -path="D:\Backups\database.sql"
```

* **Bulk Processing via Positional Arguments (Processing Multiple Targets at Once):**
    The utility natively evaluates trailing positional targets sequentially.
```cmd
    upload.bat "C:\ProjectA" "D:\Archive.zip" "E:\Images"
```

* **Routing the Generated Link Report to a Dedicated Share or Custom Log Path:**
```cmd
    upload.bat -path="C:\Logs" -report="N:\SharedLogs\upload_summary.txt"
```

### Linux & macOS Environments (Terminal)

* **Standard Execution with Inline Environment Variable Injection:**
```bash
    chmod +x ./yandex-uploader-linux
    YANDEX_TOKEN="your_oauth_token" ./yandex-uploader-linux -path="/var/www/html/uploads"
```

* **Explicit Token and Target Flag Definition:**
```bash
    chmod +x ./yandex-uploader-mac
    ./yandex-uploader-mac -token="your_oauth_token" -path="/Users/Mac/Desktop/Assets"
```

* **Batch Multi-Path Transmission under Linux:**
```bash
    YANDEX_TOKEN="your_oauth_token" ./yandex-uploader-linux "/etc/nginx/nginx.conf" "/var/log/syslog"
```

* **Automating via Headless Cron Job (Nightly Unattended Server Backups):**
    Open your crontab editor via `crontab -e` and append the following configuration to automate uploads every night at 2:00 AM:
```text
    0 2 * * * export YANDEX_TOKEN="your_oauth_token" && /usr/local/bin/yandex-uploader-linux -path="/backup/daily" -report="/var/log/uploader_report.txt"
```

---

## 📄 Generated Report Sample

Upon processing termination, the utility outputs a clean `upload_report.txt` structure where routing operators align perfectly regardless of varying source string lengths:

    === DOWNLOAD LINKS (2026-06-06 20:15:32) ===
    [Folder] Reports (Entire directory container tree) -> https://disk.yandex.ru/d/exampleRootDirID
    [Folder]   ├── Quarter1                             -> https://disk.yandex.ru/d/exampleSubDirID1
    [File]     ├── financial_statement.xlsx             -> https://disk.yandex.ru/d/exampleFileID1
    [File]     ├── summary_presentation.pdf             -> https://disk.yandex.ru/d/exampleFileID2

> 🔒 License & Disclaimer: Distributed under the MIT License. Always ensure your configuration files, secrets, and compiled binary targets remain untracked within your local `.gitignore` configuration.
```