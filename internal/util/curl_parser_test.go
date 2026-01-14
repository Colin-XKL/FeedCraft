package util

import (
	"reflect"
	"testing"
)

func TestParseCurlCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		want    *CurlParseResult
		wantErr bool
	}{
		{
			name: "Basic GET",
			cmd:  `curl https://example.com`,
			want: &CurlParseResult{
				Method:  "GET",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "",
			},
			wantErr: false,
		},
		{
			name: "Basic GET with explicit method",
			cmd:  `curl -X GET https://example.com`,
			want: &CurlParseResult{
				Method:  "GET",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "",
			},
			wantErr: false,
		},
		{
			name: "Basic POST with JSON",
			cmd:  `curl -X POST https://example.com -d '{"foo":"bar"}'`,
			want: &CurlParseResult{
				Method:  "POST",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    `{"foo":"bar"}`,
			},
			wantErr: false,
		},
		{
			name: "Implicit POST with data",
			cmd:  `curl https://example.com --data "param=value"`,
			want: &CurlParseResult{
				Method:  "POST",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "param=value",
			},
			wantErr: false,
		},
		{
			name: "Headers",
			cmd:  `curl -H "Content-Type: application/json" -H "Authorization: Bearer token" https://api.example.com`,
			want: &CurlParseResult{
				Method: "GET",
				URL:    "https://api.example.com",
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer token",
				},
				Body: "",
			},
			wantErr: false,
		},
		{
			name: "Compressed flag",
			cmd:  `curl --compressed https://example.com`,
			want: &CurlParseResult{
				Method:  "GET",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "",
			},
			wantErr: false,
		},
		{
			name: "Attached flags",
			cmd:  `curl -XPOST https://example.com`,
			want: &CurlParseResult{
				Method:  "POST",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "",
			},
			wantErr: false,
		},
		{
			name: "Attached Header",
			cmd:  `curl -H"Accept: json" https://example.com`,
			want: &CurlParseResult{
				Method: "GET",
				URL:    "https://example.com",
				Headers: map[string]string{
					"Accept": "json",
				},
				Body: "",
			},
			wantErr: false,
		},
		{
			name: "Basic Auth",
			cmd:  `curl -u user:pass https://example.com`,
			want: &CurlParseResult{
				Method: "GET",
				URL:    "https://example.com",
				Headers: map[string]string{
					"Authorization": "Basic dXNlcjpwYXNz", // base64(user:pass)
				},
				Body: "",
			},
			wantErr: false,
		},
		{
			name: "Multiple Data Flags",
			cmd:  `curl -d "p1=v1" --data "p2=v2" https://example.com`,
			want: &CurlParseResult{
				Method:  "POST",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "p1=v1&p2=v2",
			},
			wantErr: false,
		},
		{
			name: "Explicit URL flag",
			cmd:  `curl --url https://example.com`,
			want: &CurlParseResult{
				Method:  "GET",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    "",
			},
			wantErr: false,
		},
		{
			name: "Data Raw",
			cmd:  `curl --data-raw '{"a":1}' https://example.com`,
			want: &CurlParseResult{
				Method:  "POST",
				URL:     "https://example.com",
				Headers: map[string]string{},
				Body:    `{"a":1}`,
			},
			wantErr: false,
		},
		{
			name: "Invalid Command",
			cmd:  `wget https://example.com`,
			want: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCurlCommand(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCurlCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCurlCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
