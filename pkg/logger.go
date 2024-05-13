package pkg

import (
	"github.com/rs/zerolog"
	"os"
)

func InitLogger(debug bool) zerolog.Logger {
	var l zerolog.Logger

	if debug {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		l = zerolog.New(os.Stderr)
	}

	return l.With().Timestamp().Logger()
}
