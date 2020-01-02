// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package trace

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"

	opentrace "go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

// LogSpan wraps a trace.Span type and adds logging when starting/ending span.
type LogSpan struct {
	span *opentrace.Span
	name string
	err  error
}

// StringAttribute is equivalent to opentrace.StringAttribute.
func StringAttribute(key, value string) opentrace.Attribute {
	return opentrace.StringAttribute(key, value)
}

// NewSpan creates a new Span and logs information about it.
func NewSpan(ctx context.Context, args ...string) (context.Context, *LogSpan) {
	name := strings.Join(args, " ")
	ctx, span := opentrace.StartSpan(ctx, name)
	return ctx, &LogSpan{span, name, nil}
}

// MarkAsEntitySpan creates a new Span and logs information about it.
func (s *LogSpan) MarkAsEntitySpan(entityPath string, entity string) {
	s.span.AddAttributes(StringAttribute("EntityPath", entityPath))
	s.span.AddAttributes(StringAttribute("Entity", entity))
}

// Log creates a new Span and logs information about it.
func (s *LogSpan) Log(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	callerLocation := fileName + ":" + strconv.Itoa(line) + " "

	log := fmt.Sprintf(format, a...)
	s.span.Annotate([]opentrace.Attribute{StringAttribute("CallerLocation", callerLocation)}, log)
}

// NewSpanWithRemoteParent starts a new child span of the span from the given parent.
func NewSpanWithRemoteParent(ctx context.Context, parentSerialized string, args ...string) (context.Context, *LogSpan) {
	sc, _ := DeserializeSpanContext(parentSerialized)
	name := strings.Join(args, " ")
	ctx, span := opentrace.StartSpanWithRemoteParent(ctx, name, sc)
	return ctx, &LogSpan{span, name, nil}
}

// SetError sets err attribute and returns given argument.
func (s *LogSpan) SetError(err error) error {
	s.err = err
	return err
}

// EndSpan calls End() function on span and logs information about it. Sets status in case of error.
func (s *LogSpan) EndSpan() {
	s.End(s.err)
}

// End calls End() function on span and logs information about it. Sets status in case of error.
func (s *LogSpan) End(err error) {
	if err != nil {
		s.span.SetStatus(opentrace.Status{
			Code:    opentrace.StatusCodeUnknown,
			Message: err.Error(),
		})
	}
	s.span.End()
}

// Annotate runs annotate function on internal span.
func (s *LogSpan) Annotate(attributes []opentrace.Attribute, str string) {
	s.span.Annotate(attributes, str)
}

// SerializeSpanContext returns a string which is a text representation of the SpanContext
func (s *LogSpan) SerializeSpanContext() string {
	return b64.StdEncoding.EncodeToString(propagation.Binary(s.span.SpanContext()))
}

// DeserializeSpanContext returns a SpanContext for given text representation.
func DeserializeSpanContext(ssc string) (sc opentrace.SpanContext, ok bool) {
	arr, _ := b64.StdEncoding.DecodeString(ssc)
	return propagation.FromBinary(arr)
}
