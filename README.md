# Vigilant 👁️

Vigilant is a lightweight system tray application for Windows designed to keep your system active and prevent it from going into sleep mode or changing your status to "Away".

## How it works

Vigilant runs silently in the background and periodically (every 3 minutes by default) simulates a press of the **F15** key. 

The **F15** key is a non-disruptive key that is not present on most physical keyboards, meaning it won't interfere with your typing or any open applications, but it is recognized by the operating system as user activity.

## Features

- **System Tray Integration**: Easily accessible from the taskbar.
- **Visual Feedback**: The icon changes to grayscale when paused and color when active.
- **Easy Control**: Right-click the icon to Pause, Resume, or Exit.
- **Lightweight**: Minimal CPU and memory footprint.
- **Auto-stop**: Gracefully closes when you exit the app.

## Installation

### Prerequisites

- [Python 3.9+](https://www.python.org/downloads/)
- Windows (Primary target for F15 key simulation)

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/AnderCMD/Vigilant.git
   cd Vigilant
   ```

2. **Run the setup script**:
   Execute `setup.bat` on Windows to create a virtual environment and install the required dependencies:
   ```powershell
   .\setup.bat
   ```

## Usage

After installation, you can run the application using the convenient batch script:
```powershell
.\run.bat
```

Alternatively, you can run it directly with Python:
```bash
python main.py
```

## Packaging

### Building the Executable

If you want to create a standalone `.exe` file for Windows:

1. **Install PyInstaller**:
   ```bash
   pip install pyinstaller
   ```

2. Run the provided `build.bat` for an automated build:
   ```powershell
   .\build.bat
   ```

   Or run the command manually:
   ```bash
   pyinstaller --noconsole --onefile --icon=icon.ico --add-data "icon.ico;." --name Vigilant main.py
   ```
   - `--noconsole`: Prevents a terminal window from appearing.
   - `--onefile`: Bundles everything into a single `.exe`.
   - `--add-data`: Includes the icon file inside the executable.

The generated file will be in the `dist/` folder.

## Dependencies

- `pyautogui`: For simulating key presses.
- `pystray`: For creating the system tray icon.
- `Pillow`: For image handling and icon manipulation.

## Security

Vigilant does **not** collect any data, connect to the internet, or track your keystrokes. It only generates a single keypress at a fixed interval.

## License

MIT License - feel free to use and modify for personal or commercial projects.