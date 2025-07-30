package http_utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

type TestStruct struct {
	Name string
}

type mockResponseWriter struct {
	header            http.Header
	writtenBytes      []byte
	writtenStatusCode int
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		header: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write(bytes []byte) (int, error) {
	m.writtenBytes = bytes
	return len(bytes), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.writtenStatusCode = statusCode
}

func TestWriteJSONFromAny(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "valid data",
			input:      TestStruct{Name: "test"},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "nil data",
			input:      nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newMockResponseWriter()
			count, err := WriteJSONFromAny(w, tt.input, tt.wantStatus)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Greater(t, count, 0)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.wantStatus, w.writtenStatusCode)

			// Verify the written JSON is valid
			var decoded interface{}
			assert.NoError(t, json.Unmarshal(w.writtenBytes, &decoded))
		})
	}
}

func TestWriteJSONFromBytes(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "valid JSON bytes",
			input:      []byte(`{"Name":"test"}`),
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty bytes",
			input:      []byte{},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newMockResponseWriter()
			count, err := WriteJSONFromBytes(w, tt.input, tt.wantStatus)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.input), count)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.wantStatus, w.writtenStatusCode)
			assert.Equal(t, tt.input, w.writtenBytes)
		})
	}
}

func TestJSONDecode(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     TestStruct
		wantErr  bool
	}{
		{
			name:     "valid json",
			jsonData: `{"Name": "test name"}`,
			want:     TestStruct{Name: "test name"},
			wantErr:  false,
		},
		{
			name:     "invalid json",
			jsonData: `{"Name": "test name"`,
			want:     TestStruct{},
			wantErr:  true,
		},
		{
			name:     "empty json object",
			jsonData: `{}`,
			want:     TestStruct{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "http://example.com", strings.NewReader(tt.jsonData))
			var got TestStruct
			err := JSONDecode(req, &got)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJSONEncode(t *testing.T) {
	tests := []struct {
		name    string
		input   TestStruct
		want    string
		wantErr bool
	}{
		{
			name:    "valid struct",
			input:   TestStruct{Name: "test name"},
			want:    `{"Name":"test name"}` + "\n",
			wantErr: false,
		},
		{
			name:    "empty struct",
			input:   TestStruct{},
			want:    `{"Name":""}` + "\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newMockResponseWriter()
			result, err := JSONEncode(w, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			// Check the actual written bytes from the mock writer
			assert.Equal(t, tt.want, string(w.writtenBytes))
			assert.Equal(t, tt.input, result)
		})
	}
}
