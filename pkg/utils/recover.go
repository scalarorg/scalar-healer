package utils

import "github.com/rs/zerolog/log"

func Recover() {
	if r := recover(); r != nil {
		log.Error().Msgf("Recovered from panic: %v", r)
	}
}
