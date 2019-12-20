// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package trace

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	opentrace "go.opencensus.io/trace"

	"github.com/microsoft/wssdagent/pkg/exporters"
)

// Settings for the tracing.
type Settings struct {
	Verbose            bool
	JaegerAgentURI     string
	JaegerCollectorURI string
}

var fileExporter *exporters.FileExporter
var consoleExporter *exporters.ConsoleExporter

// Configure initializes tracing
func Configure(settings *Settings) error {
	opentrace.ApplyConfig(opentrace.Config{DefaultSampler: opentrace.AlwaysSample()})
	// opentrace.RegisterExporter(&exporters.ZipkinExporter{Verbose: settings.Verbose})
	fileExporter = exporters.NewFileExporter()
	opentrace.RegisterExporter(fileExporter)
	consoleExporter = exporters.NewConsoleExporter()
	opentrace.RegisterExporter(consoleExporter)
	opentrace.ApplyConfig(opentrace.Config{DefaultSampler: opentrace.AlwaysSample()})
	return registerJaeger(settings)
}

// InsertAnnotations inserts serialized span context and update uuid to the annotations map.
func InsertAnnotations(annotations map[string]string, span *LogSpan, uuid string) map[string]string {
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["trace.wssdagent.microsoft.com/spancontext"] = span.SerializeSpanContext()
	annotations["trace.wssdagent.microsoft.com/uuid"] = uuid
	return annotations
}

// GetSpanContext returns the string representing the SpanContext kept in annotations.
func GetSpanContext(annotations map[string]string) string {
	if annotations == nil {
		annotations = map[string]string{}
	}
	return annotations["trace.wssdagent.microsoft.com/spancontext"]
}

// GetUpdateUUID returns the string representing the update UUID.
func GetUpdateUUID(annotations map[string]string) string {
	if annotations == nil {
		annotations = map[string]string{}
	}
	return annotations["trace.wssdagent.microsoft.com/uuid"]
}

// CopyTracingAnnotations copies the tracing data from the source annotations map to the destination annotations map.
func CopyTracingAnnotations(src, dst map[string]string) {
	dst["trace.wssdagent.microsoft.com/spancontext"] = src["trace.wssdagent.microsoft.com/spancontext"]
	dst["trace.wssdagent.microsoft.com/uuid"] = src["trace.wssdagent.microsoft.com/uuid"]
}

func registerJaeger(settings *Settings) error {
	if settings.JaegerAgentURI == "" && settings.JaegerCollectorURI == "" {
		return nil
	}

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     settings.JaegerAgentURI,
		CollectorEndpoint: settings.JaegerCollectorURI,
		ServiceName:       "trace",
	})
	if err != nil {
		return err
	}

	opentrace.RegisterExporter(je)
	return nil
}
