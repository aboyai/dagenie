@echo off
setlocal

echo 🔧 Building Dagenie CLI...

REM Save current directory
set BASE_DIR=%CD%

REM Navigate to cmd\cli directory
if exist cmd\cli (
    cd cmd\cli
) else (
    echo ❌ Directory cmd\cli not found!
    exit /b 1
)

REM Build the executable to base directory
go build -o "%BASE_DIR%\dagenie.exe"
IF %ERRORLEVEL% NEQ 0 (
    echo ❌ Build failed!
    exit /b %ERRORLEVEL%
)

REM Return to base directory
cd "%BASE_DIR%"

echo ✅ Build successful! Executable created: dagenie.exe
