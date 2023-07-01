package sl

import (
	gologflare "github.com/stillmatic/go-logflare"
)

const (
	levelKey   = "level"
	messageKey = "msg"
)

// SlogWriter is a writer for the slog logger
type SlogWriter struct {
	*gologflare.LogflareClient
}

func NewSlogWriter(client *gologflare.LogflareClient) *SlogWriter {
	return &SlogWriter{client}
}

func (hw *SlogWriter) Write(p []byte) (n int, err error) {
	logData, err := gologflare.Convert(p, levelKey, messageKey)
	if err != nil {
		return 0, err
	}

	hw.AddLog(*logData)
	return len(p), nil
}
