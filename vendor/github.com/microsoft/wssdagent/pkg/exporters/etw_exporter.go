package exporters

import (
	"github.com/Microsoft/go-winio/pkg/etw"
	"go.opencensus.io/trace"
	"time"
)

const etwProviderName = "nodeAgent"

// trace.Exporter which allows for exporting to an etw provider
type EtwExporter struct {
	provider      *etw.Provider
	closeProvider bool
	LogChannel    chan *spanWriteDataEtw
}

type spanWriteDataEtw struct {
	Name        string    `json:"name"`
	ID          string    `json:"id"`
	StartTime   time.Time `json:"starttime"`
	EndTime     time.Time `json:"endtime"`
	ParentID    trace.SpanID
	Annotations []trace.Annotation
	TraceID     int32
}

func (sEtw *spanWriteDataEtw) createEtwFieldOpt(name string, id string, parentId string, time time.Time, message string) []etw.FieldOpt {
	return []etw.FieldOpt{
		etw.StringField("Name", name),
		etw.StringField("ActivityID", id),
		etw.StringField("ParentActivityID", parentId),
		etw.Time("Time", time),
		etw.StringField("Message", message),
	}
}

// Functon that takes the data from spanWriteDateEtw and converts it to a list of etwLogs to be written
func (sEtw *spanWriteDataEtw) createEtwEntryFields() [][]etw.FieldOpt {

	etwFieldsList := [][]etw.FieldOpt{}
	etwFieldsList = append(etwFieldsList, sEtw.createEtwFieldOpt(sEtw.Name, sEtw.ID, sEtw.ParentID.String(), sEtw.StartTime, "Starting Span"))
	for _, a := range sEtw.Annotations {
		etwFieldsList = append(etwFieldsList, sEtw.createEtwFieldOpt(sEtw.Name, sEtw.ID, sEtw.ParentID.String(), a.Time, a.Message))
	}
	etwFieldsList = append(etwFieldsList, sEtw.createEtwFieldOpt(sEtw.Name, sEtw.ID, sEtw.ParentID.String(), sEtw.EndTime, "Ending Span"))
	return etwFieldsList
}

func NewEtwExporter() *EtwExporter {
	// Provider ID: {6b6558ea-c87d-529a-d216-de7c900492dc}
	hook, _ := etw.NewProvider(etwProviderName, nil)
	ee := EtwExporter{
		provider:      hook,
		closeProvider: true,
		LogChannel:    make(chan *spanWriteDataEtw),
	}

	go ee.monitorLogs()

	return &ee
}

func (ee *EtwExporter) monitorLogs() {

	eeChannel := ee.LogChannel
	for {
		nextSpanMessage, ok := <-eeChannel
		if !ok {
			// Channel closed
			if ee.closeProvider {
				_ = ee.provider.Close()
			}
			break
		}

		etwFieldsList := nextSpanMessage.createEtwEntryFields()
		for _, etwEntryFields := range etwFieldsList {
			ee.provider.WriteEvent("TraceSpanEntry", []etw.EventOpt{}, etwEntryFields)
		}
	}
}

func (ee *EtwExporter) ExportSpan(sd *trace.SpanData) {
	traceIdentifierID := GetHashedTraceID(sd.SpanContext.TraceID)
	ee.LogChannel <- &spanWriteDataEtw{
		Name:        sd.Name,
		ID:          sd.SpanID.String(),
		StartTime:   sd.StartTime,
		EndTime:     sd.EndTime,
		ParentID:    sd.ParentSpanID,
		Annotations: sd.Annotations,
		TraceID:     traceIdentifierID,
	}
}
