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
	depths := []int{0, 1, 2, 3}

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

	for _, depth := range depths {
		createNestedDirs(b, baseDir, 0, depth)
		source := Source{OriginalPath: baseDir, Depth: uint8(depth)}

		for _, tt := range tests {
			b.Run(fmt.Sprintf("Depth_%d/%s", depth, tt.name), func(b *testing.B) {
				resultCh := make(chan string, 3)
				opts := &FinderOpts{
					Sources:  []Source{source, source, source},
					HomeDir:  baseDir,
					ResultCh: resultCh,
					SortType: tt.sortType,
					Unique:   tt.unique,
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					go Run(opts)

					for range resultCh {
					}

					//FIXME: this affects the benchmark.
					resultCh = make(chan string, 3)
					opts.ResultCh = resultCh
				}
			})
		}
	}

	assert.NoError(b, os.RemoveAll(baseDir))
}

func createNestedDirs(b *testing.B, baseDir string, currentDepth, maxDepth int) {
	if currentDepth > maxDepth {
		return
	}

	for i := 0; i < 5; i++ {
		dirPath := filepath.Join(baseDir, fmt.Sprintf("dir-%d", i))
		assert.NoError(b, os.MkdirAll(dirPath, 0755))
		createNestedDirs(b, dirPath, currentDepth+1, maxDepth)
	}
}
