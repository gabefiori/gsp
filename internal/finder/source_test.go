package finder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	depth1Dir := filepath.Join(tempDir, "depth1")
	depth2Dir := filepath.Join(tempDir, "depth1", "depth2")
	depth3Dir := filepath.Join(tempDir, "depth1", "depth2", "depth3")
	depth4Dir := filepath.Join(tempDir, "depth1", "depth2", "depth3", "depth4")

	assert.NoError(t, os.Mkdir(depth1Dir, 0755))
	assert.NoError(t, os.Mkdir(depth2Dir, 0755))
	assert.NoError(t, os.Mkdir(depth3Dir, 0755))
	assert.NoError(t, os.Mkdir(depth4Dir, 0755))

	symlinkPath := filepath.Join(depth1Dir, "symlink_to_depth2")
	assert.NoError(t, os.Symlink(depth2Dir, symlinkPath))

	tests := []struct {
		name     string
		depth    uint8
		expected []string
	}{
		{
			name:  "depthZero",
			depth: 0,
			expected: []string{
				tempDir,
			},
		},
		{
			name:  "depthOne",
			depth: 1,
			expected: []string{
				tempDir,
				depth1Dir,
			},
		},
		{
			name:  "depthGreater",
			depth: 3,
			expected: []string{
				tempDir,
				depth1Dir,
				depth2Dir,
				depth3Dir,
				symlinkPath,
				filepath.Join(symlinkPath, "depth3"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := Source{Path: tempDir, Depth: tt.depth}
			resultCh := make(chan string)

			go func() {
				defer close(resultCh)
				err := source.Find(resultCh, func(s string) string {
					return s
				})

				assert.NoError(t, err)
			}()

			var paths []string
			for path := range resultCh {
				paths = append(paths, path)
			}

			for _, expected := range tt.expected {
				assert.Contains(t, paths, expected)
			}

			assert.Len(t, paths, len(tt.expected))
		})
	}
}
