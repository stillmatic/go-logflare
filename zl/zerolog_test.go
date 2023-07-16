package zl_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	gologflare "github.com/stillmatic/go-logflare"
	"github.com/stillmatic/go-logflare/zl"
)

func TestZerolog(t *testing.T) {
	API_KEY := os.Getenv("LOGFLARE_API_KEY")
	sourceID := "236a182e-4ffb-4b72-b44a-7c21b3291a8f"
	client := gologflare.NewLogflareClient(API_KEY, &sourceID, nil, nil, 1*time.Second, 0)
	defer func() {
		client.Flush()
	}()
	t.Run("io.Writer", func(t *testing.T) {
		zlw := zl.NewZerologWriter(client)
		sw := zerolog.SyncWriter(zlw)
		mw := gologflare.NewMultiWriter(sw, os.Stdout)
		logger := zerolog.New(mw).With().Timestamp().Logger()
		for i := 0; i < 5; i++ {
			logger.Info().Int("count", i).Msg("hello world")
			time.Sleep(time.Millisecond * 500) // just for testing
		}
	})
	t.Run("zerolog.Hook", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Logger()

		logger = logger.Hook(&zl.LogflareHook{client})
		for i := 0; i < 5; i++ {
			logger.Info().Int("count", i).Msg("hook world")
			time.Sleep(time.Millisecond * 500) // just for testing
		}
	})
}
