@echo off
title Vigilant
echo Starting Vigilant...

if not exist "..\.venv" (
    echo [ERROR] Virtual environment not found. 
    echo Please run setup.bat first to install dependencies.
    pause
    exit /b
)

call ..\.venv\Scripts\activate
python ..\src\main.py
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Application crashed or Python is not installed.
    pause
)
