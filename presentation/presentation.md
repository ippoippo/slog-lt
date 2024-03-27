---
marp: true
paginate: true
---

# slog LT

A quick overview and comparison of slog vs others

---

# slog?

- Introduced in Go 1.21: [Blog: https://go.dev/blog/slog](https://go.dev/blog/slog)
  - [Docs: https://pkg.go.dev/golang.org/x/exp/slog](https://pkg.go.dev/golang.org/x/exp/slog)
- Initially routes output through the old log `handler`
- Comes with two new handlers:
  - [`TextHandler`](https://pkg.go.dev/golang.org/x/exp/slog#TextHandler)
    - Basically [`logfmt`](https://betterstack.com/community/guides/logging/logfmt/)
  - [`JSONHandler`](https://pkg.go.dev/golang.org/x/exp/slog#JSONHandler)

---

# Simple Demo (Code)

```go
func main() {
	// Default loggers
	log.Print("1a: I am a basic out of the box `log`")                                  // Pre Go1.21 slog formatting
	slog.Info("1b: I am a basic out of the box `slog`", slog.String("arg", "argValue")) // Only difference to `log` is additional INFO, and the appended text

	// Create a text handler
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}
	slogger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	log.Print("2a: After creating a TextHandler slogger `log`")                                                 // Default logger: Same as 1a
	slog.Info("2b: After creating a TextHandler slogger `slog`", slog.String("arg", "argValue"))                // Default slogger: Same as 2a
	slogger.Info("2c: After creating a TextHandler slogger `slogger instance`", slog.String("arg", "argValue")) // TextHandler output

	// Change the default logger
	slog.SetDefault(slogger)
	log.Print("3a: After changing the Default `log`")                                  // Changes to TextHandler format
	slog.Info("3b: After changing the Default `slog`", slog.String("arg", "argValue")) // "Same" as 3a

	slogger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(slogger)
	log.Print("4a: After changing the Default to JSON Handler `log`")                                  // Changes to JSON format
	slog.Info("4b: After changing the Default to JSON Handler `slog`", slog.String("arg", "argValue")) // Changes to JSON format
}
```

See: [../sldemo/cmd/basic/main.go](../sldemo/cmd/basic/main.go)

---

# Simple Demo (Output)

```text
2024/03/27 06:35:26 1a: I am a basic out of the box `log`
2024/03/27 06:35:26 INFO 1b: I am a basic out of the box `slog` arg=argValue
2024/03/27 06:35:26 2a: After creating a TextHandler slogger `log`
2024/03/27 06:35:26 INFO 2b: After creating a TextHandler slogger `slog` arg=argValue
time=2024-03-27T06:35:26.496+09:00 level=INFO source=/Users/gavinthomas/Developer/slog-lt/sldemo/cmd/basic/main.go:22 msg="2c: After creating a TextHandler slogger `slogger instance`" arg=argValue
time=2024-03-27T06:35:26.496+09:00 level=INFO source=:0 msg="3a: After changing the Default `log`"
time=2024-03-27T06:35:26.496+09:00 level=INFO source=/Users/gavinthomas/Developer/slog-lt/sldemo/cmd/basic/main.go:27 msg="3b: After changing the Default `slog`" arg=argValue
{"time":"2024-03-27T06:35:26.496552+09:00","level":"INFO","msg":"4a: After changing the Default to JSON Handler `log`"}
{"time":"2024-03-27T06:35:26.496563+09:00","level":"INFO","source":{"function":"main.main","file":"/Users/gavinthomas/Developer/slog-lt/sldemo/cmd/basic/main.go","line":32},"msg":"4b: After changing the Default to JSON Handler `slog`","arg":"argValue"}
```

---

# Differences between `log` and `slog`

- `log` has: `Print`, `Fatal` and `Panic`
  - `Fatal` calls `os.Exit(1)` after logging
  - `Panic` calls `panic(s)` after logging
- `slog` has `Debug`, `Info`, `Warn`, `Error`
  - None of them call `os.Exit(1)`, or `panic(s)`, so you need to handle this if you migrate

---

# Using with Echo?

- An Echo instance has two loggers
  - `e.StdLogger` which is just stdlib `log` with a prefix of `echo:`
  - `e.Logger` which is a [custom logger implementation: https://github.com/labstack/gommon](https://github.com/labstack/gommon)
- We can replace the implementations, but instead, I just wouldn't use them.

```go
slogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
slog.SetDefault(slogger) // Overrides the default slog.* AND! log.* functions to use the same handler above
e.StdLogger = slog.NewLogLogger(slogger.Handler(), slog.LevelInfo) // Now override Echo's `log` logger
```

---

# DEMO: `x-request-id` & `trace-id`

- `x-request-id` is OPTIONALLY supplied from calling client
  - We want to add to Context (`ctx`), so that it's not necessary to pass request to dependencies, business logic etc
  - We want to log this from the `ctx`
- We want to generate a `trace-id`, and insert into the `ctx`
  - We want to log this from the `ctx`
  - This is relative easy to achieve because `slog.Handler` has a small interface

---

# DEMO: Hiding sensitive data

- `slog` provides a mechanism, but could be forgotten about by developers.

---

# Why not Zap?

- [Zap: https://github.com/uber-go/zap](https://github.com/uber-go/zap)
- Zap claims superior performance
- Zap has a built in sampling feature
- `zap.Error(error)` is included, and is equivalent is missing in `slog`
  - In `slog` we need to do something like `slog.String("error", err.Error())`
- `zap` has `Fatal` and `Panic` similar to `log`
  - There is also `DPanic` which panics in development, but errors logs in Production

---

# DEMO: Zap differences

- Recommend the non-sugared logger for performance and type-safety reasons.
- For hiding information, zap either:
  - Asks for you to implement `fmt.Stringer`, which is overly simplistic
  - Implement `zapcore.Encoder`, which is more complex than `slog`'s mechanism
- `zap` doesn't have a easy way to access add `ctx` to a logging call, and handle in middleware

---

# DEMO: Zap customer Encoder example

```go
func (e *SensitiveFieldEncoder) EncodeEntry(
    entry zapcore.Entry,
    fields []zapcore.Field,
) (*buffer.Buffer, error) {
    filtered := make([]zapcore.Field, 0, len(fields))

    for _, field := range fields {
        account, ok := field.Interface.(Account)
        if ok {
            account.Name = "[REDACTED]"
            account.AccountNumber = "[REDACTED]"
            field.Interface = account
        }

        filtered = append(filtered, field)
    }

    return e.Encoder.EncodeEntry(entry, filtered)
}

func NewSensitiveFieldsEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
    encoder := zapcore.NewJSONEncoder(config)
    return &SensitiveFieldEncoder{encoder, config}
}

func createLogger() *zap.Logger {
    . . .

    jsonEncoder := NewSensitiveFieldsEncoder(productionCfg)

    . . .

    return zap.New(samplingCore)
}
```

---

# Recomendations

- Old projects using `log` -> Consider switching to `slog`
- New projects:
  - Start with `slog`
    - Should be possible to switch to different handlers after introducing.
  - Only if you have specific requirements, use `zap`. (Performance issues, easy to configure sampling in Production)
- Personally, I found `slog` easier to work with.
