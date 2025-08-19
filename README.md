# pprof4svc

`pprof4svc` is a Go package that integrates Go runtime profiling and statistics endpoints into a Gin-based HTTP service. It provides secure access to pprof, memory, GC, and trace data with token-based authentication and randomized route prefixes.

## Features
- Exposes standard pprof endpoints (`/debug/pprof/*`) for CPU, memory, and trace profiling.
- Custom endpoints for:
    - Memory statistics (`/debug/mem`): Displays `runtime.MemStats` in text or JSON format.
    - GC statistics (`/debug/gc`): Displays `debug.GCStats` in text or JSON format.
    - Trace control (`/debug/trace`): Starts a runtime trace for a specified duration and writes binary trace data to the response.
- Token-based authentication for secure access.
- Randomized route prefixes to prevent unauthorized access.

## Installation
1. Ensure you have Go installed (version 1.21 or higher recommended).
2. Install the Gin framework:
   ```bash
   go get github.com/gin-gonic/gin
   ```
3. Add `pprof4svc` to your project by including the package files.

## Usage
1. **Initialize the Plugin**:
   ```go
   package main

   import (
       "github.com/gin-gonic/gin"
       "github.com/go-the-way/pprof4svc"
   )

   func main() {
       engine := gin.Default()
       plugin := pprof4svc.DefaultPlugin("your-secret-token")
       plugin.Plug(engine)
       engine.Run(":8080")
   }
   ```
    - Replace `"your-secret-token"` with a secure token for authentication.
    - The plugin registers all endpoints with a random prefix (e.g., `/abc123/debug/pprof/`).

2. **Access Endpoints**:
    - Access the entrypoint with the token to redirect to the pprof index:
      ```bash
      curl "http://localhost:8080/debug/pprof/?token=your-secret-token"
      ```
    - Example endpoints (with random prefix, e.g., `/abc123`):
        - `/abc123/debug/pprof/` (pprof index)
        - `/abc123/debug/pprof/profile` (CPU profile)
        - `/abc123/debug/pprof/trace` (execution trace)
        - `/abc123/debug/mem?json=true` (memory stats in JSON)
        - `/abc123/debug/gc` (GC stats in text)
        - `/abc123/debug/trace?dur=5s` (trace for 5 seconds)

3. **Analyze Trace Data**:
    - Save the response from `/debug/trace` to a file (e.g., `trace.out`):
      ```bash
      curl "http://localhost:8080/abc123/debug/trace?dur=5s" > trace.out
      ```
    - Analyze with:
      ```bash
      go tool trace trace.out
      ```

## Endpoints
- **`/debug/pprof/`**: Pprof index listing all profiling endpoints.
- **`/debug/pprof/:name`**: Specific pprof profiles (e.g., heap, goroutine).
- **`/debug/pprof/cmdline`**: Command line arguments.
- **`/debug/pprof/profile`**: CPU profile.
- **`/debug/pprof/symbol`**: Symbol lookup.
- **`/debug/pprof/trace`**: Execution trace (binary format).
- **`/debug/mem`**: Memory statistics (`runtime.MemStats`). Use `?json=true` for JSON output.
- **`/debug/gc`**: GC statistics (`debug.GCStats`). Use `?json=true` for JSON output.
- **`/debug/trace`**: Starts a trace for a duration (default: 10s, set via `?dur=5s`). Outputs binary trace data.

## Notes
- **Authentication**: Access requires a `token` query parameter matching the plugin's token.
- **Performance**: `runtime.ReadMemStats` and `debug.ReadGCStats` may trigger Stop-The-World pauses. Use sparingly in production.
- **Thread Safety**: The `/debug/trace` endpoint uses a mutex to prevent concurrent tracing.
- **Dependencies**: Requires `github.com/gin-gonic/gin` for the HTTP server.