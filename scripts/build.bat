@echo off
SETLOCAL EnableDelayedExpansion

:: Define application name and build directory
SET APP_NAME=ServerCommander
SET BUILD_DIR=..\build
SET ICON_PATH=..\src\assets\icon.ico
SET SRC_DIR=..\src

:: Check if Go is installed
where go >nul 2>nul
IF %ERRORLEVEL% NEQ 0 (
    echo Go is not installed. Please install Go first.
    exit /b 1
)

:: Ensure the "build" directory exists
IF NOT EXIST "%BUILD_DIR%" mkdir "%BUILD_DIR%"

:: Ensure all necessary Go modules are installed
echo Installing Go modules...
go mod tidy

:: Generate rsrc.syso
echo 🔨 Creating Windows resource file...
rsrc -ico "%ICON_PATH%" -o "%SRC_DIR%\rsrc.syso"

:: Build for Windows
echo 🔨 Building for Windows...
SET GOOS=windows
SET GOARCH=amd64
go build -o "%BUILD_DIR%\%APP_NAME%.exe" "%SRC_DIR%\main.go"

IF %ERRORLEVEL% EQU 0 (
    echo ✅ Windows build completed!
) ELSE (
    echo ❌ Error during Windows build
    exit /b 1
)

:: Build for Linux
echo 🐧 Building for Linux...
SET GOOS=linux
SET GOARCH=amd64
go build -o "%BUILD_DIR%\%APP_NAME%.bin" "%SRC_DIR%\main.go"

IF %ERRORLEVEL% EQU 0 (
    echo ✅ Linux build completed!
) ELSE (
    echo ❌ Error during Linux build
    exit /b 1
)

:: Build for macOS
echo 🍏 Building for macOS...
SET GOOS=darwin
SET GOARCH=amd64
go build -o "%BUILD_DIR%\%APP_NAME%.app" "%SRC_DIR%\main.go"

IF %ERRORLEVEL% EQU 0 (
    echo ✅ macOS build completed!
) ELSE (
    echo ❌ Error during macOS build
    exit /b 1
)

:: Compress the src folder (excluding build directory)
echo 📦 Compressing src directory...
powershell -Command "Get-ChildItem -Path '../' -Exclude 'build', '.github', '.gitignore', 'LICENSE' | Compress-Archive -DestinationPath '%BUILD_DIR%\src.zip'"

echo ✅ Build process completed! Files are in "%BUILD_DIR%".

ENDLOCAL
exit /b
