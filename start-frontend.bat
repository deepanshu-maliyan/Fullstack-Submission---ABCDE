@echo off
echo Starting E-commerce Frontend Server...
echo.
echo Adding Node.js to PATH and starting server...
echo Frontend will start on http://localhost:3000
echo.

REM Add Node.js to PATH
set PATH=%PATH%;C:\Program Files\nodejs

REM Navigate to frontend directory
cd /d "%~dp0frontend"

REM Check if Node.js is accessible
node --version
if %errorlevel% neq 0 (
    echo.
    echo ERROR: Node.js is not found. Please ensure Node.js is installed in C:\Program Files\nodejs\
    echo Or manually add Node.js to your system PATH environment variable.
    pause
    exit /b 1
)

REM Install dependencies if needed
if not exist node_modules (
    echo Installing dependencies...
    npm install
)

echo.
echo Starting frontend server...
npm start
pause
