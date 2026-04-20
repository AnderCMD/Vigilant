import pystray
from PIL import Image
import pyautogui
import threading
import time
import sys
import os

# Disable pyautogui's FAILSAFE to prevent errors when the mouse is in the screen corners
pyautogui.FAILSAFE = False

# Global state variables
is_running = True
is_paused = False
CHECK_INTERVAL = 180  # 3 minutes in seconds (180s)


def get_resource_path(relative_path: str) -> str:
    """Gets the absolute path to a resource, compatible with PyInstaller.

    Args:
        relative_path: The relative path to the resource file.

    Returns:
        The absolute path to the resource.
    """
    try:
        # PyInstaller unpacks temporary files into sys._MEIPASS
        base_path = sys._MEIPASS
    except AttributeError:
        # Move up one level from src/ to reach the project root
        # This assumes main.py is in the src/ directory
        base_path = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    
    return os.path.join(base_path, relative_path)


# Load icons
icon_path = get_resource_path(os.path.join('assets', 'icon.ico'))
base_img = Image.open(icon_path)

# Active: Full color icon
icon_active = base_img.copy()
# Paused: Grayscale icon
icon_paused = base_img.convert("L")


def vigil_loop():
    """Background monitoring loop that presses F15 periodically.

    This function runs in a separate thread. It checks every second for
    program exit or pause state, but triggers the keypress only at the 
    specified interval.
    """
    global is_running, is_paused
    while is_running:
        # Sleep in 1-second chunks to keep the app responsive to user input
        for _ in range(CHECK_INTERVAL):
            if not is_running:
                return
            time.sleep(1)

        # Trigger keypress if not paused
        if not is_paused:
            try:
                # F15 is a non-disruptive key that prevents system idle/sleep
                pyautogui.press('f15')
            except Exception:
                # Ignore errors silently to prevent the app from crashing
                pass


def toggle_pause(icon: pystray.Icon, item: pystray.MenuItem):
    """Toggles the application between active and paused states.

    Updates the tray icon visually and changes the icon title.

    Args:
        icon: The system tray icon instance.
        item: The menu item that triggered this callback.
    """
    global is_paused
    is_paused = not is_paused

    if is_paused:
        icon.icon = icon_paused
        icon.title = "Vigilant - Paused"
    else:
        icon.icon = icon_active
        icon.title = "Vigilant - Active"


def quit_app(icon: pystray.Icon, item: pystray.MenuItem):
    """Stops the background loop and closes the application.

    Args:
        icon: The system tray icon instance.
        item: The menu item that triggered this callback.
    """
    global is_running
    is_running = False
    icon.stop()


def setup_gui():
    """Configures and runs the system tray application.

    Sets up the menu items, initializes the background thread, and
    starts the icon loop.
    """
    # Menu structure (Right-click menu)
    menu = pystray.Menu(
        pystray.MenuItem(
            lambda text: "Resume" if is_paused else "Pause", 
            toggle_pause
        ),
        pystray.Menu.SEPARATOR,
        pystray.MenuItem('Exit', quit_app)
    )

    # Initialize the tray icon
    icon = pystray.Icon("Vigilant", icon_active, "Vigilant - Active", menu)

    # Start the key-pressing thread in the background
    background_thread = threading.Thread(target=vigil_loop)
    background_thread.daemon = True  # Thread will exit when the main app exits
    background_thread.start()

    # Launch the system tray icon (blocks until icon.stop() is called)
    icon.run()


if __name__ == "__main__":
    setup_gui()

