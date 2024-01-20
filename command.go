package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// Runs command and the output of the command.
func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	log.Trace().Str("cmd", cmd.String()).Msg("running command")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run [%s]: %s: %w", cmd.String(), out.String(), err)
	}
	return []byte(out.String()), nil
}
