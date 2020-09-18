package main

import (
	"fmt"
	"os"

	"samhofi.us/x/infobot/cmd"
)

// Exit code on failure
const exitFail = 1

func main() {
	if err := cmd.Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}
