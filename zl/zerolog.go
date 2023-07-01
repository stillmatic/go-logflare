package zl

import (
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

// func (hw *ZerologWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
// 	return hw.Write(p)
// }
