package revenium

import "testing"

func TestNormalizeModelName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bare model name gets both prefixes",
			input:    "flux/dev",
			expected: "fal_ai/fal-ai/flux/dev",
		},
		{
			name:     "fal-ai endpoint prefix gets litellm prefix prepended",
			input:    "fal-ai/flux/dev",
			expected: "fal_ai/fal-ai/flux/dev",
		},
		{
			name:     "already correct format passes through",
			input:    "fal_ai/fal-ai/flux/dev",
			expected: "fal_ai/fal-ai/flux/dev",
		},
		{
			name:     "litellm prefix without fal-ai segment gets segment inserted",
			input:    "fal_ai/flux/dev",
			expected: "fal_ai/fal-ai/flux/dev",
		},
		{
			name:     "already correct with nested path",
			input:    "fal_ai/fal-ai/flux-pro/v1.1",
			expected: "fal_ai/fal-ai/flux-pro/v1.1",
		},
		{
			name:     "fal-ai prefix with nested path",
			input:    "fal-ai/flux-pro/v1.1",
			expected: "fal_ai/fal-ai/flux-pro/v1.1",
		},
		{
			name:     "idempotent - calling twice produces same result",
			input:    "fal_ai/fal-ai/flux/dev",
			expected: "fal_ai/fal-ai/flux/dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeModelName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeModelName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeModelNameIdempotent(t *testing.T) {
	inputs := []string{
		"flux/dev",
		"fal-ai/flux/dev",
		"fal_ai/flux/dev",
		"fal_ai/fal-ai/flux/dev",
	}

	for _, input := range inputs {
		first := normalizeModelName(input)
		second := normalizeModelName(first)
		if first != second {
			t.Errorf("not idempotent: normalizeModelName(%q) = %q, but normalizeModelName(%q) = %q",
				input, first, first, second)
		}
	}
}
