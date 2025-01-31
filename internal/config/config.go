package config

import (
	"errors"
	"os"

	"github.com/gabefiori/gsp/internal/finder"
	"github.com/mitchellh/go-homedir"
)

// Config represents the configuration structure for the application.
type Config struct {
	// List of sources to be used by the finder
	Sources []finder.Source

	// Flag to indicate if output should be expanded
	// Useful to hide the user's home directory
	ExpandOutput bool

	// Flag to indicate if measurement should be performed
	Measure bool

	// Flag to list results
	List bool

	// Selector for displaying the projects
	Selector string

	// Flag to display only unique projects.
	Unique bool

	// Type of sorting.
	Sort string
}

type LoadParams struct {
	Selector     string
	Sort         string
	Path         string
	ExpandOutput int8
	Unique       int8
	Measure      bool
	List         bool
}

// Load reads the configuration from a JSON file at the specified path.
func Load(params *LoadParams) (*Config, error) {
	path, err := homedir.Expand(params.Path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var cfg Config

	cfg.ExpandOutput = true

	parser := NewParser(file, &cfg)
	if err := parser.Run(); err != nil {
		return nil, err
	}

	cfg.Measure = params.Measure
	cfg.List = params.List

	if params.ExpandOutput != 0 {
		cfg.ExpandOutput = params.ExpandOutput == 1
	}

	if params.Unique != 0 {
		cfg.Unique = params.Unique == 1
	}

	if params.Selector != "" {
		cfg.Selector = params.Selector
	}

	if cfg.Selector == "" {
		return nil, errors.New("invalid selector")
	}

	if params.Sort != "" {
		cfg.Sort = params.Sort
	}

	return &cfg, nil
}
