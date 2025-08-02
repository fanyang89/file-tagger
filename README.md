# File Tagger

A Go-based CLI tool that allows you to tag files with key-value pairs and mount them as a virtual filesystem organized by tags.

## Features

- **File Tagging**: Add multiple key-value tags to any file
- **Tag Management**: View, clear, and delete tags from files
- **Virtual Filesystem**: Mount files organized by their tags using FUSE
- **SQLite Database**: Persistent storage for file metadata
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation

### Prerequisites

- Go 1.24.4 or later
- FUSE compatible filesystem (WinFsp on Windows, fuse macOS/macFUSE on macOS, libfuse on Linux)

### Build from Source

```bash
git clone <repository-url>
cd file-tagger
go build ./cmd/file_tagger
```

## Quick Start

### 1. Tag Files

```bash
# Tag a file with multiple key-value pairs
./file_tagger tag document.txt category=report type=pdf status=draft

# Tag multiple files
./file_tagger tag image.jpg category=image format=jpeg
./file_tagger tag spreadsheet.xlsx category=data format=xlsx
```

### 2. View Tags

```bash
# Show all tags for a file
./file_tagger show document.txt
# Output:
# category=report
# type=pdf
# status=draft
```

### 3. Manage Tags

```bash
# Remove specific tags
./file_tagger delete document.txt status type

# Clear all tags from a file
./file_tagger clear document.txt
```

### 4. Mount Virtual Filesystem

```bash
# Create mount point
mkdir ./mnt

# Mount filesystem
./file_tagger mount ./mnt

# Browse files by tags
ls ./mnt/
# Output:
# category  format  type

ls ./mnt/category/
# Output:
# document.txt  image.jpg  spreadsheet.xlsx

# Files appear as symlinks to their original locations
ls -la ./mnt/category/document.txt
# Output:
# document.txt -> /path/to/original/document.txt

# Access files through the virtual filesystem
cat ./mnt/category/document.txt
```

## Configuration

The tool uses a configuration file located at `~/.config/file-tagger/config.yaml`:

```yaml
dsn: "file-tagger.db"
```

The `dsn` field specifies the SQLite database file path. The configuration file is automatically created on first run.

## Command Reference

### `tag`

Add tags to a file.

```bash
./file_tagger tag <path> <tags...>
```

**Arguments:**

- `path`: File path to tag
- `tags`: One or more key-value pairs in format `key=value`

**Example:**

```bash
./file_tagger tag report.pdf category=document type=pdf year=2024
```

### `show`

Display all tags for a file.

```bash
./file_tagger show <path>
```

**Arguments:**

- `path`: File path to show tags for

**Example:**

```bash
./file_tagger show report.pdf
```

### `clear`

Remove all tags from a file.

```bash
./file_tagger clear <path>
```

**Arguments:**

- `path`: File path to clear tags from

**Example:**

```bash
./file_tagger clear report.pdf
```

### `delete`

Remove specific tags from a file.

```bash
./file_tagger delete <path> <tags...>
```

**Arguments:**

- `path`: File path to delete tags from
- `tags`: One or more tag names to delete

**Example:**

```bash
./file_tagger delete report.pdf year type
```

### `mount`

Mount the virtual filesystem.

```bash
./file_tagger mount <mountpoint>
```

**Arguments:**

- `mountpoint`: Directory to mount the filesystem

**Example:**

```bash
./file_tagger mount ./mnt
```

## Virtual Filesystem Structure

The mounted filesystem organizes files by their tags:

```bash
/mnt/
├── category/
│   ├── document.txt -> /path/to/document.txt
│   ├── image.jpg -> /path/to/image.jpg
│   └── spreadsheet.xlsx -> /path/to/spreadsheet.xlsx
├── type/
│   ├── pdf -> /path/to/document.txt
│   ├── jpeg -> /path/to/image.jpg
│   └── xlsx -> /path/to/spreadsheet.xlsx
└── status/
    └── draft -> /path/to/document.txt
```

- **Root directory**: Contains directories for each unique tag name
- **Tag directories**: Contain files that have that specific tag
- **Files**: Appear as symlinks to their original locations

## Architecture

### Core Components

- **`cmd/file_tagger/ft.go`**: Main CLI application using urfave/cli/v3
- **`ft/v1/`**: Core business logic package
  - `db.go`: Database operations and Tagger struct
  - `file.go`: FileEntry model (ID, Path)
  - `tag.go`: Tag model (ID, FileID, Name, Value) and parsing logic
  - `fs.go`: FUSE filesystem interface implementation

### Data Model

The application uses GORM with SQLite to store:

- **FileEntry**: Represents files with unique paths
- **Tag**: Key-value tags associated with files via foreign key

## Dependencies

- `github.com/urfave/cli/v3`: CLI framework
- `gorm.io/gorm`: ORM with SQLite driver
- `github.com/winfsp/cgofuse`: FUSE bindings
- `github.com/rs/zerolog`: Structured logging
- `github.com/fioepq9/pzlog`: Custom pretty console logging

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./ft/v1
```

### Building

```bash
# Build the application
go build ./cmd/file_tagger

# Build for different platforms
GOOS=linux GOARCH=amd64 go build ./cmd/file_tagger
GOOS=darwin GOARCH=amd64 go build ./cmd/file_tagger
GOOS=windows GOARCH=amd64 go build ./cmd/file_tagger
```

## Platform-Specific Notes

### Windows

- Requires [WinFsp](https://winfsp.dev/rel/) to be installed
- Use Windows paths with backslashes or forward slashes

### macOS

- Requires [macFUSE](https://osxfuse.github.io/) to be installed
- May require additional permissions for filesystem access

### Linux

- Requires `fuse` package to be installed
- User may need to be in the `fuse` group

## License

MIT
