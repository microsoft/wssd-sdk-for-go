// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package exporters

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"go.opencensus.io/trace"
)

// FileExporter is a trace exporter and conforms to the inferface trace.Exporter. Prints tracing data to logs.
type FileExporter struct {
	Use      bool
	Filepath string
	Spans    []*SpanWriteData
	// Writer   *bufio.Writer
}

// SpanWriteData keeps information about spans which should be written to a file.
type SpanWriteData struct {
	Name     string
	ID       string
	ParentID string
	Duration string
}

// NewFileExporter does
func NewFileExporter(filepath string) *FileExporter {
	return &FileExporter{filepath != "", filepath, make([]*SpanWriteData, 0)}
}

// End writes JSON object to a file and closes it.
func (cse *FileExporter) End() error {
	if !cse.Use {
		return nil
	}
	file, err := os.Create(cse.Filepath)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(cse.Spans)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		return err
	}
	writer.Flush()
	file.Close()
	return nil
}

// ExportSpan prints information about given span to logs.
func (cse *FileExporter) ExportSpan(sd *trace.SpanData) {
	if !cse.Use {
		return
	}
	cse.Spans = append(cse.Spans, &SpanWriteData{sd.Name, sd.SpanID.String(), sd.ParentSpanID.String(),
		strconv.FormatFloat(sd.EndTime.Sub(sd.StartTime).Seconds(), 'f', 2, 64) + "s"})
}
