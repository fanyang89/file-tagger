# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

File Tagger is a Go-based CLI tool that allows users to tag files with key-value pairs stored in a SQLite database. The tool supports mounting a virtual filesystem that reflects file tags, though the FUSE implementation is currently incomplete (all operations return ENOSYS).

## Architecture

### Core Components

- **`cmd/file_tagger/ft.go`**: Main CLI application using urfave/cli/v3
- **`ft/v1/`**: Core business logic package
  - `db.go`: Database operations and Tagger struct
  - `file.go`: FileEntry model (ID, Path)
  - `tag.go`: Tag model (ID, FileID, Name, Value) and parsing logic
  - `fs.go`: FUSE filesystem interface (stub implementation)
- **`pzlog/`**: Custom logging package with pterm-based pretty console output

### Data Model

The application uses GORM with SQLite to store:

- **FileEntry**: Represents files with unique paths
- **Tag**: Key-value tags associated with files via foreign key

### Configuration

- Config file location: `~/.config/file-tagger/config.yaml`
- Required field: `dsn` (SQLite database connection string)
- Auto-creates config file if it doesn't exist

## Common Commands

### Build

```bash
go build ./cmd/file_tagger
```

### Run

```bash
# Build and run the CLI tool
go run ./cmd/file_tagger [command]

# Example usage
go run ./cmd/file_tagger tag /path/to/file key=value
go run ./cmd/file_tagger show /path/to/file
go run ./cmd/file_tagger clear /path/to/file
go run ./cmd/file_tagger delete /path/to/file key
go run ./cmd/file_tagger mount /mount/point
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./ft/v1
```

### Database Migration

The application automatically runs GORM migrations on startup using `tagger.Migrate()`.

## CLI Commands

- `tag <path> <tags...>`: Add tags to a file (format: key=value)
- `show <path>`: Display all tags for a file
- `clear <path>`: Remove all tags from a file
- `delete <path> <tags...>`: Remove specific tags from a file
- `mount <mountpoint>`: Mount FUSE filesystem (currently stub implementation)

## Development Notes

- Uses Go 1.24.4 with go.work for multi-module workspace
- Includes local pzlog package replacement in go.mod
- Logging uses zerolog with custom pterm writer for colored output
- Error handling uses cockroachdb/errors for wrapped error messages
- All file paths are resolved to absolute paths before storage
- FUSE filesystem implementation exists but is incomplete - all operations return ENOSYS

## Dependencies

Key libraries:

- `github.com/urfave/cli/v3`: CLI framework
- `gorm.io/gorm`: ORM with SQLite driver
- `github.com/winfsp/cgofuse`: FUSE bindings
- `github.com/rs/zerolog`: Structured logging
- `github.com/fioepq9/pzlog`: Custom pretty console logging (local package)
