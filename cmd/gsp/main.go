package main

import (
	"fmt"
	"os"

	"github.com/gabefiori/gsp/internal/cli"
)

var version = "unknown"

func main() {
	if err := cli.Run(version); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
