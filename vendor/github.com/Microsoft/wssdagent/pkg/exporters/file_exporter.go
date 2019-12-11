// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package exporters

import (
	"encoding/json"
	"go.opencensus.io/trace"
	"io/ioutil"
	"k8s.io/klog"
	"os"
	"path"
	"strconv"
)

const (
	INITIAL_LOG_NAME   = "init-log"
	EVENT_LOGFILE_NAME = "event-log"
)

// FileExporter is a trace exporter and conforms to the inferface trace.Exporter. Prints tracing data to logs.
type FileExporter struct {
	Use         bool
	FilepathMap map[string]*string
	LogChannel  chan *spanWriteDataFile
}

// spanWriteDataFile keeps information about spans which should be written to a file.
type spanWriteDataFile struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	ParentID    string `json:"parentid"`
	Duration    string `json:"duration"`
	FilePath    string `json:"-"`
	Entity      string `json:"entity"`
	Annotations []trace.Annotation
}

// NewFileExporter does
func NewFileExporter() *FileExporter {
	fe := FileExporter{
		true,
		make(map[string]*string),
		make(chan *spanWriteDataFile),
	}

	go fe.monitorLogs()

	return &fe
}

func (fe *FileExporter) monitorLogs() {

	feChannel := fe.LogChannel
	for {
		nextSpanMessage, ok := <-feChannel
		if !ok {
			// Channel has been closed
			break
		}

		jsonBytes, err := json.Marshal(nextSpanMessage)
		if err != nil {
			klog.Infof("Json Marshaling Failed with %v. Log and Continue", err)
			continue
		}

		logFile, err := os.OpenFile(nextSpanMessage.FilePath, os.O_APPEND, 0644)
		if err != nil {
			klog.Infof("Open File Failed with %v. Log and Continue", err)
			continue
		}

		defer logFile.Close()

		_, err = logFile.Write([]byte(jsonBytes))
		if err != nil {
			klog.Infof("Writing File Failed with %v. Log and Continue", err)
			continue
		}
	}
}

// ExportSpan prints information about given span to logs.
func (fe *FileExporter) ExportSpan(sd *trace.SpanData) {

	if sd.Attributes["EntityPath"] == nil {
		return
	}

	entityPath, ok := sd.Attributes["EntityPath"].(string)
	if !ok {
		return
	}

	if sd.Attributes["Entity"] == nil {
		return
	}

	entity, ok := sd.Attributes["EntityPath"].(string)
	if !ok {
		return
	}

	filePath := path.Join(entityPath, EVENT_LOGFILE_NAME)

	fe.LogChannel <- &spanWriteDataFile{
		sd.Name,
		sd.SpanID.String(),
		sd.ParentSpanID.String(),
		strconv.FormatFloat(sd.EndTime.Sub(sd.StartTime).Seconds(), 'f', 2, 64) + "s",
		filePath,
		entity,
		sd.Annotations,
	}
}

// WriteInitialLog bootstraps the log process
func WriteInitialLog(basepath string) error {
	initLog := spanWriteDataFile{
		Name:     INITIAL_LOG_NAME,
		ID:       "",
		ParentID: "",
		Duration: "",
		FilePath: "",
	}

	out, err := json.Marshal(initLog)
	if err != nil {
		return nil
	}

	ioutil.WriteFile(
		path.Join(basepath, EVENT_LOGFILE_NAME),
		out,
		0644)

	return nil
}
