@echo off
:: Disable command echoing to keep the console output clean and readable

:: -----------------------------------------------------------------------------
:: CONFIGURATION (LOAD LOCAL SECRETS IF EXIST)
:: -----------------------------------------------------------------------------
if exist "%~dp0config.local.bat" (
    call "%~dp0config.local.bat"
) else (
    echo [Warning] config.local.bat not found. Trying system environment...
)

:: Check if token was successfully set either from file or system environment
if "%YANDEX_TOKEN%"=="" (
    echo [Error] YANDEX_TOKEN is not set. Please create config.local.bat
    pause
    exit /b 1
)

:: -----------------------------------------------------------------------------
:: ARCHITECTURE DETECTION
:: -----------------------------------------------------------------------------
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" goto x64
if "%PROCESSOR_ARCHITEW6432%"=="AMD64" goto x64

:x86
:: Fallback configuration: Target the 32-bit architecture executable binary
set BINARY_NAME=yandex-uploader32.exe
goto run

:x64
:: Target configuration: Target the native 64-bit architecture executable binary
set BINARY_NAME=yandex-uploader.exe
goto run

:run
:: -----------------------------------------------------------------------------
:: EXECUTION PIPELINE
:: -----------------------------------------------------------------------------
:: Change the console active code page to UTF-8 (65001) to correctly render
:: Cyrillic file structures, directory names, and folder path payloads in logs
chcp 65001 > nul

:: Execute the determined binary located directly in the script's directory (%~dp0).
:: Forwards the authentication token credentials and appends all command-line 
:: arguments (%*) passed externally to this batch wrapper script.
"%~dp0%BINARY_NAME%" -token="%YANDEX_TOKEN%" %*

:: Restore the standard Russian OEM code page (866) to reset terminal state environments
chcp 866 > nul

echo.
:: Keep the command prompt window open after execution completes,
:: allowing the user to review the processing metrics and link tables before exiting.
pause