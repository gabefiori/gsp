package main

import (
	"log"

	"github.com/gabefiori/gsp/internal/cli"
)

var version = "unknown"

func main() {
	if err := cli.Run(version); err != nil {
		log.Fatal(err)
	}
}
