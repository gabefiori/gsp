package finder

import (
	"log"
	"strings"
	"sync"
)

type FinderOpts struct {
	Sources  []Source
	HomeDir  string
	ResultCh chan string
	SortType SortType
	Unique   bool
}

// Run executes the package finder using the provided options.
// Any error encountered within this function is considered fatal and will terminate the program.
//
// Each source runs its [Find] method in a separate goroutine.
func Run(opts *FinderOpts) {
	var wg sync.WaitGroup
	var pipeCh chan string

	ch := opts.ResultCh
	usePipe := opts.SortType != NoSort || opts.Unique

	if usePipe {
		pipeCh = make(chan string, cap(opts.ResultCh))
		ch = pipeCh
	}

	for _, source := range opts.Sources {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := source.Find(ch, func(s string) string {
				if strings.HasPrefix(source.OriginalPath, "~") {
					return "~" + strings.TrimPrefix(s, opts.HomeDir)
				}

				return s
			})

			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	if !usePipe {
		wg.Wait()
		close(opts.ResultCh)

		return
	}

	go func() {
		defer close(opts.ResultCh)

		unique := make(map[string]struct{})
		results := make([]string, 0, 50)

		for r := range pipeCh {
			if opts.Unique {
				if _, exists := unique[r]; exists {
					continue
				}

				unique[r] = struct{}{}
			}

			results = append(results, r)
		}

		if opts.SortType != NoSort {
			sortResults(results, opts.SortType)
		}

		for _, r := range results {
			opts.ResultCh <- r
		}
	}()

	wg.Wait()
	close(pipeCh)
}
