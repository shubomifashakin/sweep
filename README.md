# sweep

A CLI tool that automatically organizes files in a directory by sorting them into categorized subfolders.

## What it does

`sweep` scans a directory and moves files into the following subfolders based on their extension:

| Folder       | Extensions              |
| ------------ | ----------------------- |
| `images/`    | `.jpg`, `.png`, `.gif`  |
| `videos/`    | `.mp4`, `.avi`, `.mkv`  |
| `documents/` | `.pdf`, `.doc`, `.docx` |
| `others/`    | Everything else         |

After every sweep, a timestamped JSON report is saved to `~/.sweep/` detailing every file that was moved and any failures.

## Installation

### Download binary (recommended)

Download the latest binary for your platform from the [releases page](https://github.com/shubomifashakin/sweep/releases):

| Platform              | Binary                    |
| --------------------- | ------------------------- |
| Windows               | `sweep-windows-amd64.exe` |
| macOS (Intel)         | `sweep-darwin-amd64`      |
| macOS (Apple Silicon) | `sweep-darwin-arm64`      |
| Linux                 | `sweep-linux-amd64`       |

### Build from source

Requires Go 1.21+

```bash
git clone https://github.com/shubomifashakin/sweep
cd sweep
go build -o sweep
```

## Usage

```bash
sweep --dir <path-to-directory> [flags]
```

### Flags

| Flag        | Type   | Default  | Description                    |
| ----------- | ------ | -------- | ------------------------------ |
| `--dir`     | string | required | Path to the directory to sweep |
| `--verbose` | bool   | false    | Enable verbose/debug logging   |

### Examples

**Organize your Downloads folder:**

```bash
sweep --dir ~/Downloads
```

**Organize with verbose output to see every step:**

```bash
sweep --dir ~/Downloads --verbose
```

**Organize a specific project folder:**

```bash
sweep --dir /path/to/messy/folder
```

## Output

Running sweep produces structured JSON logs to stdout:

```json
{
  "level": "info",
  "timestamp": "2026-05-20T10:30:00.000Z",
  "msg": "Report saved",
  "path": "/Users/john/.sweep/sweep-2026-05-20-10-30-00.json"
}
```

With `--verbose` enabled you'll also see debug logs for every file processed:

```json
{"level":"debug","timestamp":"2026-05-20T10:30:00.000Z","msg":"Moving file","file":"invoice.pdf"}
{"level":"debug","timestamp":"2026-05-20T10:30:00.000Z","msg":"Moving file","file":"photo.jpg"}
```

## Reports

After every sweep a JSON report is automatically saved to `~/.sweep/`:

```
~/.sweep/
  sweep-2026-05-20-10-30-00.json
  sweep-2026-05-19-09-15-00.json
```

Example report:

```json
{
  "total_files": 10,
  "date": "2026-05-20 10:30:00",
  "directory": "/Users/john/Downloads",
  "success": 9,
  "successfulMoves": [
    { "source": "invoice.pdf", "destination": "documents/invoice.pdf" },
    { "source": "photo.jpg", "destination": "images/photo.jpg" }
  ],
  "failures": 1,
  "failedMoves": [
    { "source": "locked-file.pdf", "destination": "documents/locked-file.pdf" }
  ]
}
```

## Notes

- Files with no extension are skipped
- `.ini` files are skipped
- Files already inside `images/`, `videos/`, `documents/` or `others/` subfolders are skipped automatically
- The target subfolders are created automatically if they don't exist

## Built with

- [Go](https://golang.org)
- [Uber Zap](https://github.com/uber-go/zap) — structured logging
