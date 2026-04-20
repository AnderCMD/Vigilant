@echo off
title Vigilant
echo Starting Vigilant...
if exist ".venv" (
    call .venv\Scripts\activate
)
python main.py
pause
