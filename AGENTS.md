# AGENTS.md - UCD Project Guidelines

This file provides guidelines for agents operating on the UCD (Unnecessary Chdir) codebase.

## Project Overview

UCD is a Go CLI tool that wraps the shell's `cd` command with MRU (Most Recently Used) tracking and path stashing capabilities. It outputs the target path to stdout, which is consumed by the shell's builtin `cd` command.

## Build Commands

```bash
# Build the binary
go build .

# Build and install to /usr/local/bin
go build . && sudo chmod +x ucd && sudo mv ucd /usr/local/bin/ucd

# Run all tests
go test -v ./...

# Run tests in a specific package
go test -v ./records
go test -v ./utilities
go test -v ./configurations

# Run a single test
go test -v -run TestRecordHasCorrectCount ./records
go test -v -run TestSortPathRecords ./records
go test -v -run TestIsInvalidPath ./utilities
```

## Code Style Guidelines

### General

- **No comments** unless explicitly required by the user
- Use `log` for logging with prefix `[ucd]` via `util.Clog()`
- Use `fmt.Print()` for stdout output (consumed by shell's builtin `cd`)

### Imports

Group imports in the following order:
1. Standard library packages
2. Third-party packages
3. Local packages (github.com/ztcjoe93/ucd/...)

```go
import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strings"

    "github.com/jedib0t/go-pretty/table"
    "github.com/ztcjoe93/ucd/configurations"
    "github.com/ztcjoe93/ucd/records"
    util "github.com/ztcjoe93/ucd/utilities"
)
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Packages | lowercase, short | `records`, `utilities` |
| Types | PascalCase | `Records`, `PathRecord`, `Configuration` |
| Functions | PascalCase | `GetConfigurations`, `SortRecords` |
| Variables | camelCase | `aliasFlag`, `targetPath` |
| Constants | PascalCase or UPPER_SNAKE | `APPLICATION_NAME`, `MaxMRUDisplay` |
| JSON fields | PascalCase | `json:"MaxMRUDisplay"` |
| Interfaces | PascalCase | `Record` |

### Types and Structures

- Use struct tags for JSON serialization
- Use generics when applicable (see `records.SortRecords`)
- Avoid unused variables - remove or prefix with `_`

```go
type Records struct {
    PathRecords  map[string]PathRecord  `json:"paths"`
    StashRecords map[string]StashRecord `json:"stash"`
}

type PathRecord struct {
    Timestamp string `json:"ts"`
}

func SortRecords[v Record](r map[string]v) []string {
    // generic function
}
```

### Error Handling

- Use `log.Fatalf()` for fatal errors that should stop execution
- Use `util.ReturnCwd()` to exit gracefully when path resolution fails (outputs `.` to stdout)
- Ignore errors only when explicitly safe to do so (e.g., `os.UserHomeDir()`)

```go
// Fatal error
if err != nil {
    log.Fatalf("%v\n", err)
}

// Graceful exit (for cd wrapper)
if !found {
    util.Clog(fmt.Sprintf("Alias %v not found", aliasPathFlag))
    util.ReturnCwd()
}
```

### Flag Definition

Define CLI flags using the `flag` package at the start of `main()`:

```go
flag.BoolVar(&helpFlag, "h", false, "display help")
flag.BoolVar(&versionFlag, "v", false, "display ucd version")
flag.BoolVar(&listFlag, "l", false, "List history of cd paths")
flag.IntVar(&historyPathFlag, "p", 0, "Chdir to the # path from the list")
flag.StringVar(&aliasPathFlag, "pa", "", "Chdir to path with provided alias")
flag.Parse()
```

### Testing

- Use `testing` package
- Test file naming: `<package>_test.go`
- Test function naming: `Test<FunctionName>`
- Use `t.Fatalf()` for assertion failures with descriptive messages
- Table-driven tests are acceptable

```go
func TestSortPathRecords(t *testing.T) {
    pr := map[string]PathRecord{
        "path1": {Timestamp: "1990-01-01 00:00:01 +08"},
        "path3": {Timestamp: "2024-04-01 12:34:12 +08"},
    }

    records := SortRecords(pr)

    if records[0] != "path3" {
        t.Fatalf(`path3 is not the latest record`)
    }
}
```

### File Organization

```
ucd/
├── main.go                    # Entry point, flag parsing, path resolution
├── configurations/
│   └── configurations.go      # Config loading from ~/.config/ucd/ucd.conf
├── records/
│   ├── records.go             # Data structures and MRU logic
│   └── records_test.go        # Tests
├── utilities/
│   ├── utilities.go          # Helper functions
│   └── utilities_test.go     # Tests
├── go.mod
└── README.md
```

### Key Functions Reference

| Function | Location | Purpose |
|----------|----------|---------|
| `util.ReturnCwd()` | utilities.go | Exit gracefully, output `.` to stdout |
| `util.Clog()` | utilities.go | Log with `[ucd]` prefix |
| `util.IsInvalidPath()` | utilities.go | Check if path exists |
| `util.GetParentDir()` | utilities.go | Get parent directory for file fallback |
| `util.TimeNow()` | utilities.go | Get current timestamp |
| `util.AutoClear()` | utilities.go | Trim MRU records to limit |
| `r.ListRecords()` | records.go | Display MRU/stash list |
| `r.SortRecords()` | records.go | Sort records by timestamp |
| `r.AliasExists()` | records.go | Check if alias exists in stash |
| `c.GetConfigurations()` | configurations.go | Load/create config file |

### Configuration

On first run, UCD creates `~/.config/ucd/ucd.conf`:

```json
{
    "MaxMRUDisplay": -1,
    "FileFallbackBehavior": true
}
```

### Common Patterns

**Path Resolution Priority** (from main.go):
1. `-pa <alias>` - Find by stash alias
2. `-p <n>` - Navigate by MRU index
3. `-ps <n>` - Navigate by stash index
4. Argument or home directory

**Cache Storage**: `~/.ucd-cache` (JSON)

**Exit Codes**:
- `0` - Success (outputs path to stdout)
- Non-zero - Error (via `log.Fatalf` or `os.Exit`)
