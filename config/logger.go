package config

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/openobserve"
)

func InitLogger() {
	if Env.IS_TEST {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimestampFieldName

	var writer io.Writer
	var subWriters []io.Writer = []io.Writer{}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	subWriters = append(subWriters, consoleWriter)

	if !Env.IS_TEST {
		o2Writer := openobserve.NewLogWriter(zerolog.InfoLevel)
		subWriters = append(subWriters, o2Writer)
	}

	writer = zerolog.MultiLevelWriter(subWriters...)
	log.Logger = log.Output(writer)
	log.Info().Msg("Logger initialized")
}
