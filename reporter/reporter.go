package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tommyvercetti89/garpscan"
)

// ExportJSON writes a stream of results to the provided io.Writer in JSON lines format.
func ExportJSON(writer io.Writer, results <-chan *garpscan.Result) error {
	encoder := json.NewEncoder(writer)
	for res := range results {
		if err := encoder.Encode(res); err != nil {
			return fmt.Errorf("failed to encode json: %w", err)
		}
	}
	return nil
}

// ExportCSV writes results to the provided io.Writer in CSV format.
func ExportCSV(writer io.Writer, results <-chan *garpscan.Result) error {
	// Write header
	if _, err := fmt.Fprintln(writer, "Target,PluginName,Status,Data,Timestamp"); err != nil {
		return err
	}

	for res := range results {
		dataStr := fmt.Sprintf("%v", res.Data)
		row := fmt.Sprintf("%s,%s,%s,\"%s\",%s",
			res.Target,
			res.PluginName,
			res.Status,
			dataStr,
			res.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		)
		if _, err := fmt.Fprintln(writer, row); err != nil {
			return err
		}
	}
	return nil
}
