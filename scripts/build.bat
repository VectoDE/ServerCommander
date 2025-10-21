@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "SCRIPT_DIR=%~dp0"
for %%I in ("%SCRIPT_DIR%..") do set "ROOT_DIR=%%~fI"
set "BUILD_DIR=%ROOT_DIR%\build"
set "APP_NAME=ServerCommander"
set "MAIN_PKG=%ROOT_DIR%\src"
set "ICON_PATH=%MAIN_PKG%\assets\icon.ico"
set "RESOURCE_SYSO=%MAIN_PKG%\assets\resource.syso"

echo =============================================
echo   Building %APP_NAME% (Windows batch script)
echo =============================================

call :require_tool go "https://go.dev/dl/" || goto :fail
call :require_tool git "https://git-scm.com/" || goto :fail

if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%" || goto :fail
)

set "CGO_ENABLED=0"

call :build_target windows amd64 "%BUILD_DIR%\%APP_NAME%-windows-amd64.exe" || goto :fail
call :build_target linux amd64 "%BUILD_DIR%\%APP_NAME%-linux-amd64" || goto :fail
call :build_target darwin amd64 "%BUILD_DIR%\%APP_NAME%-darwin-amd64" || goto :fail

call :archive_sources || goto :fail

echo.
echo ✅ Build completed successfully. Artifacts are located in "%BUILD_DIR%".
endlocal
exit /b 0

:fail
echo.
echo ❌ Build failed.
endlocal
exit /b 1

:require_tool
where %1 >nul 2>nul
if errorlevel 1 (
    echo %1 is required but was not found in PATH.
    echo Please install %1 from %2 and try again.
    exit /b 1
)
exit /b 0

:build_target
set "TARGET_OS=%~1"
set "TARGET_ARCH=%~2"
set "TARGET_OUTPUT=%~3"

echo.
echo Building !TARGET_OS!/!TARGET_ARCH!...

if /I "!TARGET_OS!"=="windows" (
    call :ensure_windows_resources || exit /b 1
)

set "GOOS=!TARGET_OS!"
set "GOARCH=!TARGET_ARCH!"
set "BUILD_TAG_ARGS="
if defined GO_BUILD_TAGS (
    set "BUILD_TAG_ARGS=-tags ""!GO_BUILD_TAGS!"""
)
go build !BUILD_TAG_ARGS! -o "!TARGET_OUTPUT!" "!MAIN_PKG!"
if errorlevel 1 (
    echo Failed to build !TARGET_OS!/!TARGET_ARCH!.
    exit /b 1
)

if /I "!TARGET_OS!"=="windows" (
    call :restore_windows_resources
)

exit /b 0

:ensure_windows_resources
if exist "%RESOURCE_SYSO%" (
    set "CLEAN_RESOURCE="
    exit /b 0
)

where rsrc >nul 2>nul
if errorlevel 1 (
    echo Windows resources not found and rsrc.exe is not installed.
    echo Install rsrc from https://github.com/akavel/rsrc/releases to embed the application icon.
    exit /b 1
)

echo Generating Windows resources...
rsrc -ico "%ICON_PATH%" -o "%RESOURCE_SYSO%"
if errorlevel 1 (
    echo Failed to generate Windows resources.
    exit /b 1
)
set "CLEAN_RESOURCE=1"
exit /b 0

:restore_windows_resources
if defined CLEAN_RESOURCE (
    del "%RESOURCE_SYSO%" >nul 2>nul
    set "CLEAN_RESOURCE="
)
exit /b 0

:archive_sources
set "ARCHIVE=%BUILD_DIR%\%APP_NAME%-src.zip"
echo.
echo Creating source archive...
powershell -NoLogo -NoProfile -Command ^
    "if (Test-Path '%ARCHIVE%') { Remove-Item '%ARCHIVE%' -Force };" ^
    "Compress-Archive -Path '%MAIN_PKG%\*' -DestinationPath '%ARCHIVE%'"
if errorlevel 1 (
    echo Failed to create source archive.
    exit /b 1
)
exit /b 0
