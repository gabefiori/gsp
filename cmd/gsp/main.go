package main

import (
	"log"

	"github.com/gabefiori/gsp/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
