// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package exporters

import (
	"go.opencensus.io/trace"
	"strconv"

	"contrib.go.opencensus.io/exporter/zipkin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"k8s.io/klog"
)

// spanWriteData keeps information about spans which should be written to a file.
type spanWriteDataZipkin struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	ParentID string `json:"parentid"`
	Duration string `json:"duration"`
}

// ZipkinExporter is a trace exporter and conforms to the inferface trace.Exporter. Prints tracing data to logs.
type ZipkinExporter struct {
	Verbose bool
	Zipkin  trace.Exporter
	Spans   []*spanWriteDataZipkin
}

// NewZipkinExporter does .. well not sure yet
func NewZipkinExporter(filepath string) *ZipkinExporter {
	localEndpoint, err := openzipkin.NewEndpoint("wssdagent", "0.0.0.0:5455")
	if err != nil {
		klog.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	ze := zipkin.NewExporter(reporter, localEndpoint)

	return &ZipkinExporter{Zipkin: ze}
}

// ExportSpan prints information about given span to logs.
func (cse *ZipkinExporter) ExportSpan(sd *trace.SpanData) {
	cse.Spans = append(cse.Spans, &spanWriteDataZipkin{sd.Name, sd.SpanID.String(), sd.ParentSpanID.String(),
		strconv.FormatFloat(sd.EndTime.Sub(sd.StartTime).Seconds(), 'f', 2, 64) + "s"})
}
