package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/system"
)

// WithLogger will create an event suitable Options from CLI ones.
func (opts *Options) WithLogger(outputs system.Outputs) (*event.Properties, error) {
	zc := zap.NewProductionConfig()
	cfg := zap.NewProductionEncoderConfig()
	if opts.Verbose {
		cfg = zap.NewDevelopmentEncoderConfig()
	}
	cfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	var encoder zapcore.Encoder
	switch opts.Output {
	case HumanReadable:
		if !opts.Verbose {
			cfg.CallerKey = ""
		}
		cfg.ConsoleSeparator = " "
		cfg.EncodeLevel = alignCapitalColorLevelEncoder
		cfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.StampMilli)
		encoder = zapcore.NewConsoleEncoder(cfg)
	case YAML:
		encoder = &yamlEncoder{zapcore.NewJSONEncoder(cfg)}
	case JSON:
		encoder = zapcore.NewJSONEncoder(cfg)
	}
	sink := zapcore.AddSync(outputs.OutOrStdout())
	errSink := zapcore.AddSync(outputs.ErrOrStderr())
	zcore := zapcore.NewCore(encoder, sink, zc.Level)
	log := zap.New(
		zcore, buildOptions(zc, errSink)...,
	)

	return &event.Properties{
		KnPluginOptions: opts.KnPluginOptions,
		Log:             log.Sugar(),
	}, nil
}

func alignCapitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	spaces := len(zapcore.FatalLevel.CapitalString()) - len(l.CapitalString())
	if spaces > 0 {
		enc.AppendString(strings.Repeat(" ", spaces))
	}
	zapcore.CapitalColorLevelEncoder(l, enc)
}

func buildOptions(cfg zap.Config, errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	return opts
}

type yamlEncoder struct {
	zapcore.Encoder
}

func (y *yamlEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := y.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, unexpected(err)
	}
	var v interface{}
	err = json.Unmarshal(buf.Bytes(), &v)
	if err != nil {
		return nil, unexpected(err)
	}
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return nil, unexpected(err)
	}
	buf = buffer.NewPool().Get()
	_, _ = buf.Write([]byte("---\n"))
	if _, err = buf.Write(bytes); err != nil {
		return nil, unexpected(err)
	}
	return buf, nil
}

func unexpected(err error) error {
	if errors.Is(err, event.ErrUnexpected) {
		return err
	}
	return fmt.Errorf("%w: %w", event.ErrUnexpected, err)
}
