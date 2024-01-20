package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	alwaysYes := flag.Bool("y", false, "skip confirmations")
	force := flag.Bool("force", false, "force download even when up-to-date")
	level := flag.String("l", "info", "log level (trace, debug, info, warn, error, fatal, panic)")
	flag.Parse()

	lvl, err := zerolog.ParseLevel(*level)
	if err != nil {
		flag.Usage()
		return
	}

	// Set up the zerologger.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(lvl)
	log.Info().Bool("skip-confirmation", *alwaysYes).Msg("Updating/Installing proton-ge")

	log.Debug().
		Str("proton-dir", cfg.ProtonDir()).
		Str("compatibility-tools-dir", cfg.CompatibilityToolsDir()).
		Str("home", cfg.Home()).
		Msg("config loaded")

	a := NewApp(*alwaysYes, *force)
	if err := a.InstallOrUpdate(); err != nil {
		log.Fatal().Err(err).Msg("install or update")
	}
}
