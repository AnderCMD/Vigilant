@echo off
title Vigilant - Dependency Installer
echo ==============================================
echo [ Vigilant ] Installing environment and dependencies
echo ==============================================
echo.

if not exist ".venv" (
    echo [1/3] Creating local virtual environment (.venv)...
    python -m venv .venv
) else (
    echo [1/3] Virtual environment (.venv) already exists.
)

echo [2/3] Activating environment and installing dependencies (requirements.txt)...
call .venv\Scripts\activate
pip install -r requirements.txt

echo.
echo ==============================================
echo [3/3] Installation completed!
echo You can now close this window and execute run.bat to start.
echo ==============================================
pause
