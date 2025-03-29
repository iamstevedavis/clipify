# Clipify

Clipify is a clipboard history manager built with Go. It allows users to monitor clipboard activity, save clipboard history, and view it through a system tray interface.

## Features

- **Clipboard Monitoring**: Automatically tracks clipboard changes.
- **History Management**: Saves clipboard history to a JSON file for persistence.
- **System Tray Integration**: Provides a tray icon with options to view history or quit the application.
- **Cross-Platform Support**: Works on Windows, macOS, and Linux.

## Project Structure

```
.
├── assets/             # Contains application assets like icons
├── bin/                # Build output directory
├── cmd/                # Main application entry point
├── internal/           # Internal packages
│   ├── clipboard/      # Clipboard monitoring logic
│   ├── history/        # History management logic
│   ├── tray/           # System tray integration
│   ├── ui/             # User interface logic
├── runtime/            # Runtime files (e.g., history.json)
├── .vscode/            # VS Code configuration files
├── build.ps1           # Build script for Windows
├── go.mod              # Go module file
├── go.sum              # Go dependencies checksum
├── LICENSE             # License file
├── README.md           # Project documentation
```

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/your-username/clipify.git
   cd clipify
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Build the project:

   ```sh
   go build -o bin/Clipify.exe cmd/main.go
   ```

## Usage

1. Run the application:

   ```sh
   ./bin/Clipify.exe
   ```

2. Access the system tray icon to view clipboard history or quit the application.

## Configuration

- **Data Directory**: Stores runtime files like `history.json`. Default: `runtime/`.
- **Assets Directory**: Stores application assets like icons. Default: `assets/icons`.

## Development

### Prerequisites

- Go 1.20 or later
- A code editor (e.g., Visual Studio Code)

### Running Locally

1. Start the application:

   ```sh
   go run cmd/main.go
   ```

2. Modify the code in the `internal/` directory to add new features or fix bugs.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

- [systray](https://github.com/getlantern/systray) for system tray integration.