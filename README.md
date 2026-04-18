# GarpScan 🚀

[![Go](https://github.com/tommyvercetti89/garpscan/actions/workflows/test.yml/badge.svg)](https://github.com/tommyvercetti89/garpscan/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GarpScan is a high-performance, concurrent, and modular open-source cybersecurity scanning framework SDK written in Go.

It acts as a library that developers can import into their own tools to dispatch asynchronous scanning tasks with custom plugins, easily integrating with any existing Golang CLI or backend project.

## Features ✨
- **Modular Plugin Engine**: Easily write and inject your own custom scanner behaviors (e.g. VulnCheckers, BannerGrabbers).
- **Concurrent Worker Pool**: Execute thousands of target checks using an asynchronous worker pool that eliminates Goroutine leaks memory spikes.
- **Context-Aware Architecture**: Built-in graceful cancellation (`CTRL+C` / Timeouts) across all workers and dialers.
- **Data Exporting Module**: Streams data sequentially to be easily converted and written into JSON or CSV structures payload files.
- **Include Battery Plugins**: Comes shipped with a configurable default TCP Port Scanner plugin.

## Installation 📦

```bash
go get -u github.com/tommyvercetti89/garpscan
```

## Basic Usage (SDK) 🛠️

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tommyvercetti89/garpscan"
	"github.com/tommyvercetti89/garpscan/plugins/portscanner"
	"github.com/tommyvercetti89/garpscan/reporter"
)

func main() {
    // 1. Initialize the engine with a 50 routines worker pool
	engine := garpscan.NewEngine(garpscan.WithWorkers(50))
	
	// 2. Register your plugins (Built-in TCP port scanner)
	engine.AddPlugin(portscanner.New([]int{22, 80, 443, 3306}, 2*time.Second))

	// 3. Begin scan and listen asynchronously
	results := engine.Scan(context.Background(), []string{"scanme.nmap.org"})
	
	// 4. Output results directly to Stdout in JSON Format
	reporter.ExportJSON(os.Stdout, results)
}
```

## Examples 📂
For more detailed implementations, check the [examples/](examples/) directory.
- [Basic Scan Example](examples/basic_scan/main.go)

## Development & Testing 🏗️
This project includes a `Makefile` to simplify common development tasks:

```bash
# Run all tests
make test

# Tidy module dependencies
make tidy

# Run the basic example
make example
```

## Creating Your Own Plugins 🔌
You can create a custom plugin by implementing the `garpscan.Plugin` Interface in your tool:
```go
type Plugin interface {
	Name() string
	Run(ctx context.Context, target string) (*garpscan.Result, error)
}
```

## Documentation
- [Usage Guide](GUIDE.md) — Start here! A beginner-friendly guide to GarpScan.
- [API Reference](API.md) — Technical details for all functions and methods.

## License 📄
Distributed under the MIT License. See `LICENSE` for more information.
