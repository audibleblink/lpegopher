# lpegopher - PE Analysis and Privilege Escalation Path Discovery

lpegopher is a Go application that collects import/export data from PE files, file-system security descriptors, and automatic file runners to determine privilege escalation paths and persistence mechanisms via relationship graphs using Neo4j.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Prerequisites and Setup
- Install Go 1.23 or higher (minimum Go 1.23.0 required)
- The application has two primary modes: collector (Windows-only) and processor (cross-platform)
- Neo4j database required for the processor mode

### Building the Application
**NEVER CANCEL builds** - Build times are predictable and should complete successfully.

- Clean dependencies: `go mod tidy` (under 1 second after initial download)
- Build for Linux: `go build -o lpegopher .` - takes approximately **1 second, NEVER CANCEL**. Set timeout to 60+ seconds.
- Build for Windows: `GOOS=windows GOARCH=amd64 go build -o lpegopher.exe .` - takes approximately **1 second, NEVER CANCEL**. Set timeout to 60+ seconds.

### Testing
**NEVER CANCEL tests** - Test suites are fast and reliable.

- Run all tests: `go test ./...` - takes approximately **1 second, NEVER CANCEL**. Set timeout to 30+ seconds.
- Tests cover: collectors, cypher, node, processor, and util packages
- All tests pass consistently on Linux environments

### Code Quality and Linting
- Run formatting: `go fmt ./...` (fixes code formatting, under 1 second)
- Run static analysis: `go vet ./...` - takes approximately **under 1 second, NEVER CANCEL**. Set timeout to 30+ seconds.
- No golangci-lint configuration exists; use standard Go tools

### Running the Application

#### Command Structure
```
lpegopher [--debug] [--nocolor] <command> [<args>]

Commands:
  getsystem    Utility for acquiring SYSTEM before collection (Windows-only)
  collect      Collect Windows PE and Runner data (Windows-only)
  process      Run Post-Processing tasks and populate neo4j (cross-platform)
```

#### Collector Mode (Windows-only)
```bash
# This command only works on Windows
./lpegopher collect <root_dir>
```
- Collects PEs, file tree, OS Principals, and Runners from Windows systems
- Outputs CSV files for later processing
- Sources include: Services, Run Keys, Tasks, Currently running processes

#### Processor Mode (Cross-platform)
```bash
# Basic usage with local Neo4j
./lpegopher process /path/to/csv/data

# With database connection options
./lpegopher process --user neo4j --pass password --host localhost --port 7687 --db neo4j /path/to/csv/data

# Drop database before processing
./lpegopher process --drop /path/to/csv/data

# Serve files via HTTP instead of Neo4j import directory
./lpegopher process --http localhost:8888 /path/to/csv/data
```

## Database Configuration
The processor requires Neo4j with the following default connection settings:
- Protocol: bolt
- Host: localhost
- Port: 7687
- Username: neo4j
- Password: neo4j
- Database: neo4j

Override via CLI flags or environment variables: NEO_USER, NEO_PASSWORD, NEO_HOST, NEO_PORT, NEO_DBNAME, NEO_PROTO

## Validation Scenarios

### Building and Testing Validation
Always validate changes by running this complete sequence:
1. `go mod tidy` - verify dependencies resolve
2. `go build -o lpegopher .` - verify Linux build succeeds
3. `GOOS=windows GOARCH=amd64 go build -o lpegopher.exe .` - verify Windows cross-compilation
4. `go test ./...` - verify all tests pass
5. `go vet ./...` - verify static analysis passes
6. `go fmt ./...` - apply standard formatting

### Functional Validation
After making changes to the CLI or core functionality:
1. Test help system: `./lpegopher --help`
2. Test subcommand help: `./lpegopher process --help`
3. Test debug output: `./lpegopher --debug process /tmp/empty 2>&1 | head -10`
   - Should show connection error to Neo4j (expected when database not running)
   - Validates argument parsing and initialization code paths

### When to Run Full Validation
- **ALWAYS** run the complete build and test sequence before committing changes
- **ALWAYS** test both Linux and Windows builds when changing core functionality
- **ALWAYS** run functional validation when modifying CLI argument parsing or main application flow

## Key Projects and Code Structure

### Directory Structure
```
.
├── args/           # Command-line argument definitions and parsing
├── collectors/     # Windows PE and runner data collection (Windows-only)
├── cypher/         # Neo4j query building and database operations
├── node/           # Graph node definitions and schema management
├── processor/      # CSV data processing and Neo4j population
├── util/           # Shared utility functions
├── main.go         # Application entry point and command routing
├── controller*.go  # Platform-specific controller implementations
└── queries.md      # Example Cypher queries for analysis
```

### Important Files to Review When Making Changes
- `args/args.go` - When modifying command-line interface
- `main.go` - When changing application flow or adding new commands
- `cypher/cypher.go` - When modifying database operations
- `processor/pe.go` and `processor/runner.go` - When changing data processing logic
- `node/node.go` - When modifying graph schema or node definitions

### Test Data and Examples
- `collectors/testdata/mock_runners.go` - Contains test data generators
- `queries.md` - Contains example Cypher queries for analysis
- `.vscode/launch.json` - Contains example command-line usage patterns

## Development Notes

### Platform Considerations
- Collection functionality (`collect` command) only works on Windows
- Processing functionality (`process` command) works on all platforms
- Cross-compilation to Windows works reliably from Linux

### Dependencies
- Uses Go modules for dependency management
- Key dependencies: Neo4j Go driver, go-arg for CLI parsing, platform-specific Windows APIs
- No external build tools required beyond Go toolchain

### Debugging
- Use `--debug` flag for verbose output
- Use `--nocolor` flag to disable colored output for log parsing
- VS Code launch configurations available in `.vscode/launch.json`

## Common Operations

### Adding New CLI Commands
1. Modify `args/args.go` to add new command structure
2. Update `main.go` to handle new command in switch statement
3. Implement command logic in appropriate package
4. **ALWAYS** test with `./lpegopher <newcommand> --help`
5. **ALWAYS** run full validation sequence

### Modifying Database Schema
1. Update node definitions in `node/` package
2. Modify Cypher queries in `cypher/` package
3. Update processor logic in `processor/` package
4. **ALWAYS** test with `--drop` flag to recreate schema
5. **ALWAYS** run tests to verify schema operations

### Adding New Data Collectors (Windows-only)
1. Add collector logic to `collectors/` package
2. Update `controller_windows.go` for new collection types
3. Ensure CSV output format matches processor expectations
4. **ALWAYS** add corresponding tests
5. **CANNOT VALIDATE** collection on Linux - document Windows testing requirements