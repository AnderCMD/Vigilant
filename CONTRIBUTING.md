# Contributing to Vigilant

Thanks for your interest in contributing! This document explains how to get set up and how to submit changes.

## Getting Started

1. Fork the repository and clone your fork:
   ```bash
   git clone https://github.com/<your-username>/Vigilant.git
   cd Vigilant
   ```
2. Run the setup script to create a virtual environment and install dependencies:
   ```powershell
   .\scripts\setup.bat
   ```
3. Create a branch for your change:
   ```bash
   git checkout -b my-feature
   ```

## Development

- Source code lives in `src/`.
- Windows helper scripts live in `scripts/`.
- Run the app locally with:
  ```powershell
  .\scripts\run.bat
  ```
  or directly with `python src/main.py`.

## Submitting Changes

1. Keep pull requests focused on a single change or feature.
2. Write clear, descriptive commit messages.
3. Make sure the application still runs correctly before submitting.
4. Open a pull request describing what you changed and why.

## Reporting Bugs

Please use the [issue templates](.github/ISSUE_TEMPLATE) to report bugs or request features. Include as much detail as possible: your OS/Windows version, Python version, and steps to reproduce.

## Code Style

- Follow existing code conventions in the file you're editing.
- Prefer clear, readable code over cleverness.
- Add docstrings for new public functions, matching the existing style in `src/main.py`.

## License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
