package finder

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var ErrInvalidFormatFn = errors.New("invalid formatFn")
var ErrInvalidRoot = errors.New("invalid root")

// Source represents a directory source for finding paths.
type Source struct {
	Path         string
	OriginalPath string `json:"path"`
	Depth        uint8  `json:"depth"`

	// Function to format the output path.
	// Allows flexibility in other parts of the codebase (e.g., for testing).
	formatFn func(string) string
	resultCh chan<- string
}

// Find initiates the search based on the specified depth and format function.
func (s *Source) Find(resultCh chan<- string, formatFn func(string) string) error {
	if formatFn == nil {
		return ErrInvalidFormatFn
	}

	s.formatFn = formatFn

	expanded, err := homedir.Expand(s.OriginalPath)
	if err != nil {
		return err
	}

	s.Path = expanded
	s.resultCh = resultCh

	err = s.walkZero(s.Path)
	if err != nil {
		return err
	}

	if s.Depth == 0 {
		return nil
	}

	return s.walk(s.Path, 0)
}

func (s *Source) walkZero(root string) error {
	isDir, err := isPathDir(root)
	if err != nil {
		return err
	}

	if isDir {
		s.resultCh <- s.formatFn(root)
		return nil
	}

	return ErrInvalidRoot
}

func (s *Source) walk(root string, currDepth uint8) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	walkNext := func(p string) error {
		s.resultCh <- s.formatFn(p)

		if currDepth+1 < s.Depth {
			return s.walk(p, currDepth+1)
		}

		return nil
	}

	for _, entry := range entries {
		joined := filepath.Join(root, entry.Name())

		if entry.IsDir() {
			if err := walkNext(joined); err != nil {
				return err
			}

			continue
		}

		// is a symlink
		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}

		isDir, err := isPathDir(joined)
		if err != nil {
			return err
		}

		if isDir {
			if err := walkNext(joined); err != nil {
				return err
			}
		}
	}

	return nil
}

func isPathDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
