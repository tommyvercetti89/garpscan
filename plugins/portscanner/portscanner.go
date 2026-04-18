package portscanner

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tommyvercetti89/garpscan"
)

// PortScanner checks standard TCP ports for a given target.
type PortScanner struct {
	ports   []int
	timeout time.Duration
}

// New creates a new PortScanner instance.
func New(ports []int, timeout time.Duration) *PortScanner {
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	return &PortScanner{
		ports:   ports,
		timeout: timeout,
	}
}

func (p *PortScanner) Name() string {
	return "portscanner"
}

// Run checks the target ports and returns an array of open ports in the Data field.
func (p *PortScanner) Run(ctx context.Context, target string) (*garpscan.Result, error) {
	var openPorts []int

	for _, port := range p.ports {
		// Respect context cancellation before each port
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		address := fmt.Sprintf("%s:%d", target, port)
		
		// Use a custom dialer with context
		dialer := net.Dialer{Timeout: p.timeout}
		conn, err := dialer.DialContext(ctx, "tcp", address)
		
		if err == nil {
			openPorts = append(openPorts, port)
			conn.Close()
		}
	}

	status := "completed"
	if len(openPorts) == 0 {
		status = "no_open_ports"
	} else {
		status = "ports_found"
	}

	return &garpscan.Result{
		Target:     target,
		PluginName: p.Name(),
		Status:     status,
		Data:       openPorts,
		Timestamp:  time.Now(),
	}, nil
}
