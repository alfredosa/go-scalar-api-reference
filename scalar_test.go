package scalar

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGetScalarHTMLContent(t *testing.T) {
	tests := []struct {
		name        string
		options     *Options
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "with nil spec content and empty spec URL",
			options: &Options{
				SpecContent: nil,
				SpecURL:     "",
			},
			want:        "",
			wantErr:     true,
			errContains: "specURL or specContent must be provided",
		},
		{
			name: "with map spec content",
			options: &Options{
				SpecContent: map[string]interface{}{
					"test": "value",
				},
			},
			want:    `{"test":"value"}`,
			wantErr: false,
		},
		{
			name: "with function spec content",
			options: &Options{
				SpecContent: func() map[string]interface{} {
					return map[string]interface{}{
						"test": "value",
					}
				},
			},
			want:    `{"test":"value"}`,
			wantErr: false,
		},
		{
			name: "with string spec content",
			options: &Options{
				SpecContent: "test content",
			},
			want:    "test content",
			wantErr: false,
		},
		{
			name: "with invalid spec content type",
			options: &Options{
				SpecContent: 123, // integer is not a valid type
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetScalarHTMLContent(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScalarHTMLContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("GetScalarHTMLContent() error = %v, want error containing %v", err, tt.errContains)
				return
			}
			if got != tt.want {
				t.Errorf("GetScalarHTMLContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCustomCSS(t *testing.T) {
	tests := []struct {
		name    string
		options *Options
		want    string
	}{
		{
			name: "with theme set",
			options: &Options{
				Theme: "dark",
			},
			want: "",
		},
		{
			name: "with custom CSS",
			options: &Options{
				CustomCss: "body { color: red; }",
			},
			want: "<style>body { color: red; }</style>",
		},
		{
			name:    "with default options",
			options: &Options{},
			want:    "<style>" + CustomThemeCSS + "</style>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCustomCSS(tt.options)
			if got != tt.want {
				t.Errorf("GetCustomCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetScalarCDN(t *testing.T) {
	tests := []struct {
		name    string
		options *Options
		want    string
	}{
		{
			name: "with CDN URL",
			options: &Options{
				CDN: "https://cdn.example.com/scalar.js",
			},
			want: `<script src="https://cdn.example.com/scalar.js"></script>`,
		},
		{
			name:    "with empty CDN",
			options: &Options{},
			want:    `<script src=""></script>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetScalarCDN(tt.options)
			if got != tt.want {
				t.Errorf("GetScalarCDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetScalarScriptWithHTMLContent(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		specContent string
		want        string
	}{
		{
			name:        "basic content",
			config:      "testConfig",
			specContent: "testContent",
			want:        `<script id="api-reference" type="application/json" data-configuration="testConfig">testContent</script>`,
		},
		{
			name:        "empty content",
			config:      "",
			specContent: "",
			want:        `<script id="api-reference" type="application/json" data-configuration=""></script>`,
		},
		{
			name:        "with special characters",
			config:      "config&quot;test",
			specContent: "content<>test",
			want:        `<script id="api-reference" type="application/json" data-configuration="config&quot;test">content<>test</script>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetScalarScriptWithHTMLContent(tt.config, tt.specContent)
			if got != tt.want {
				t.Errorf("GetScalarScriptWithHTMLContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare JSON strings regardless of formatting
func compareJSON(t *testing.T, got, want string) bool {
	var gotMap, wantMap map[string]interface{}

	if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
		t.Errorf("Failed to parse got JSON: %v", err)
		return false
	}

	if err := json.Unmarshal([]byte(want), &wantMap); err != nil {
		t.Errorf("Failed to parse want JSON: %v", err)
		return false
	}

	gotJSON, _ := json.Marshal(gotMap)
	wantJSON, _ := json.Marshal(wantMap)

	return string(gotJSON) == string(wantJSON)
}
