package config

import (
	"io"
	"os"

	"github.com/0xdavid7/goes-template/pkg/openobserve"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	if Env.IsTest {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	var writer io.Writer
	var subWriters []io.Writer = []io.Writer{}

	o2Writer := openobserve.NewLogWriter(zerolog.InfoLevel)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	if Env.IsProd || Env.IsStaging {
		subWriters = append(subWriters, os.Stdout)
		subWriters = append(subWriters, o2Writer)
	} else {
		subWriters = append(subWriters, consoleWriter)
		subWriters = append(subWriters, o2Writer)
	}

	if len(subWriters) >= 1 {
		writer = zerolog.MultiLevelWriter(subWriters...)
	} else {
		writer = os.Stdout
	}

	log.Logger = log.Output(writer)
	log.Info().Msg("Logger initialized")
}
