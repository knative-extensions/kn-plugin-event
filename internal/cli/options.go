package cli

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cardil/kn-event/internal/event"
	"github.com/ghodss/yaml"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// WithLogger will create an event suitable OptionsArgs from CLI ones.
func (opts *OptionsArgs) WithLogger() *event.Properties {
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
	sink := zapcore.AddSync(opts.OutWriter)
	errSink := zapcore.AddSync(opts.ErrWriter)
	zcore := zapcore.NewCore(encoder, sink, zc.Level)
	log := zap.New(
		zcore, buildOptions(zc, errSink)...,
	)

	return &event.Properties{
		KnPluginOptions: opts.KnPluginOptions,
		Log:             log.Sugar(),
	}
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
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(buf.Bytes(), &v)
	if err != nil {
		return nil, err
	}
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	buf = buffer.NewPool().Get()
	_, _ = buf.Write([]byte("---\n"))
	_, err = buf.Write(bytes)
	return buf, err
}
