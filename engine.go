package garpscan

import (
	"context"
	"sync"
	"time"
)

// Engine is the core structure that orchestrates plugins and targets.
type Engine struct {
	workers int
	plugins []Plugin
}

// Option configures the Engine.
type Option func(*Engine)

// WithWorkers sets the maximum number of concurrent workers.
func WithWorkers(n int) Option {
	return func(e *Engine) {
		if n > 0 {
			e.workers = n
		}
	}
}

// NewEngine creates a new scanning engine with the provided options.
func NewEngine(options ...Option) *Engine {
	e := &Engine{
		workers: 10, // default workers
		plugins: make([]Plugin, 0),
	}

	for _, opt := range options {
		opt(e)
	}

	return e
}

// AddPlugin registers a plugin with the engine.
func (e *Engine) AddPlugin(p Plugin) {
	e.plugins = append(e.plugins, p)
}

// scanTask defines an internal unit of work.
type scanTask struct {
	plugin Plugin
	target string
}

// Scan begins the scanning process for all registered plugins against all targets.
// It returns a read-only channel where results will be streamed as they complete.
// The channel will be closed automatically when all scanning operations are done.
func (e *Engine) Scan(ctx context.Context, targets []string) <-chan *Result {
	resultsChan := make(chan *Result, e.workers*2)
	tasksChan := make(chan scanTask, e.workers*2)

	var wg sync.WaitGroup

	// 1. Spawn workers
	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done(): // Cancelled or timed out
					return
				case task, ok := <-tasksChan:
					if !ok {
						return // tasks channel closed, we are done
					}

					// Execute the plugin run
					res, err := task.plugin.Run(ctx, task.target)
					if err != nil {
						// Optionally we could wrap the error in a generic result struct
						res = &Result{
							Target:     task.target,
							PluginName: task.plugin.Name(),
							Status:     "error",
							Data:       err.Error(),
							Timestamp:  time.Now(),
						}
					}
					
					// Ignore nil results (if plugin chooses to return nothing)
					if res != nil {
						select {
						case resultsChan <- res:
						case <-ctx.Done():
							return // abort writing if context is cancelled
						}
					}
				}
			}
		}()
	}

	// 2. Feed the tasks pipeline, wait, and auto-close channels
	go func() {
		// Feed tasks
		for _, p := range e.plugins {
			for _, t := range targets {
				select {
				case <-ctx.Done(): // Abort feeding early
					break // Instead of break, let's just use a flag or label
				case tasksChan <- scanTask{plugin: p, target: t}:
				}
			}
		}
		
		// Close tasks channel to signal workers there's no more work
		close(tasksChan)
		
		// Wait for all workers to finish their current tasks
		wg.Wait()
		
		// Everything is done, close the overall results pipeline
		close(resultsChan)
	}()

	return resultsChan
}
