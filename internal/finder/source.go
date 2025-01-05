package finder

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charlievieth/fastwalk"
	"github.com/mitchellh/go-homedir"
)

var ErrInvalidFormatFn = errors.New("invalid formatFn")

// Source represents a directory source for finding paths.
type Source struct {
	Path         string
	OriginalPath string `json:"path"`
	Depth        uint8  `json:"depth"`

	// Function to format the output path.
	// Allows flexibility in other parts of the codebase (e.g., for testing).
	formatFn func(string) string
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

	// This prevents a wrong calculation of depth
	s.Path = strings.TrimSuffix(s.Path, "/")

	// Fastwalk is generally faster for deep directory structures,
	// but for shallow searches, using just [os.ReadDir] or [os.Stat] is more efficient.
	if s.Depth == 0 {
		return s.depthZero(resultCh)
	}

	if s.Depth == 1 {
		return s.depthOne(resultCh)
	}

	return s.depthGreater(resultCh)
}

func (s *Source) depthZero(resultCh chan<- string) error {
	isDir, err := isPathDir(s.Path)
	if err != nil {
		return err
	}

	if isDir {
		resultCh <- s.formatFn(s.Path)
	}

	return nil
}

func (s *Source) depthOne(resultCh chan<- string) error {
	entries, err := os.ReadDir(s.Path)
	if err != nil {
		return err
	}

	resultCh <- s.formatFn(s.Path)

	for _, entry := range entries {
		path := filepath.Join(s.Path, entry.Name())
		isDir, err := isPathDir(path)

		if err != nil {
			return err
		}

		if isDir {
			resultCh <- s.formatFn(path)
		}
	}
	return nil
}

func (s *Source) depthGreater(resultCh chan<- string) error {
	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if currentDepth(s.Path, path) > s.Depth {
			return fs.SkipDir
		}

		resultCh <- s.formatFn(path)
		return nil
	}

	// TODO: Check other options that fastwalk has for performance.
	err := fastwalk.Walk(
		&fastwalk.Config{Follow: true},
		s.Path,
		walkFn,
	)

	return err
}

// Calculate the depth of the 'curr' path relative to the 'root' directory.
// The depth is defined as the number of directory levels between 'root' and 'curr'.
//
// For example:
//
// If root is "/home/user" and curr is "/home/user/documents",
// the relative path is "documents", which means a depth of 1.
func currentDepth(root, curr string) uint8 {
	relPath, _ := filepath.Rel(root, curr)
	return uint8(len(strings.Split(relPath, string(os.PathSeparator))))
}

func isPathDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
