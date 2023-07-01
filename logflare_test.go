package gologflare_test

import (
	"testing"

	gologflare "github.com/stillmatic/go-logflare"
)

func BenchmarkConvert(b *testing.B) {
	b.Run("explicit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := gologflare.ConvertExplicit(
				[]byte(`{"level":"info","msg":"hello slog","count":1,"key":"value"}`),
				"level",
				"msg",
			)
			if err != nil {
				b.Fatal(err)
			}
			if res == nil {
				b.Fatal("expected log data to be non-nil")
			}
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := gologflare.Convert[*gologflare.SlogData](
				[]byte(`{"level":"info","message":"hello slog","count":1,"key":"value"}`),
				"level",
				"message",
			)
			if err != nil {
				b.Fatal(err)
			}
			if res == nil {
				b.Fatal("expected log data to be non-nil")
			}
		}
	})

	b.Run("generic z", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := gologflare.Convert[*gologflare.ZerologData](
				[]byte(`{"level":"info","msg":"hello slog","count":1,"key":"value"}`),
				"level",
				"msg",
			)
			if err != nil {
				b.Fatal(err)
			}
			if res == nil {
				b.Fatal("expected log data to be non-nil")
			}
		}
	})
}

func TestConvert(t *testing.T) {
	t.Run("explicit", func(t *testing.T) {
		res, err := gologflare.ConvertExplicit(
			[]byte(`{"level":"info","message":"hello slog","count":1,"key":"value"}`),
			"level",
			"message",
		)
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected log data to be non-nil")
		}
		if res.Message != "INFO: hello slog" {
			t.Fatalf("expected message to be 'INFO: hello slog', got '%s'", res.Message)
		}
	})

	t.Run("generic", func(t *testing.T) {
		res, err := gologflare.Convert[*gologflare.SlogData](
			[]byte(`{"level":"info","message":"hello slog","count":1,"key":"value"}`),
			"level",
			"message",
		)
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected log data to be non-nil")
		}
		if res.Message != "INFO: hello slog" {
			t.Fatalf("expected message to be 'INFO: hello slog', got '%s'", res.Message)
		}
	})

	t.Run("generic", func(t *testing.T) {
		res, err := gologflare.Convert[*gologflare.ZerologData](
			[]byte(`{"level":"info","msg":"hello slog","count":1,"key":"value"}`),
			"level",
			"msg",
		)
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected log data to be non-nil")
		}
		if res.Message != "INFO: hello slog" {
			t.Fatalf("expected message to be 'INFO: hello slog', got '%s'", res.Message)
		}
	})
}
