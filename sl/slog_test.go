package sl_test

import (
	"os"
	"testing"
	"time"

	gologflare "github.com/stillmatic/go-logflare"
	"github.com/stillmatic/go-logflare/sl"
	"golang.org/x/exp/slog"
)

func TestSlog(t *testing.T) {
	API_KEY := os.Getenv("LOGFLARE_API_KEY")
	sourceID := "236a182e-4ffb-4b72-b44a-7c21b3291a8f"
	client := gologflare.NewLogflareClient(API_KEY, &sourceID, nil, nil, 1*time.Second, 2)
	defer func() {
		client.Flush()
	}()
	sw := sl.NewSlogWriter(client)
	mw := gologflare.NewMultiWriter(sw, os.Stdout)
	logger := slog.New(slog.NewJSONHandler(mw, nil))
	for i := 0; i < 5; i++ {
		logger.Info("hello slog", "count", i, "logger", "slog")
		time.Sleep(time.Millisecond * 500) // just for testing
	}
}
