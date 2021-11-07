package main

import (
	"os"

	"github.com/sample/sample-server/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
