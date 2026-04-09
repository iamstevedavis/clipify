# Clipify

Clipify is a small clipboard history app written in Go. It watches the clipboard, stores copied text in a local JSON history file, and exposes a tray menu action to open a history window where previous entries can be reviewed and copied back to the clipboard.

## Current status

This repo looks like an early desktop prototype.

What is implemented today:
- clipboard polling and change detection
- persistent history storage in `runtime/history.json`
- tray icon with:
  - `Show History`
  - `Quit`
- Gio-based history window
- click an item in the history window to copy it back to the clipboard

## How it works

- `internal/clipboard` polls the system clipboard and emits changes
- `internal/history` loads and saves clipboard entries as JSON
- `internal/tray` starts the system tray icon and menu
- `internal/ui` renders the clipboard history window
- `cmd/main.go` wires everything together

## Project structure

```text
.
├── assets/                # Icons and app assets
├── build.ps1              # Windows build helper
├── cmd/                   # Application entry point
├── internal/
│   ├── clipboard/         # Clipboard monitoring + read/write helpers
│   ├── history/           # JSON persistence for clipboard history
│   ├── tray/              # System tray setup and menu actions
│   └── ui/                # Gio-based history window
├── LICENSE
└── README.md
```

## Runtime behavior

At startup, Clipify:
- creates a `runtime/` directory in the current working directory if needed
- stores clipboard history in `runtime/history.json`
- loads its tray icon from `assets/icons/icon.ico`
- opens a history window when `Show History` is clicked from the tray

## Build

### Windows PowerShell

The repo includes a helper script for Windows:

```powershell
./build.ps1
```

That script builds:

```text
bin/Clipify.exe
```

and copies the `assets/` directory into `bin/assets`.

### Manual build

```sh
go build -o bin/Clipify.exe ./cmd
```

## Run

### From source

```sh
go run ./cmd
```

### From the built executable

```powershell
./bin/Clipify.exe
```

## Notes and limitations

- The current implementation is primarily shaped around a Windows-style build flow.
- The tray icon expects `assets/icons/icon.ico`.
- Clipboard history is stored as plain JSON locally.
- The app currently focuses on text clipboard history, not rich content.
- This README intentionally describes the code that is present today, not a broader future vision.

## Development ideas

Likely next improvements:
- deduplicate repeated clipboard entries
- add history limits and pruning
- improve search/filtering in the history window
- support deleting or pinning entries
- package the app cleanly for Windows and other desktop platforms

## License

This project is licensed under the [MIT License](LICENSE).
