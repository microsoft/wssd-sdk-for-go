// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package exporters

import (
	"bytes"
	"fmt"
	"go.opencensus.io/trace"
	"strconv"
)

// ConsoleExporter is a trace exporter and conforms to the inferface trace.Exporter. Prints tracing data to logs.
type ConsoleExporter struct {
	Use        bool
	ParentMap  map[string][]string
	LogChannel chan *spanWriteDataConsole
}

// SpanWriteData keeps information about spans which should be written to a file.
type spanWriteDataConsole struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	StartTime   string `json:"starttime"`
	EndTime     string `json:"endtime"`
	Duration    string `json:"duration"`
	ParentID    trace.SpanID
	Annotations []trace.Annotation
	TraceID     int32
}

// NewConsoleExporter does
func NewConsoleExporter() *ConsoleExporter {
	ce := ConsoleExporter{
		true,
		make(map[string][]string),
		make(chan *spanWriteDataConsole),
	}

	go ce.monitorLogs()

	return &ce
}

func (ce *ConsoleExporter) monitorLogs() {

	ceChannel := ce.LogChannel
	for {
		nextSpanMessage, ok := <-ceChannel
		if !ok {
			// Channel has been closed
			break
		}

		var buffer bytes.Buffer

		buffer.WriteString(fmt.Sprintf("[SPAN_START: %+v %v] Name: %s  SpanID: %x  ParentSpanID: %x  Duration: %s\n",
			nextSpanMessage.StartTime,
			nextSpanMessage.TraceID,
			nextSpanMessage.Name,
			nextSpanMessage.ID,
			nextSpanMessage.ParentID.String(),
			nextSpanMessage.Duration))

		for _, anno := range nextSpanMessage.Annotations {

			if anno.Attributes["CallerLocation"] == nil {
				return
			}

			callerLocation, ok := anno.Attributes["CallerLocation"].(string)
			if !ok {
				return
			}

			buffer.WriteString(
				fmt.Sprintf(
					"[LOG: %+v %s %v] %s\n", anno.Time.Format("01-02 15:04:0000005"),
					callerLocation,
					nextSpanMessage.TraceID,
					anno.Message))
		}

		if ce.ParentMap[nextSpanMessage.ID] != nil {
			for _, span := range ce.ParentMap[nextSpanMessage.ID] {
				buffer.WriteString(span)
			}
		}

		buffer.WriteString(fmt.Sprintf("[SPAN_END: %+v %v] Name: %s  Duration: %s\n",
			nextSpanMessage.EndTime,
			nextSpanMessage.TraceID,
			nextSpanMessage.Name,
			nextSpanMessage.Duration))

		if nextSpanMessage.ParentID == [8]byte{} {
			fmt.Printf(buffer.String())
		} else if ce.ParentMap[nextSpanMessage.ParentID.String()] != nil {
			ce.ParentMap[nextSpanMessage.ParentID.String()] = append(ce.ParentMap[nextSpanMessage.ParentID.String()], buffer.String())
		} else {
			ce.ParentMap[nextSpanMessage.ParentID.String()] = []string{buffer.String()}
		}
	}
}

// ExportSpan prints information about given span to logs.
func (ce *ConsoleExporter) ExportSpan(sd *trace.SpanData) {

	traceIdentiferID := GetHashedTraceID(sd.SpanContext.TraceID)
	ce.LogChannel <- &spanWriteDataConsole{
		sd.Name,
		sd.SpanID.String(),
		sd.StartTime.Format("01-02 15:04:0000005"),
		sd.EndTime.Format("01-02 15:04:0000005"),
		strconv.FormatFloat(sd.EndTime.Sub(sd.StartTime).Seconds(), 'f', 2, 64) + "s",
		sd.ParentSpanID,
		sd.Annotations,
		traceIdentiferID,
	}
}
