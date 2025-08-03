@echo off
echo Starting E-commerce Backend Server...
echo.
echo Backend will start on http://localhost:8080
echo.

REM Navigate to backend directory
cd /d "%~dp0backend"

REM Check if Go is accessible
go version
if %errorlevel% neq 0 (
    echo.
    echo ERROR: Go is not found. Please ensure Go is installed and in PATH.
    pause
    exit /b 1
)

echo.
echo Installing dependencies...
go mod tidy

echo.
echo Starting backend server...
go run main.go
pause
