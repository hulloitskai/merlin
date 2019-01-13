package cmd

import (
	"fmt"

	"github.com/spf13/pflag"
)

func showHelp() {
	fmt.Print("Merlin is a system for accessing company finance data from " +
		"EDGAR.\n\n")
	pflag.Usage()
}
