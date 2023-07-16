package discord_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	gologflare "github.com/stillmatic/go-logflare"
	"github.com/stillmatic/go-logflare/discord"
	"github.com/stillmatic/go-logflare/zl"
)

func TestZerolog(t *testing.T) {
	API_KEY := os.Getenv("DISCORD_WEBHOOK_KEY")
	url := "https://discord.com/api/webhooks/" + API_KEY
	client := discord.NewDiscordClient(url, "test", 1*time.Second, 5, nil)
	defer func() {
		client.Flush()
	}()
	zlw := zl.NewZerologWriter(client)
	sw := zerolog.SyncWriter(zlw)
	mw := gologflare.NewMultiWriter(sw, os.Stdout)
	logger := zerolog.New(mw).With().Timestamp().Logger()
	for i := 0; i < 5; i++ {
		logger.Info().Int("count", i).Msg("hello world")
		time.Sleep(time.Millisecond * 500) // just for testing
	}
}
