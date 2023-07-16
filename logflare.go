package gologflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	BASE_URL = "https://api.logflare.app/logs/"
)

type LogData struct {
	Message string `json:"message"`
	// Level is not actually in the spec and is dropped, but we need elsewhere.
	Level    string                 `json:"level,omitempty"`
	Metadata map[string]interface{} `json:"metadata"`
}

type LogPayload struct {
	Batch []LogData `json:"batch"`
}

type LogflareClient struct {
	url         string
	apiKey      string
	client      *http.Client
	buffer      []LogData
	mu          sync.Mutex
	flushPeriod time.Duration
	flushSize   int
}

type LogFlarer interface {
	AddLog(log LogData)
	Flush() error
}

func NewLogflareClient(
	apiKey string,
	sourceID *string,
	sourceName *string,
	httpClient *http.Client,
	flushPeriod time.Duration,
	flushSize int,
) *LogflareClient {
	baseURL, err := url.Parse(BASE_URL)
	if err != nil {
		panic(err)
	}
	q := baseURL.Query()
	if sourceID != nil {
		q.Set("source", *sourceID)
	} else if sourceName != nil {
		q.Set("source_name", *sourceName)
	} else {
		panic("Must specify either sourceID or sourceName")
	}
	baseURL.RawQuery = q.Encode()
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	lc := &LogflareClient{
		url:         baseURL.String(),
		apiKey:      apiKey,
		client:      httpClient,
		buffer:      make([]LogData, 0),
		flushPeriod: flushPeriod,
		flushSize:   flushSize,
	}
	go lc.startTimer()
	return lc
}

// startTimer starts a timer that flushes the log buffer every tick
func (c *LogflareClient) startTimer() {
	ticker := time.NewTicker(c.flushPeriod)
	for range ticker.C {
		err := c.Flush()
		if err != nil {
			fmt.Printf("Error flushing logs: %s\n", err)
		}
	}
}

func (c *LogflareClient) AddLog(log LogData) {
	c.mu.Lock()
	c.buffer = append(c.buffer, log)
	c.mu.Unlock()
	if len(c.buffer) >= c.flushSize {
		c.Flush()
	}
}

func (c *LogflareClient) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.buffer) == 0 {
		return nil
	}
	payload, err := json.Marshal(LogPayload{Batch: c.buffer})
	if err != nil {
		return fmt.Errorf("error marshalling logs: %s", err)
	}
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error posting logs: %s", err)
	}
	defer resp.Body.Close()
	c.buffer = make([]LogData, 0)
	// TODO: this block should be deferred and sent somewhere else
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error posting logs: %s", string(respBody))
	}
	return nil
}
