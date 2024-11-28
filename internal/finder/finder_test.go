package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRun(b *testing.B) {
	tempDir := b.TempDir()
	baseDir := filepath.Join(tempDir, "base")
	numDirs := 200

	for i := 0; i < numDirs; i++ {
		dirPath := filepath.Join(baseDir, fmt.Sprintf("dir-%d", i))
		assert.NoError(b, os.MkdirAll(dirPath, 0755))
	}

	source := Source{Path: baseDir, Depth: 3}

	tests := []struct {
		name     string
		sortType SortType
		unique   bool
	}{
		{"NoSort_NonUnique", NoSort, false},
		{"NoSort_Unique", NoSort, true},
		{"AscSort_NonUnique", AscSort, false},
		{"AscSort_Unique", AscSort, true},
		{"DescSort_NonUnique", DescSort, false},
		{"DescSort_Unique", DescSort, true},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultCh := make(chan string)

				opts := &FinderOpts{
					Sources:  []Source{source},
					HomeDir:  baseDir,
					ResultCh: resultCh,
					SortType: tt.sortType,
					Unique:   tt.unique,
				}

				go Run(opts)

				for range resultCh {
				}
			}
		})
	}

	assert.NoError(b, os.RemoveAll(baseDir))
}
