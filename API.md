# GarpScan API Reference

This document provides a detailed overview of the functions, methods, and structures available in the GarpScan SDK.

## Core Package: `garpscan`

### `type Engine`
The core scanner engine that orchestrates the concurrent worker pool.

- **`func NewEngine(options ...Option) *Engine`**
  Initializes and returns a new Engine. Accepts functional options to configure its behavior dynamically.
  *Example:* `engine := garpscan.NewEngine(garpscan.WithWorkers(100))`

- **`func WithWorkers(n int) Option`**
  A configuration option representing the maximum number of concurrent workers (goroutines) to spawn. Default is `10`.

- **`func (e *Engine) AddPlugin(p Plugin)`**
  Registers a struct that implements the `Plugin` interface into the Engine. You can register multiple plugins in a chain before starting a scan.

- **`func (e *Engine) Scan(ctx context.Context, targets []string) <-chan *Result`**
  Starts the active scanning process. For every target, the Engine asynchronously runs all registered plugins. 
  Returns a read-only channel `<-chan *Result` where outcomes are pushed stream-like. The channel is safely closed once all targets finish.

### `type Result`
Structure returned by the Engine's `Scan` channel.
```go
type Result struct {
	Target     string      `json:"target"`      // The scanned target string (e.g. "127.0.0.1")
	PluginName string      `json:"plugin_name"` // Name of the plugin that generated this result
	Status     string      `json:"status"`      // Final status of the check (e.g., "completed", "error")
	Data       interface{} `json:"data"`        // The raw output data/payload of the plugin check
	Timestamp  time.Time   `json:"timestamp"`   // The exact time the plugin finished the run
}
```

### `type Plugin`
The central interface required to build custom plugins for Engine execution.
```go
type Plugin interface {
	Name() string
	Run(ctx context.Context, target string) (*Result, error)
}
```

---

## Package: `plugins/portscanner`
Built-in simple TCP port scanning plugin.

- **`func New(ports []int, timeout time.Duration) *PortScanner`**
  Constructs a new TCP Port Scanner plugin. Takes an array limit of ports you wish to scan and a dialer connection timeout (e.g., `2 * time.Second`).

---

## Package: `reporter`
Helper package to format and process the `<-chan *garpscan.Result` values that stream out of the `Scan()` module.

- **`func ExportJSON(writer io.Writer, results <-chan *garpscan.Result) error`**
  Consumes the results channel completely and writes each object as a JSON payload row (JSON Lines format) to the provided `io.Writer` (for example `os.Stdout` or an `os.File`).

- **`func ExportCSV(writer io.Writer, results <-chan *garpscan.Result) error`**
  Consumes the results channel and parses the internal data into CSV lines format, writing them dynamically to the assigned `io.Writer`.
