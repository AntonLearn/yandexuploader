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
The repository is designed to prevent your private OAuth credentials from being committed to source control. Create a file named config.local.bat in the root directory:

    set YANDEX_TOKEN=your_actual_oauth_token_here

#### 3. Compile the Binary
You can cross-compile the program for any specific target system directly from your terminal:

* For Current Development OS:
    go build -o uploader ./cmd/uploader

* For 32-bit Windows Target (Windows 7/8/10 x86):
    GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o yandex-uploader32.exe ./cmd/uploader

* For Native 64-bit Windows:
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader.exe ./cmd/uploader

### 💻 Usage

Execute the pipeline via terminal or leverage the automated run.bat wrapper, which implicitly evaluates your processor architecture and seamlessly injects the token from config.local.bat.

#### Operational Examples:
    # Upload a single directory container
    run.bat -path="C:\Users\User\Documents\Reports"

    # Transmit an isolated archive and route logs to a custom destination
    run.bat -path="C:\Photos\vacation.zip" -report="D:\logs\my_links.txt"

    # Process multiple filesystem path arguments sequentially
    run.bat "C:\Data" "D:\Backup.rar"

#### Available CLI Flag Parameters:
| Flag | Data Type | Description |
| :--- | :--- | :--- |
| -token | string | Yandex.Disk OAuth validation token (falls back to YANDEX_TOKEN environment variable) |
| -path | string | Local filesystem identifier path targeting a file or a folder for processing |
| -report | string | Explicit destination output path for compilation metrics (defaults to upload_report.txt) |

### 📄 Generated Report Sample

Upon processing termination, the utility outputs a clean upload_report.txt structure where routing operators align perfectly regardless of varying source string lengths:

    === DOWNLOAD LINKS (2026-06-06 20:15:32) ===
    [Folder] Reports (Entire directory container tree) -> https://disk.yandex.ru/d/exampleRootDirID
    [Folder]   ├── Quarter1                            -> https://disk.yandex.ru/d/exampleSubDirID1
    [File]     ├── financial_statement.xlsx            -> https://disk.yandex.ru/d/exampleFileID1
    [File]     ├── summary_presentation.pdf            -> https://disk.yandex.ru/d/exampleFileID2

> 🔒 License & Disclaimer: Distributed under the MIT License. Always ensure your config.local.bat and compiled *.exe binaries remain untracked within your local .gitignore setup prior to shifting upstream pushes.

---

## Русский

Утилита командной строки (CLI) на языке **Go**, предназначенная для высокопроизводительной массовой загрузки файлов и папок на Яндекс.Диск. Программа автоматически воссоздает локальную структуру директорий в облаке, публикует загруженные объекты и генерирует аккуратно выровненный текстовый отчет с публичными ссылками, совместимый с любыми текстовыми редакторами на Windows.

### ✨ Основные возможности

* **Рекурсивный импорт папок:** Полное сохранение иерархии локальных директорий при переносе на удаленный сервер.
* **Потоковый индикатор прогресса:** Отслеживание процентов выполнения и скорости передачи данных в МБ/сек в реальном времени с оптимизированным I/O оверхедом (обновление раз в 200 мс).
* **Автоматическая публикация:** Автономное получение публичных ссылок на каждый файл и папку сразу после завершения их трансфера.
* **Умное логирование:** Создание кастомизированного текстового отчета. Все колонки выравниваются по ширине с использованием UTF-8 рун, а строки разделяются по стандарту Windows (CRLF).
* **Безопасная архитектура:** Готовая конфигурация для изоляции OAuth-токенов в локальной среде (.gitignore), защищающая от случайной утечки секретов в публичный репозиторий.
* **Высокая совместимость:** Оптимизировано под архитектуру 386 для бесперебойной работы на старых 32-битных операционных системах (включая Windows 7).

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
Проект настроен так, чтобы ваш OAuth-токен никогда не попал в историю коммитов. Создайте в корне проекта файл config.local.bat:

    set YANDEX_TOKEN=ваш_реальный_oauth_токен_здесь

#### 3. Сборка приложения
Вы можете скомпилировать проект под любую целевую платформу:

* Для текущей ОС (разработка):
    go build -o uploader ./cmd/uploader

* Для 32-битной Windows (Windows 7/8/10 x86):
    GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o yandex-uploader32.exe ./cmd/uploader

* Для 64-битной Windows:
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o yandex-uploader.exe ./cmd/uploader

### 💻 Использование

Запуск осуществляется через командную строку или с помощью универсального скрипта run.bat, который сам определит архитектуру вашего процессора и подтянет токен из config.local.bat.

#### Примеры команд:
    # Загрузка одной конкретной папки
    run.bat -path="C:\Users\User\Documents\Reports"

    # Загрузка одного файла с указанием кастомного пути для сохранения отчета
    run.bat -path="C:\Photos\vacation.zip" -report="D:\logs\my_links.txt"

    # Массовая загрузка нескольких аргументов подряд
    run.bat "C:\Data" "D:\Backup.rar"

#### Доступные флаги CLI:
| Флаг | Тип данных | Описание |
| :--- | :--- | :--- |
| -token | string | OAuth-токен Яндекс.Диска (если не задан, берется из переменной YANDEX_TOKEN) |
| -path | string | Путь к локальному файлу или папке для загрузки |
| -report | string | Кастомный путь для генерации файла отчета (по умолчанию: upload_report.txt рядом с exe) |

### 📄 Пример генерируемого отчета

После завершения работы утилита формирует файл upload_report.txt, где все стрелочки выровнены строго по вертикали, независимо от длины названий файлов:

    === DOWNLOAD LINKS (2026-06-06 20:15:32) ===
    [Folder] Reports (Entire directory container tree) -> https://disk.yandex.ru/d/exampleRootDirID
    [Folder]   ├── Quarter1                            -> https://disk.yandex.ru/d/exampleSubDirID1
    [File]     ├── financial_statement.xlsx            -> https://disk.yandex.ru/d/exampleFileID1
    [File]     ├── summary_presentation.pdf            -> https://disk.yandex.ru/d/exampleFileID2

> 🔒 Лицензия: Код распространяется под лицензией MIT. Перед отправкой изменений на GitHub убедитесь, что файлы config.local.bat, *.exe и upload_report.txt находятся в списке исключений вашего .gitignore.