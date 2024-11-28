package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gabefiori/gsp/internal/config"
	"github.com/gabefiori/gsp/internal/finder"
	"github.com/gabefiori/gsp/internal/selector"
	"github.com/mitchellh/go-homedir"
)

// Run executes the main logic of the application using the provided configuration.
func Run(cfg *config.Config) error {
	home, err := homedir.Dir()

	if err != nil {
		return err
	}

	// Channel to receive output (string) from the finder.
	//
	// This channel is also passed to the selector to populate its input.
	resultCh := make(chan string, 3)
	measureStart := time.Now()

	go finder.Run(&finder.FinderOpts{
		ResultCh: resultCh,
		HomeDir:  home,
		Sources:  cfg.Sources,
		SortType: finder.SortTypeFromStr(cfg.Sort),
		Unique:   true,
	})

	// If output expansion is not enabled, set the home directory to "~".
	// This is useful for hiding the user's home directory.
	if !cfg.ExpandOutput {
		home = "~"
	}

	// If measurement is enabled, count the number of projects found
	// and the time taken to find the projects.
	if cfg.Measure {
		var count int

		for range resultCh {
			count++
		}

		measureEnd := time.Since(measureStart).String()
		msg := fmt.Sprintf("Took %s (%d projects)", measureEnd, count)

		_, err = os.Stdout.WriteString(msg)
		return err
	}

	// If listing is enabled, print the results to stdout in batches.
	//
	// Using io.Copy is more efficient for larger batches of data, as it minimizes
	// the number of system calls and leverages internal buffering.
	if cfg.List {
		batchSize := 50
		batchCount := 0

		buf := new(bytes.Buffer)

		for r := range resultCh {
			if _, err := buf.WriteString(r + "\n"); err != nil {
				return err
			}

			batchCount++

			if batchCount >= batchSize {
				if _, err := io.Copy(os.Stdout, buf); err != nil {
					return err
				}

				buf.Reset()
				batchCount = 0
			}
		}

		_, err = io.Copy(os.Stdout, buf)
		return err
	}

	t := selector.TypeFromStr(cfg.Selector)
	s, err := selector.New(t)

	if err != nil {
		return err
	}

	result, err := s.Run(resultCh)

	// If the selector is canceled, result will be empty.
	if err != nil || result == "" {
		return err
	}

	// The first character ("~") of the result is skipped.
	// It's only used for display inside the selector.
	//
	// The expanded version of the result must be used;
	// otherwise, it will not be able to be consumed by other programs.
	_, err = os.Stdout.WriteString(home + result[1:])
	return err
}
