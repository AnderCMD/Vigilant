@echo off
title Vigilant - Build Executable
echo ==============================================
echo [ Vigilant ] Building Standalone Executable
echo ==============================================
echo.

if not exist ".venv" (
    echo [ERROR] Virtual environment not found. Please run setup.bat first.
    pause
    exit /b
)

echo [1/2] Activating environment and installing PyInstaller...
call .venv\Scripts\activate
pip install pyinstaller

echo.
echo [2/2] Generating .exe file in /dist folder...
pyinstaller --noconsole --onefile --icon=icon.ico --add-data "icon.ico;." --name Vigilant main.py

echo.
echo ==============================================
echo Build completed! Check the "dist" folder.
echo ==============================================
pause
