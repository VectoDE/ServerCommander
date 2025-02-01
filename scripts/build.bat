@echo off
SETLOCAL EnableDelayedExpansion

:: Define application name and directories
SET APP_NAME=ServerCommander
SET BUILD_DIR=..\build
SET SRC_DIR=..\src
SET ICON_PATH=%SRC_DIR%\assets\icon.ico
SET RESOURCE_FILE=%SRC_DIR%\rsrc.syso

:: Function to check if Go is installed
echo Checking if Go is installed...

:: Check if Go is installed
where go >nul 2>nul
IF %ERRORLEVEL% NEQ 0 (
    echo Go is not installed. Downloading Go...
    call :InstallGo
)

:: Check if Git is installed (necessary for downloading dependencies)
where git >nul 2>nul
IF %ERRORLEVEL% NEQ 0 (
    echo Git is not installed. Please download and install Git from https://git-scm.com/
    exit /b 1
) ELSE (
    echo Git is already installed.
)

:: Ensure the "build" directory exists
IF NOT EXIST "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
)

:: Ensure all necessary Go modules are installed
echo Installing Go modules...
go mod tidy

:: Create and set the cache folder
SET CACHE_DIR=%USERPROFILE%\AppData\Local\Temp\build_cache
IF NOT EXIST "%CACHE_DIR%" (
    mkdir "%CACHE_DIR%"
)
SET GOPATH=%CACHE_DIR%

:: Convert icon.ico to rsrc.syso (inside src)
echo 🔨 Creating resource file (rsrc.syso) with icon...
rsrc -ico "%ICON_PATH%" -o "%RESOURCE_FILE%"

:: Verify that rsrc.syso was created successfully
IF NOT EXIST "%RESOURCE_FILE%" (
    echo ❌ Error: rsrc.syso was not created. Build cannot continue.
    exit /b 1
)

:: Build for Windows
echo 🔨 Building for Windows (.exe)...
SET GOOS=windows
SET GOARCH=amd64
go build -o "%BUILD_DIR%\%APP_NAME%.exe" ../src/main.go

:: Check if the Windows build was successful
IF %ERRORLEVEL% EQU 0 (
    echo ✅ Windows build completed!
) ELSE (
    echo ❌ Error during Windows build
    exit /b 1
)

:: Build for Linux
echo 🐧 Building for Linux (.bin)...
SET GOOS=linux
SET GOARCH=amd64
go build -o "%BUILD_DIR%\%APP_NAME%.bin" ../src/main.go

:: Check if the Linux build was successful
IF %ERRORLEVEL% EQU 0 (
    echo ✅ Linux build completed!
) ELSE (
    echo ❌ Error during Linux build
    exit /b 1
)

:: Build for macOS
echo 🍏 Building for macOS (.app)...
SET GOOS=darwin
SET GOARCH=amd64
go build -o "%BUILD_DIR%\%APP_NAME%.app" ../src/main.go

:: Check if the macOS build was successful
IF %ERRORLEVEL% EQU 0 (
    echo ✅ macOS build completed!
) ELSE (
    echo ❌ Error during macOS build
    exit /b 1
)

:: Compress the entire src folder, excluding the build directory
echo 📦 Compressing src directory (excluding build)...

:: Remove existing src.zip if it exists
IF EXIST "%BUILD_DIR%\src.zip" (
    echo Deleting existing src.zip...
    del "%BUILD_DIR%\src.zip"
)

:: Compress the src folder
powershell -Command "Compress-Archive -Path '%SRC_DIR%\*' -DestinationPath '%BUILD_DIR%\src.zip' -Force"

:: Final message
echo ✅ Build process completed for all platforms! Your build files and the compressed src folder are located in the "%BUILD_DIR%" folder.

ENDLOCAL
exit /b

:InstallGo
:: Install Go by downloading it from the official site
echo Downloading Go from the official website...

:: Download Go binary (Windows 64-bit as an example)
powershell -Command "Invoke-WebRequest -Uri https://go.dev/dl/go1.20.5.windows-amd64.msi -OutFile '%TEMP%\go_installer.msi'"

:: Install Go using MSI package
msiexec /i "%TEMP%\go_installer.msi" /quiet

:: Wait for Go installation to complete
echo Waiting for Go installation to complete...
timeout /t 10 /nobreak

:: Check if Go is installed correctly now
where go >nul 2>nul
IF %ERRORLEVEL% EQU 0 (
    echo Go has been successfully installed!
) ELSE (
    echo Failed to install Go. Please install it manually from https://golang.org/dl/
    exit /b 1
)

:: Add Go bin path to system environment variables (if not present)
SET GO_BIN_PATH=%USERPROFILE%\Go\bin
echo Checking if Go bin path is in system environment variables...

:: Check if Go bin directory is in PATH
echo %PATH% | findstr /I "%GO_BIN_PATH%" >nul
IF %ERRORLEVEL% NEQ 0 (
    echo Adding Go bin path to system environment variables...
    setx PATH "%PATH%;%GO_BIN_PATH%"
    echo Go bin path has been added to PATH.
) ELSE (
    echo Go bin path is already in PATH.
)

:: Clean up the installer
del "%TEMP%\go_installer.msi"
exit /b