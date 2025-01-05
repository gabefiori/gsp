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

type Mode int8

const (
	ModeSelector Mode = iota
	ModeList
	ModeMeasure
)

type App struct {
	// Channel to receive output (string) from the finder.
	// This channel is also passed to the selector to populate its input.
	ch           chan string
	home         string
	sources      []finder.Source
	selectorType selector.Type
	sortType     finder.SortType
	Mode
}

func New(cfg *config.Config) (*App, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	// If output expansion is not enabled, set the home directory to "~".
	// This is useful for hiding the user's home directory.
	if !cfg.ExpandOutput {
		home = "~"
	}

	var m Mode
	if cfg.List {
		m = ModeList
	} else if cfg.Measure {
		m = ModeMeasure
	}

	return &App{
		Mode:         m,
		home:         home,
		sources:      cfg.Sources,
		ch:           make(chan string, len(cfg.Sources)),
		selectorType: selector.TypeFromStr(cfg.Selector),
		sortType:     finder.SortTypeFromStr(cfg.Sort),
	}, nil
}

// Run executes the main logic of the application.
func (a *App) Run() error {
	measureStart := time.Now()

	go finder.Run(&finder.FinderOpts{
		ResultCh: a.ch,
		HomeDir:  a.home,
		Sources:  a.sources,
		SortType: a.sortType,
		Unique:   true,
	})

	switch a.Mode {
	case ModeMeasure:
		return a.measure(measureStart)
	case ModeList:
		return a.list()
	default:
		return a.selector()
	}
}

func (a *App) selector() error {
	s, err := selector.New(a.selectorType)
	if err != nil {
		return err
	}

	result, err := s.Run(a.ch)
	// If the selector is canceled, result will be empty.
	if err != nil || result == "" {
		return err
	}

	// The first character ("~") of the result is skipped.
	// It's only used for display inside the selector.
	//
	// The expanded version of the result must be used;
	// otherwise, it will not be able to be consumed by other programs.
	_, err = os.Stdout.WriteString(a.home + result[1:])
	return err
}

func (a *App) measure(start time.Time) error {
	var count int

	for range a.ch {
		count++
	}

	measureEnd := time.Since(start).String()
	msg := fmt.Sprintf("Took %s (%d projects)", measureEnd, count)

	_, err := os.Stdout.WriteString(msg)
	return err
}

func (a *App) list() error {
	size, count := 50, 0
	buf := new(bytes.Buffer)

	for r := range a.ch {
		if _, err := buf.WriteString(r + "\n"); err != nil {
			return err
		}

		count++

		if count >= size {
			if _, err := io.Copy(os.Stdout, buf); err != nil {
				return err
			}

			buf.Reset()
			count = 0
		}
	}

	_, err := io.Copy(os.Stdout, buf)
	return err
}
