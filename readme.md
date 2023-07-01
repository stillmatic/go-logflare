# go-logflare

`go-logflare` implements transport for various Golang logging libraries to send logs to [Logflare](https://logflare.app). The underlying transport contains a configurable buffer and flush frequency to reduce the number of requests sent.

## Usage

Currently, we support the following logging libraries:

- [stdlib `log`](https://pkg.go.dev/log)
- [stdlib `slog`](https://pkg.go.dev/golang.org/x/exp/slog)
- [zerolog](https://github.com/rs/zerolog)

The standardlib logger is provided in the main package. `slog` and `zerolog` are provided in their respective subpackages, to isolate dependencies. The main package is dependency-free. `slog` will eventually be promoted to the main package, once it is promoted to the standard library.

Note that both `slog` and `zerolog` require overhead for serdes - internally, we unmarshal the structs in order to send them to Logflare. The standardlib logger does not require this and should be the fastest.