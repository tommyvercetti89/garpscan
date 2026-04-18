package garpscan

import (
	"context"
	"testing"
	"time"
)

// mockPlugin is a test plugin that returns a predictable result.
type mockPlugin struct {
	name      string
	shouldErr bool
}

func (m *mockPlugin) Name() string {
	return m.name
}

func (m *mockPlugin) Run(ctx context.Context, target string) (*Result, error) {
	time.Sleep(10 * time.Millisecond) // Simulate work
	
	if m.shouldErr {
		return nil, context.DeadlineExceeded // arbitrary error for testing
	}
	
	return &Result{
		Target:     target,
		PluginName: m.name,
		Status:     "success",
		Data:       "test-data",
		Timestamp:  time.Now(),
	}, nil
}

func TestEngine_Scan(t *testing.T) {
	engine := NewEngine(WithWorkers(2))
	
	engine.AddPlugin(&mockPlugin{name: "mock-1", shouldErr: false})
	
	targets := []string{"target-1", "target-2"}
	ctx := context.Background()
	
	resultsChan := engine.Scan(ctx, targets)
	
	var results []*Result
	for res := range resultsChan {
		results = append(results, res)
	}
	
	expectedResults := len(targets) * 1 // 1 plugin * 2 targets
	if len(results) != expectedResults {
		t.Errorf("Expected %d results, got %d", expectedResults, len(results))
	}
}

func TestEngine_Scan_Cancellation(t *testing.T) {
	engine := NewEngine(WithWorkers(2))
	engine.AddPlugin(&mockPlugin{name: "slow-mock", shouldErr: false})
	
	targets := []string{"target-1", "target-2", "target-3", "target-4"}
	
	// Fast timeout to trigger cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	
	resultsChan := engine.Scan(ctx, targets)
	
	var results []*Result
	for res := range resultsChan {
		results = append(results, res)
	}
	
	// Some tasks might get cancelled, so we should expect fewer results than targets
	if len(results) == len(targets) {
		t.Errorf("Expected some tasks to be cancelled, but all completed")
	}
}
