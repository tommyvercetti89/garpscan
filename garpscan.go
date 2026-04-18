package garpscan

import (
	"context"
	"time"
)

// Result represents the outcome of a single scan task by a plugin.
type Result struct {
	Target     string      `json:"target"`
	PluginName string      `json:"plugin_name"`
	Status     string      `json:"status"` // e.g., "open", "closed", "error"
	Data       interface{} `json:"data"`   // Any specific data returned by the plugin
	Timestamp  time.Time   `json:"timestamp"`
}

// Plugin is the interface that all scanner modules must implement.
type Plugin interface {
	// Name returns the unique name of the plugin (e.g., "portscanner").
	Name() string
	// Run executes the scanning logic against a single target.
	// It should respect context cancellation.
	Run(ctx context.Context, target string) (*Result, error)
}
