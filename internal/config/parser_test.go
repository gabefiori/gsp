package config

import (
	"strings"
	"testing"

	"github.com/gabefiori/gsp/internal/finder"
	"github.com/stretchr/testify/assert"
)

func TestParser_Run(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  *Config
		expectErr bool
	}{
		{
			name: "Valid config",
			input: `
				source = 1:~/test_1/test_1
				source = 2:~/test_2/test_2
				expand-output = true
				selector = test-selector
				unique = true
				sort = desc
			`,
			expected: &Config{
				Sources: []finder.Source{
					{OriginalPath: "~/test_1/test_1", Depth: 1},
					{OriginalPath: "~/test_2/test_2", Depth: 2},
				},
				ExpandOutput: true,
				Selector:     "test-selector",
				Unique:       true,
				Sort:         "desc",
			},
			expectErr: false,
		},
		{
			name: "Invalid source format",
			input: `
				source = invalid-source
			`,
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Missing equals sign",
			input: `
				source 1:~/test_1/test_1
			`,
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Empty line",
			input: `
				source = 1:~/test_1/test_1

				selector = test-selector
			`,
			expected: &Config{
				Sources: []finder.Source{
					{OriginalPath: "~/test_1/test_1", Depth: 1},
				},
				Selector: "test-selector",
			},
			expectErr: false,
		},
		{
			name: "Comment line",
			input: `
				# This is a comment
				source = 1:~/test_1/test_1
			`,
			expected: &Config{
				Sources: []finder.Source{
					{OriginalPath: "~/test_1/test_1", Depth: 1},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			parser := NewParser(strings.NewReader(tt.input), cfg)

			err := parser.Run()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, cfg)
			}
		})
	}
}
