package main

import (
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/rs/zerolog/log"
)

// Returns true if the user responded with Yes/Yy
//
// You can provide defaultYes false to make No/Nn the default option.
//
// When the user hits enter it chooses the default option.
func prompt(msg string, defaultYes bool) bool {
	var decisions strings.Builder
	if defaultYes {
		decisions.WriteString("[Yy(default)] [Nn]")	
	} else {
		decisions.WriteString("[Yy] [Nn(default)]")
	}

	fmt.Printf("\n%s: %s\n", msg, decisions.String())
	for true {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			log.Warn().Err(err).Msg("get single key")
			break
		}
		// Check the response type.
		switch char {
		case 'Y' | 'y':
			return true
		case 'N' | 'n':
			return false
		default:
			if key == keyboard.KeyEnter {
				return defaultYes
			}
		}
	}
	return false
}
