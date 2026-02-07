@echo off
TITLE WebApp Aviator
CLS

ECHO ==========================================
ECHO    AVIATOR LAUNCHER
ECHO ==========================================

CD /D "%~dp0windows-server"

REM 1. Check if we have a python environment
IF EXIST "venv\Scripts\python.exe" (
    GOTO :LAUNCH
)

ECHO [SETUP] Entry point not found. Checking system Python...
where python >nul 2>nul
IF %ERRORLEVEL% NEQ 0 (
    ECHO [ERROR] Python not found in PATH!
    ECHO Please install Python from https://www.python.org/
    ECHO Make sure to check "Add Python to PATH" during installation.
    PAUSE
    EXIT /B
)

ECHO [SETUP] Creating virtual environment...
python -m venv venv
IF %ERRORLEVEL% NEQ 0 (
    ECHO [ERROR] Failed to create venv.
    PAUSE
    EXIT /B
)

ECHO [SETUP] Installing dependencies...
.\venv\Scripts\pip install -r requirements.txt
IF %ERRORLEVEL% NEQ 0 (
    ECHO [ERROR] Failed to install dependencies.
    PAUSE
    EXIT /B
)

:LAUNCH
CLS
ECHO ==========================================
ECHO    AVIATOR IS RUNNING
ECHO ==========================================
ECHO.
ECHO Use the popup window to manage apps.
ECHO Keep this terminal open.
ECHO.
ECHO [LOGS]
ECHO.

.\venv\Scripts\python.exe main.py
IF %ERRORLEVEL% NEQ 0 (
    ECHO.
    ECHO [ERROR] The application crashed with code %ERRORLEVEL%
    PAUSE
)
EXIT
