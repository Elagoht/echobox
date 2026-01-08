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

The server will start on port 5867 by default. You can customize the port using the `PORT` environment variable:

```bash
PORT=8080 echobox
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
