# mock-grpc

A simple gRPC server for todo management with support for multiple transport modes.

## Features

- gRPC server with full CRUD operations for todos
- Multiple transport modes:
  - **HTTP** (TCP) - default, works on all platforms
  - **Unix socket** - Linux/macOS only
  - **Named pipe** - Windows only
- Server-side logging for all requests
- CLI with configurable options

## Prerequisites

### Go

Install Go 1.21 or later: https://go.dev/doc/install

```bash
# Verify installation
go version
```

### Task (task runner)

Install Task: https://taskfile.dev/installation/

```bash
# macOS
brew install go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Windows (scoop)
scoop install task

# Verify installation
task --version
```

### Protocol Buffers (optional, for regenerating protos)

Only needed if you want to modify the protobuf definitions.

```bash
# macOS
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure $GOPATH/bin is in your PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Setup

```bash
# Clone the repository
git clone https://github.com/chirag-bruno/mock-grpc.git
cd mock-grpc

# Download dependencies
task tidy

# Build the binary
task build
```

## Usage

### Available Tasks

```bash
task --list
```

| Task | Description |
|------|-------------|
| `task build` | Build the server binary |
| `task run` | Run server in HTTP mode (default) |
| `task run:unix` | Run server in Unix socket mode (Linux/macOS) |
| `task run:pipe` | Run server in named pipe mode (Windows) |
| `task proto` | Regenerate Go code from protobuf |
| `task test` | Run tests |
| `task clean` | Clean build artifacts |
| `task tidy` | Tidy go modules |
| `task all` | Clean, tidy, and build |

### Running the Server

```bash
# Default: HTTP mode on localhost:50051
task run

# Or run the binary directly with options
./bin/todo-server --help
```

**CLI Options:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--mode` | `-m` | `http` | Server mode: `http`, `unix`, or `pipe` |
| `--address` | `-a` | `localhost:50051` | Server address |

**Examples:**

```bash
# HTTP mode on custom port
./bin/todo-server --mode http --address 0.0.0.0:8080

# Unix socket (Linux/macOS)
./bin/todo-server --mode unix --address /tmp/todo.sock

# Named pipe (Windows)
./bin/todo-server --mode pipe --address '\\.\pipe\todo'
```

## API

The server implements a `TodoService` with the following RPCs:

| Method | Description |
|--------|-------------|
| `CreateTodo` | Create a new todo |
| `GetTodo` | Get a todo by ID |
| `ListTodos` | List all todos |
| `UpdateTodo` | Update an existing todo |
| `DeleteTodo` | Delete a todo by ID |

See [`proto/todo.proto`](proto/todo.proto) for the full schema.

## Project Structure

```
mock-grpc/
├── cmd/server/
│   └── main.go              # CLI entrypoint
├── internal/
│   ├── server/
│   │   ├── server.go        # TodoServer implementation
│   │   └── grpc.go          # Server runner and logging interceptor
│   └── transport/
│       ├── listener.go      # Listener factory and path validation
│       ├── pipe_windows.go  # Windows named pipe support
│       └── pipe_other.go    # Unix stub
├── pkg/todo/                # Generated protobuf code
├── proto/
│   └── todo.proto           # Protobuf definitions
├── Taskfile.yml             # Task runner configuration
└── go.mod
```

## License

MIT
