package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tommyvercetti89/garpscan"
)

// ExamplePlugin is a simple implementation of garpscan.Plugin
type ExamplePlugin struct{}

func (p *ExamplePlugin) Name() string {
	return "example-plugin"
}

func (p *ExamplePlugin) Run(ctx context.Context, target string) (*garpscan.Result, error) {
	// Simulate some scanning work
	time.Sleep(100 * time.Millisecond)
	
	return &garpscan.Result{
		Target:     target,
		PluginName: p.Name(),
		Status:     "success",
		Data:       "Hello from example plugin!",
		Timestamp:  time.Now(),
	}, nil
}

func main() {
	fmt.Println("Starting GarpScan Engine Example...")

	// 1. Create a new engine with 5 workers
	engine := garpscan.NewEngine(garpscan.WithWorkers(5))

	// 2. Add an example plugin
	engine.AddPlugin(&ExamplePlugin{})

	// 3. Define targets
	targets := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	// 4. Start the scan with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultsChan := engine.Scan(ctx, targets)

	// 5. Read results as they stream in
	for res := range resultsChan {
		fmt.Printf("[RESULT] Target: %s | Plugin: %s | Status: %s | Data: %v\n", 
			res.Target, res.PluginName, res.Status, res.Data)
	}

	fmt.Println("Scan completed!")
}
