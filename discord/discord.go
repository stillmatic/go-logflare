package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	gologflare "github.com/stillmatic/go-logflare"
	"golang.org/x/sync/errgroup"
)

type DiscordClient struct {
	url         string
	name        string
	client      *http.Client
	buffer      []gologflare.LogData
	mu          sync.Mutex
	flushPeriod time.Duration
	bufPool     sync.Pool
	flushSize   int
}

func NewDiscordClient(url, name string, flushPeriod time.Duration, flushSize int, httpClient *http.Client) *DiscordClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	dc := &DiscordClient{
		url:         url,
		name:        name,
		client:      httpClient,
		buffer:      make([]gologflare.LogData, 0),
		flushPeriod: flushPeriod,
		flushSize:   flushSize,
		bufPool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
	go dc.startTimer()
	return dc
}

// startTimer starts a timer that flushes the log buffer every tick
func (c *DiscordClient) startTimer() {
	ticker := time.NewTicker(c.flushPeriod)
	for range ticker.C {
		err := c.Flush()
		if err != nil {
			fmt.Printf("Error flushing logs: %s\n", err)
		}
	}
}

func (c *DiscordClient) AddLog(log gologflare.LogData) {
	c.mu.Lock()
	c.buffer = append(c.buffer, log)
	c.mu.Unlock()
	if len(c.buffer) >= c.flushSize {
		c.Flush()
	}
}

type Message struct {
	Username  *string  `json:"username,omitempty"`
	AvatarUrl *string  `json:"avatar_url,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Embeds    *[]Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       *string    `json:"title,omitempty"`
	Url         *string    `json:"url,omitempty"`
	Description *string    `json:"description,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Fields      *[]Field   `json:"fields,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
}

type Field struct {
	Name   *string `json:"name,omitempty"`
	Value  *string `json:"value,omitempty"`
	Inline *bool   `json:"inline,omitempty"`
}

func strPtr(s string) *string {
	return &s
}

func (c *DiscordClient) convertMessageToDiscord(msg gologflare.LogData) Message {
	var true = true
	var message Message
	var embed Embed
	var fields []Field
	content := msg.Message
	embedTitle := "metadata"

	message.Username = &c.name
	if msg.Level != "" {
		embedTitle = msg.Level
		// strip from content
		splits := strings.SplitN(content, ": ", 2)
		if len(splits) > 1 {
			content = splits[1]
		}
	}
	message.Content = &content
	for key, value := range msg.Metadata {
		key := key
		value := value
		v := fmt.Sprintf("%v", value)
		switch key {
		case "time", "timestamp":
			ts, _ := time.Parse(time.RFC3339, v)
			embed.Timestamp = &ts
		default:
			fields = append(fields, Field{Name: &key, Value: &v, Inline: &true})
		}
	}
	embed.Title = &embedTitle
	embed.Fields = &fields
	message.Embeds = &[]Embed{embed}
	return message
}

func (c *DiscordClient) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.buffer) == 0 {
		return nil
	}

	var eg errgroup.Group
	eg.SetLimit(4)
	for _, log := range c.buffer {
		log := log
		eg.Go(func() error {
			message := c.convertMessageToDiscord(log)
			buf := c.bufPool.Get().(*bytes.Buffer)
			if buf == nil {
				buf = new(bytes.Buffer)
			}
			err := json.NewEncoder(buf).Encode(message)
			if err != nil {
				return fmt.Errorf("error marshalling logs: %s", err)
			}
			resp, err := http.Post(c.url, "application/json", buf)
			if err != nil {
				return fmt.Errorf("error posting logs: %s", err)
			}
			if resp.StatusCode != 200 && resp.StatusCode != 204 {
				defer resp.Body.Close()

				responseBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				return fmt.Errorf(string(responseBody))
			}
			c.bufPool.Put(buf)
			return nil
		})
	}
	c.buffer = make([]gologflare.LogData, 0)
	err := eg.Wait()
	if err != nil {
		return err
	}
	return nil
}
