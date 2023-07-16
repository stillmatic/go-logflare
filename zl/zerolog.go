package zl

import (
	"github.com/rs/zerolog"
	gologflare "github.com/stillmatic/go-logflare"
)

const (
	levelKey   = "level"
	messageKey = "message"
)

type ZerologWriter struct {
	*gologflare.LogflareClient
}

func NewZerologWriter(client *gologflare.LogflareClient) *ZerologWriter {
	return &ZerologWriter{client}
}

func (hw *ZerologWriter) Write(p []byte) (n int, err error) {
	logData, err := gologflare.Convert(p, levelKey, messageKey)
	if err != nil {
		return 0, err
	}

	hw.AddLog(*logData)
	return len(p), nil
}

type LogflareHook struct {
	*gologflare.LogflareClient
}

func (t *LogflareHook) Run(
	e *zerolog.Event,
	level zerolog.Level,
	message string,
) {
	logData := gologflare.LogData{
		Message:  message,
		Metadata: make(map[string]interface{}),
	}
	t.AddLog(logData)
}

// func (hw *ZerologWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
// 	return hw.Write(p)
// }
