# Yandex.Disk Bulk Uploader

[English](#english) | [Русский](#русский)

---

## English

A high-performance command-line interface (CLI) utility written in **Go**, designed for automated bulk uploading of files and directories to Yandex.Disk cloud storage. The application mirrors your local filesystem structures in the cloud, automatically publishes uploaded assets, and generates a perfectly aligned ASCII text report.

### ✨ Features

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
    └── run.bat               # Cross-architecture Windows automation wrapper script

### 🚀 Quick Start & Local Setup

#### 1. Clone the repository
    git clone https://github.com/YOUR_USERNAME/yandex-uploader.git
    cd yandex-uploader

#### 2. Configure Local Secrets (Security Isolation)
The repository is designed to prevent your private OAuth credentials from being committed to source control. For Windows users, create a file named `config.local.bat` in the root directory:

    set YANDEX_TOKEN=your_actual_oauth_token_here

#### 3. Compile the Binary
You can cross-compile the program for any specific target system directly from your terminal:

* **For 32-bit Windows Target (Windows 7/8/10 x86):**
    GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o yandex-uploader32.exe ./cmd/uploader

* **For Native 64-bit Windows:**
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader.exe ./cmd/uploader

* **For Linux (64-bit Server/Desktop):**
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader-linux ./cmd/uploader

* **For macOS (Apple Silicon M1/M2/M3):**
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o yandex-uploader-mac ./cmd/uploader

### 💻 Usage & Operational Examples

#### Available CLI Flag Parameters:

| Flag&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; | Data Type | Description |
| :--- | :--- | :--- |
| `-token` | `string` | Yandex.Disk OAuth validation token (falls back to YANDEX_TOKEN environment variable) |
| `-path` | `string` | Local filesystem identifier path targeting a file or a folder for processing |
| `-report` | `string` | Explicit destination output path for compilation metrics (defaults to upload_report.txt) |

#### Windows Environments (CMD / PowerShell)

* **Standard Directory Backup (Implicit Token from config.local.bat):**
```cmd
    run.bat -path="C:\Users\User\Documents\Reports"
    ```

* **Passing the OAuth Token Directly via CLI Flag (Bypassing config files):**
```cmd
    yandex-uploader.exe -token="AgAAAAA..." -path="D:\Backups\database.sql"
    ```

* **Bulk Processing via Positional Arguments (Processing Multiple Targets at Once):**
    The utility natively evaluates trailing positional targets sequentially.
```cmd
    run.bat "C:\ProjectA" "D:\Archive.zip" "E:\Images"
    ```

* **Routing the Generated Link Report to a Dedicated Share or Custom Log Path:**
```cmd
    run.bat -path="C:\Logs" -report="N:\SharedLogs\upload_summary.txt"
    ```

#### Linux & macOS Environments (Terminal)

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

### 📄 Generated Report Sample

Upon processing termination, the utility outputs a clean `upload_report.txt` structure where routing operators align perfectly regardless of varying source string lengths:

    === DOWNLOAD LINKS (2026-06-06 20:15:32) ===
    [Folder] Reports (Entire directory container tree) -> https://disk.yandex.ru/d/exampleRootDirID
    [Folder]   ├── Quarter1                            -> https://disk.yandex.ru/d/exampleSubDirID1
    [File]     ├── financial_statement.xlsx            -> https://disk.yandex.ru/d/exampleFileID1
    [File]     ├── summary_presentation.pdf            -> https://disk.yandex.ru/d/exampleFileID2

> 🔒 License & Disclaimer: Distributed under the MIT License. Always ensure your configuration files, secrets, and compiled binary targets remain untracked within your local `.gitignore` configuration.

---

## Русский

Утилита командной строки (CLI) на языке **Go**, предназначенная для высокопроизводительной массовой загрузки файлов и папок на Яндекс.Диск. Программа автоматически воссоздает локальную структуру директорий в облаке, публикует загруженные объекты и генерирует аккуратно выровненный текстовый отчет с публичными ссылками, совместимый с любыми текстовыми редакторами.

### ✨ Основные возможности

* **Рекурсивный импорт папок:** Полное сохранение иерархии локальных директорий при переносе на удаленный сервер.
* **Потоковый индикатор прогресса:** Отслеживание процентов выполнения и скорости передачи данных в МБ/сек в реальном времени с оптимизированным I/O оверхедом (обновление раз в 200 мс).
* **Автоматическая публикация:** Автономное получение публичных ссылок на каждый файл и папку сразу после завершения их трансфера.
* **Умное логирование:** Создание кастомизированного текстового отчета. Все колонки выравниваются по ширине с использованием UTF-8 рун, а строки разделяются по стандарту Windows (CRLF).
* **Безопасная архитектура:** Готовая конфигурация для изоляции OAuth-токенов в локальной среде (.gitignore), защищающая от случайной утечки секретов в публичный репозиторий.
* **Высокая совместимость:** Оптимизировано под архитектуру 386 для бесперебойной работы на старых 32-битных операционных системах (включая Windows 7).
* **Полная кроссплатформенность:** Написано на чистом Go без платформозависимых зависимостей, благодаря чему утилиту можно собрать под Windows, Linux и macOS.

### 📁 Структура проекта

Архитектура выполнена в соответствии с общепринятым стандартом Go Standard Project Layout:

    yandex-uploader/
    ├── cmd/
    │   └── uploader/
    │       └── main.go       # Точка входа, обработка флагов CLI и запуск конвейера
    ├── internal/
    │   ├── progress/
    │   │   └── progress.go   # Обёртка над io.Reader для калькуляции метрик и вывода прогресса
    │   ├── report/
    │   │   └── report.go     # Движок форматирования и генерации выровненных отчетов
    │   └── yandex/
    │       └── client.go     # Низкоуровневый REST API клиент для работы с Яндекс.Диском
    ├── .gitignore            # Список исключений для Git (скрывает токены, бинарники и логи)
    ├── go.mod                # Манифест Go-модуля
    └── run.bat               # Кросс-архитектурный скрипт автоматизации запуска для Windows

### 🚀 Быстрый старт и локальная настройка

#### 1. Клонирование репозитория
    git clone https://github.com/YOUR_USERNAME/yandex-uploader.git
    cd yandex-uploader

#### 2. Настройка безопасности (Секреты)
Проект настроен так, чтобы ваш OAuth-токен никогда не попал в историю коммитов. Пользователям Windows необходимо создать в корне проекта файл `config.local.bat`:

    set YANDEX_TOKEN=ваш_реальный_oauth_токен_здесь

#### 3. Сборка приложения
Вы можете скомпилировать проект под любую целевую платформу прямо из вашей текущей консоли:

* **Для 32-битной Windows (Windows 7/8/10 x86):**
    GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o yandex-uploader32.exe ./cmd/uploader

* **Для 64-битной Windows:**
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader.exe ./cmd/uploader

* **Для Linux (64-bit Server/Desktop):**
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader-linux ./cmd/uploader

* **Для macOS (Apple Silicon M1/M2/M3):**
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o yandex-uploader-mac ./cmd/uploader

### 💻 Использование и примеры работы

#### Доступные флаговые параметры CLI:

| Флаг&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; | Тип данных | Описание |
| :--- | :--- | :--- |
| `-token` | `string` | OAuth-токен Яндекс.Диска (если не задан, берется из переменной YANDEX_TOKEN) |
| `-path` | `string` | Путь к локальному файлу или папке для загрузки |
| `-report` | `string` | Кастомный путь для генерации файла отчета (по умолчанию: upload_report.txt рядом с exe) |

#### Окружение Windows (Командная строка / PowerShell)

* **Стандартная загрузка папки (токен неявно подтягивается из config.local.bat):**
```cmd
    run.bat -path="C:\Users\User\Documents\Reports"
    ```

* **Прямая передача OAuth-токена через флаг (без использования файлов конфигурации):**
```cmd
    yandex-uploader.exe -token="AgAAAAA..." -path="D:\Backups\database.sql"
    ```

* **Пакетная массовая загрузка нескольких папок и файлов за один запуск:**
    Утилита нативно поддерживает передачу списка путей в качестве последующих аргументов.
```cmd
    run.bat "C:\ProjectA" "D:\Archive.zip" "E:\Images"
    ```

* **Сохранение итогового отчета на сетевой диск или в выделенную папку логов:**
```cmd
    run.bat -path="C:\Logs" -report="N:\SharedLogs\upload_summary.txt"
    ```

#### Окружение Linux и macOS (Терминал)

* **Стандартный запуск с inline-передачей токена в переменной окружения:**
```bash
    chmod +x ./yandex-uploader-linux
    YANDEX_TOKEN="ваш_oauth_токен" ./yandex-uploader-linux -path="/var/www/html/uploads"
    ```

* **Запуск с явным указанием токена и целевого пути через флаги:**
```bash
    chmod +x ./yandex-uploader-mac
    ./yandex-uploader-mac -token="ваш_oauth_токен" -path="/Users/Mac/Desktop/Assets"
    ```

* **Массовая параллельно-последовательная загрузка нескольких путей в Linux:**
```bash
    YANDEX_TOKEN="ваш_oauth_токен" ./yandex-uploader-linux "/etc/nginx/nginx.conf" "/var/log/syslog"
    ```

* **Автоматизация через Cron (настройка ежедневных бэкапов сервера по расписанию):**
    Откройте редактор планировщика командой `crontab -e` и добавьте следующую строку для автоматического запуска утилиты каждую ночь в 02:00:
```text
    0 2 * * * export YANDEX_TOKEN="ваш_oauth_токен" && /usr/local/bin/yandex-uploader-linux -path="/backup/daily" -report="/var/log/uploader_report.txt"
    ```

### 📄 Пример генерируемого отчета

После завершения работы утилита формирует файл `upload_report.txt`, где все стрелочки выровнены строго по вертикали, независимо от длины названий файлов:

    === DOWNLOAD LINKS (2026-06-06 20:15:32) ===
    [Folder] Reports (Entire directory container tree) -> https://disk.yandex.ru/d/exampleRootDirID
    [Folder]   ├── Quarter1                            -> https://disk.yandex.ru/d/exampleSubDirID1
    [File]     ├── financial_statement.xlsx            -> https://disk.yandex.ru/d/exampleFileID1
    [File]     ├── summary_presentation.pdf            -> https://disk.yandex.ru/d/exampleFileID2

> 🔒 Лицензия: Код распространяется под лицензией MIT. Перед отправкой изменений на GitHub убедитесь, что локальные файлы конфигурации с секретами и скомпилированные исполняемые файлы добавлены в ваш `.gitignore`.