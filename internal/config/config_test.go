package config

import (
	"os"
	"testing"

	"github.com/gabefiori/gsp/internal/finder"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config")
	assert.NoError(t, err)

	defer os.Remove(tempFile.Name())

	sampleConfig := `
		source = 1:~/test_1/test_1
		source = 2:~/test_2/test_2
		expand-output = false
		selector = test-selector
		unique = false
		sort = asc
	`

	_, err = tempFile.WriteString(sampleConfig)
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	sources := []finder.Source{
		{OriginalPath: "~/test_1/test_1", Depth: uint8(1)},
		{OriginalPath: "~/test_2/test_2", Depth: uint8(2)},
	}

	t.Run("With parameters specified", func(t *testing.T) {
		params := &LoadParams{
			Selector:     "other",
			Sort:         "asc",
			Path:         tempFile.Name(),
			ExpandOutput: 1,
			Unique:       1,
			Measure:      true,
			List:         true,
		}

		cfg, err := Load(params)
		assert.NoError(t, err)

		assert.Equal(t, true, cfg.ExpandOutput)
		assert.Equal(t, true, cfg.Measure)
		assert.Equal(t, true, cfg.List)
		assert.Equal(t, params.Selector, cfg.Selector)
		assert.Equal(t, true, cfg.Unique)
		assert.Equal(t, "asc", cfg.Sort)
		assert.Equal(t, sources, cfg.Sources)
	})

	t.Run("With minimal parameters", func(t *testing.T) {
		params := &LoadParams{
			Path: tempFile.Name(),
		}

		cfg, err := Load(params)
		assert.NoError(t, err)

		assert.Equal(t, false, cfg.ExpandOutput)
		assert.Equal(t, false, cfg.Measure)
		assert.Equal(t, false, cfg.List)
		assert.Equal(t, "test-selector", cfg.Selector)
		assert.Equal(t, false, cfg.Unique)
		assert.Equal(t, "asc", cfg.Sort)
		assert.Equal(t, sources, cfg.Sources)
	})
}

func BenchmarkLoad(b *testing.B) {
	tempFile, err := os.CreateTemp("", "config.json")
	assert.NoError(b, err)

	defer os.Remove(tempFile.Name())
	sampleConfig := `
		source = 1:~/test_1/test_1
		source = 2:~/test_2/test_2
		expand-output = false
		selector = test-selector
		unique = false
		sort = asc
	`
	_, err = tempFile.WriteString(sampleConfig)
	assert.NoError(b, err)

	err = tempFile.Close()
	assert.NoError(b, err)

	params := &LoadParams{
		Path: tempFile.Name(),
	}

	for i := 0; i < b.N; i++ {
		_, _ = Load(params)
	}
}
