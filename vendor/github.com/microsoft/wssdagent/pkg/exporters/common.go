package exporters

import (
	"go.opencensus.io/trace"
)

// TODO Decide on a common tracing data

// SpanWriteData keeps information about spans which should be written to a file.
type SpanWriteData struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	ParentID string `json:"parentid"`
	Duration string `json:"duration"`
}

// GetHashedTraceID gets a human readable identifer from the trace guid.
func GetHashedTraceID(traceID trace.TraceID) int32 {
	return int32(traceID[15])
}
