package gologflare

type LogClient interface {
	AddLog(LogData)
	Flush() error
}
