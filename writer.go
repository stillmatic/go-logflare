package gologflare

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// StdWriter is a writer for the standard library logger
type StdWriter struct {
	*LogflareClient
}

func NewStdWriter(client *LogflareClient) *StdWriter {
	return &StdWriter{client}
}

func (hw *StdWriter) Write(p []byte) (n int, err error) {
	log := LogData{Message: string(p)}
	hw.AddLog(log)
	return len(p), nil
}

// MultiWriter is a helper for writing to multiple writers
type MultiWriter struct {
	writers []io.Writer
}

func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, writer := range mw.writers {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}
	return len(p), nil
}

type SlogData struct {
	Level    string                 `json:"level"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"-"`
}

func (d *SlogData) GetLevel() string {
	return d.Level
}

func (d *SlogData) GetMessage() string {
	return d.Message
}

func (d *SlogData) GetMetadata() map[string]interface{} {
	return d.Metadata
}

type ZerologData struct {
	Level    string                 `json:"level"`
	Message  string                 `json:"msg"`
	Metadata map[string]interface{} `json:"-"`
}

func (d *ZerologData) GetLevel() string {
	return d.Level
}

func (d *ZerologData) GetMessage() string {
	return d.Message
}

func (d *ZerologData) GetMetadata() map[string]interface{} {
	return d.Metadata
}

type AuxData interface {
	GetLevel() string
	GetMessage() string
	GetMetadata() map[string]interface{}
}

func Convert[T AuxData](p []byte, levelKey, messageKey string) (logData *LogData, err error) {
	var t T
	if err := json.Unmarshal(p, &t); err != nil {
		return nil, err
	}

	sb := &strings.Builder{}
	sb.WriteString(strings.ToUpper(t.GetLevel()))
	sb.WriteString(": ")
	sb.WriteString(t.GetMessage())
	logData = &LogData{
		Message:  sb.String(),
		Metadata: t.GetMetadata(),
	}
	return logData, nil
}

// ConvertExplicit converts a json log line to a LogData struct.
// It is provided for completeness but is not recommended for use - the
// Convert function should be used instead, as it is about 2x as fast.
func ConvertExplicit(p []byte, levelKey, messageKey string) (*LogData, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal(p, &metadata); err != nil {
		return nil, err
	}

	sb := &strings.Builder{}
	if level, ok := metadata[levelKey]; ok {
		levelStr, ok := level.(string)
		if !ok {
			return nil, fmt.Errorf("expected level to be string, got: %v", reflect.TypeOf(level))
		}
		sb.WriteString(strings.ToUpper(levelStr))
		sb.WriteString(": ")
		delete(metadata, levelKey)
	}

	if msg, ok := metadata[messageKey]; ok {
		msgStr, ok := msg.(string)
		if !ok {
			return nil, fmt.Errorf("expected message to be string, got: %v", reflect.TypeOf(msg))
		}
		sb.WriteString(msgStr)
		delete(metadata, messageKey)
	}
	logData := &LogData{
		Message:  sb.String(),
		Metadata: metadata,
	}
	return logData, nil
}
