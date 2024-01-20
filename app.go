package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Sentinel errors.
var (
	ErrUpToDate = errors.New("already up-to-date")
)

// App acts as an updater and installer for proton-ge
type App struct {
	alwaysYes bool
	force     bool
}

// NewApp creates a new application that's configured to install or update
// proton-ge
func NewApp(alwaysYes bool, force bool) *App {
	return &App{
		alwaysYes: alwaysYes,
		force:     force,
	}
}

// InstallOrUpdate will install proton-ge for the first time or update it if
// there is a new version available. Optionally you can force the app to download
// the latest proton-ge version even if it's the same as your last installed.
func (a *App) InstallOrUpdate() error {
	if err := a.initializeDirectories(); err != nil {
		return fmt.Errorf("initialize directories: %w", err)
	}

	log.Trace().Msg("downloading latest release info")
	releaseInfo, err := downloadLatestReleaseInfo()
	if err != nil {
		return fmt.Errorf("download release info: %w", err)
	}

	// Print out currently installed versions and latest version and check if we're up-to-date
	if err := a.displayVersions(releaseInfo, a.force); errors.Is(err, ErrUpToDate) {
		log.Info().Msg("Already up-to-date!")
		return nil
	}

	if !a.alwaysYes && !prompt("Would you like to continue with installation?", true) {
		return nil
	}

	log.Info().Msg("Downloading latest proton-ge release...")
	if err := downloadLatestProtonRelease(releaseInfo); err != nil {
		return err
	}

	tarBallName := releaseInfo.TarName()
	if !a.alwaysYes && !prompt(fmt.Sprintf("Would you like to install %s?", color.GreenString(tarBallName)), true) {
		log.Info().Msg("Cleaning up...")
		// Remove all the downloaded files
		for _, asset := range releaseInfo.Assets {
			if err := os.Remove(asset.Name); err != nil {
				return fmt.Errorf("remove %s: %w", asset.Name, err)
			}
		}
		return nil
	}

	log.Info().Msg("Installing latest Proton-GE")
	if err := a.decompressProtonGE(tarBallName); err != nil {
		return err
	}

	if err := a.moveDownloadedToArchive(tarBallName); err != nil {
		return err
	}

	// Remove the sha512sum file
	sha512FileName := releaseInfo.Sha512SumName()
	if err := os.Remove(sha512FileName); err != nil {
		return fmt.Errorf("remove %s: %w", sha512FileName, err)
	}

	log.Info().Msg("Installed, please restart Steam")
	if log.Logger.GetLevel() <= zerolog.DebugLevel {
		entries, err := os.ReadDir(cfg.CompatibilityToolsDir())
		if err != nil {
			return fmt.Errorf("read dir: %s: %w", cfg.CompatibilityToolsDir(), err)
		}

		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		log.Debug().Strs("installed", names).Msg("compatibility tools")
	}

	return nil
}

func (a *App) moveDownloadedToArchive(n string) error {
	path := cfg.ProtonDir() + "/"
	_, err := runCommand("mv", n, path)
	if err != nil {
		return fmt.Errorf("move file: %w", err)
	}

	log.Debug().Str("path", path+n).Msg("moved directory")

	return nil
}

func (a *App) decompressProtonGE(name string) error {
	if _, err := runCommand("tar", "-xf", name, "-C", cfg.CompatibilityToolsDir()); err != nil {
		return err
	}

	return nil
}

func printDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}
	for _, e := range entries {
		fmt.Printf("    %s", color.BlueString(e.Name()+"/"))
	}
	return nil
}

// displays the currently installed version of protonge and the latest
func (a *App) displayVersions(releaseInfo *Response, force bool) error {
	entries, err := os.ReadDir(cfg.ProtonDir())
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	var latestTarBall string
	// Print the latest version's tar ball file name as comparison.
	for _, a := range releaseInfo.Assets {
		if strings.Contains(a.Name, "tar.gz") {
			latestTarBall = a.Name
		}
	}

	var names strings.Builder
	for i := range entries {
		// Skip any directories in the path
		if entries[i].IsDir() {
			continue
		}

		if entries[i].Name() == latestTarBall && !force {
			return ErrUpToDate
		}

		names.WriteString("    ")
		names.WriteString(color.BlueString(entries[i].Name()))
	}

	fmt.Println("Found the following installed last:")
	fmt.Println(names.String())
	fmt.Printf("To be downloaded and installed:\n    %s\n",
		color.GreenString(latestTarBall))

	return nil
}

// Create working directories if they don't exist yet.
func (a *App) initializeDirectories() error {
	if err := createIfNotExists(cfg.ProtonDir()); err != nil {
		return fmt.Errorf("create tmp-dir: %w", err)
	}
	if err := createIfNotExists(cfg.CompatibilityToolsDir()); err != nil {
		return fmt.Errorf("create steam-compat dir: %w", err)
	}
	return nil
}

// Create a directory recursively if it does not yet exist.
func createIfNotExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}
	return nil
}
