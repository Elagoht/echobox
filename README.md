# Echobox

A simple HTTP echo server for testing and debugging HTTP requests.

## Installation

```bash
go install github.com/Elagoht/echobox/cmd/echobox@latest
```

## Usage

Run directly:
```bash
go run github.com/Elagoht/echobox/cmd/echobox@latest
```

Or if installed:
```bash
echobox
```

### Configuration

The server will start on port 5867 by default. You can customize the port and timeouts using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `5867` |
| `READ_TIMEOUT` | Read timeout in seconds | `30` |
| `WRITE_TIMEOUT` | Write timeout in seconds | `30` |

```bash
PORT=8080 echobox
PORT=3000 READ_TIMEOUT=60 WRITE_TIMEOUT=60 echobox
```

### Using Make

If you have cloned the repository, you can use the Makefile:

```bash
make build    # Build the binary
make run      # Run the application
make test     # Run tests
make install  # Install to GOPATH/bin
make help     # Show all available targets
```

## Endpoints

| Endpoint | Description |
|----------|-------------|
| `/` | Full echo (method, path, query, headers, body) |
| `/headers` | Returns only the request headers |
| `/body` | Returns the request body as-is |
| `/queries` | Returns only the query parameters |
| `/200-699` | Any 3-digit status code (e.g., `/404`, `/500`) |

All endpoints accept any HTTP method (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS).

## Examples

```bash
# Get full request echo
curl localhost:5867

# Get headers
curl -H "Custom-Header: test" localhost:5867/headers

# Get body
curl -X POST -d '{"message":"hello"}' localhost:5867/body

# Get query parameters
curl "localhost:5867/queries?foo=bar&baz=qux"

# Test status codes
curl localhost:5867/404
curl localhost:5867/500
```

## Project Structure

```
echobox/
├── cmd/
│   └── echobox/          # Application entry point
│       └── main.go
├── internal/
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── handler/          # HTTP handlers
│   │   └── handler.go
│   └── router/           # Routing setup
│       └── router.go
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

## Development

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Vet code
make vet

# Build
make build
```

## License

MIT License - see [LICENSE](LICENSE) for details.
