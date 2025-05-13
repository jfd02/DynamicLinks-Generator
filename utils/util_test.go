package utils

import (
	"os"
	"testing"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	_ = log
	os.Exit(m.Run())
}

func TestValidateURLScheme(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "invalid scheme - ftp",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "invalid scheme - file",
			url:     "file:///path/to/file",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "not a url",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURLScheme(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURLScheme() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid https URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "valid http URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "valid URL with path",
			input:    "https://example.com/path",
			expected: true,
		},
		{
			name:     "valid URL with query params",
			input:    "https://example.com?param=value",
			expected: true,
		},
		{
			name:     "valid URL with port",
			input:    "https://example.com:8080",
			expected: true,
		},
		{
			name:     "missing scheme",
			input:    "example.com",
			expected: false,
		},
		{
			name:     "missing host",
			input:    "https://",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid URL format",
			input:    "not a url",
			expected: false,
		},
		{
			name:     "valid URL with subdomain",
			input:    "https://sub.example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsURL(tt.input)
			if result != tt.expected {
				t.Errorf("IsURL(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsNumericString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "empty string",
			s:    "",
			want: true,
		},
		{
			name: "single digit",
			s:    "5",
			want: true,
		},
		{
			name: "multiple digits",
			s:    "12345",
			want: true,
		},
		{
			name: "contains letter",
			s:    "123a45",
			want: false,
		},
		{
			name: "contains space",
			s:    "123 45",
			want: false,
		},
		{
			name: "contains special character",
			s:    "123-45",
			want: false,
		},
		{
			name: "all non-numeric",
			s:    "abc",
			want: false,
		},
		{
			name: "unicode numbers",
			s:    "１２３", // Full-width numbers
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumericString(tt.s); got != tt.want {
				t.Errorf("IsNumericString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDomainAllowed(t *testing.T) {
	allowList := []string{
		"example.com",
		"allowed.com",
		"sub.allowed.com",
	}

	tests := []struct {
		name      string
		rawLink   string
		allowList []string
		want      bool
	}{
		{
			name:      "exact match",
			rawLink:   "https://example.com",
			allowList: allowList,
			want:      true,
		},
		{
			name:      "subdomain match",
			rawLink:   "https://sub.example.com",
			allowList: allowList,
			want:      false,
		},
		{
			name:      "not in allow list",
			rawLink:   "https://notallowed.com",
			allowList: allowList,
			want:      false,
		},
		{
			name:      "invalid URL",
			rawLink:   "not a valid url",
			allowList: allowList,
			want:      false,
		},
		{
			name:      "case insensitive",
			rawLink:   "https://EXAMPLE.com",
			allowList: allowList,
			want:      true,
		},
		{
			name:      "empty allow list",
			rawLink:   "https://example.com",
			allowList: []string{},
			want:      false,
		},
		{
			name:      "with path and query",
			rawLink:   "https://example.com/path?query=value",
			allowList: allowList,
			want:      true,
		},
		{
			name:      "with port",
			rawLink:   "https://example.com:8080",
			allowList: allowList,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDomainAllowed(tt.allowList, tt.rawLink)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateDynamicLinkPath(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "short path",
			length: 6,
		},
		{
			name:   "medium path",
			length: 12,
		},
		{
			name:   "long path",
			length: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := make(map[string]bool)
			const iterations = 1000

			for range iterations {
				path := GenerateDynamicLinkPath(tt.length)

				// Test length
				if len(path) != tt.length {
					t.Errorf("generateDynamicLinkPath(%d) length = %d, want %d", tt.length, len(path), tt.length)
				}

				// Test character set
				for _, r := range path {
					if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
						t.Errorf("generateDynamicLinkPath(%d) contains invalid character: %c", tt.length, r)
					}
				}

				// Test uniqueness
				if paths[path] {
					t.Errorf("generateDynamicLinkPath(%d) generated duplicate path: %s", tt.length, path)
				}
				paths[path] = true
			}

			// Test distribution of characters
			charCount := make(map[rune]int)
			for path := range paths {
				for _, r := range path {
					charCount[r]++
				}
			}

			// Check if we have both letters and numbers
			hasLetters := false
			hasNumbers := false
			for r := range charCount {
				if unicode.IsLetter(r) {
					hasLetters = true
				}
				if unicode.IsDigit(r) {
					hasNumbers = true
				}
			}

			if !hasLetters || !hasNumbers {
				t.Errorf("generateDynamicLinkPath(%d) does not generate both letters and numbers", tt.length)
			}
		})
	}
}

func TestCleanHost(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid host with https",
			raw:     "https://example.com",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "valid host with http",
			raw:     "http://example.com",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "host without scheme",
			raw:     "example.com",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "host with trailing slash",
			raw:     "https://example.com/",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "host with path",
			raw:     "https://example.com/path",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "host with port",
			raw:     "https://example.com:8080",
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			raw:     "not a valid url",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty host",
			raw:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "host with spaces",
			raw:     "  example.com  ",
			want:    "example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanHost(tt.raw)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
