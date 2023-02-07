package main

import (
	"enforce-tool-versions/identifier"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	zlog.Info().Msg("starting program")

	for _, p := range []identifier.Program{identifier.Make, identifier.Git} {
		version, err := identifier.Identify(p, &zlog)
		if err != nil {
			zlog.Error().Err(err).Msg("failed to identify program")
			continue
		}
		zlog.Info().
			Str("program", identifier.GetProgramName(p)).
			Str("version", string(version)).
			Msg("identified program")
	}
}
