package http_utils

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

type TestData struct {
	Name string
}

type mockResponseWriter struct {
	header              http.Header
	WriteMockFunc       func(bytes []byte) (int, error)
	WriteHeaderMockFunc func(statusCode int)
}

func (m mockResponseWriter) Header() http.Header {
	return m.header
}

func (m mockResponseWriter) Write(bytes []byte) (int, error) {
	return m.WriteMockFunc(bytes)
}

func (m mockResponseWriter) WriteHeader(statusCode int) {
	m.WriteHeaderMockFunc(statusCode)
}

func TestWriteJSONFromAny(t *testing.T) {
	mockWriter := mockResponseWriter{
		header: http.Header{},
		WriteMockFunc: func(bytes []byte) (int, error) {
			return len(bytes), nil
		},
		WriteHeaderMockFunc: func(statusCode int) {
			slog.Info("status received")

		},
	}

	data := TestData{Name: "test"}

	count, err := WriteJSONFromAny(mockWriter, &data, http.StatusOK)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Greater(t, count, 0)
	assert.Equal(t, mockWriter.Header().Get("Content-Type"), "application/json")

}

func TestWriteJSONFromBytes(t *testing.T) {

	mockWriter := mockResponseWriter{
		header: http.Header{},
		WriteMockFunc: func(bytes []byte) (int, error) {
			return len(bytes), nil
		},
		WriteHeaderMockFunc: func(statusCode int) {
			slog.Info("status received")

		},
	}

	data := TestData{Name: "test"}

	dataBytes, _ := json.Marshal(&data)

	count, err := WriteJSONFromBytes(mockWriter, dataBytes, http.StatusOK)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Greater(t, count, 0)
	assert.Equal(t, mockWriter.Header().Get("Content-Type"), "application/json")
}

type TestStruct struct {
	Name string
}

func TestJSONDecode(t *testing.T) {
	// Test cases
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
			jsonData: `{"Name": "test name"`, // missing closing brace
			want:     TestStruct{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new response with a ReadCloser body
			resp := &http.Response{
				Body: io.NopCloser(strings.NewReader(tt.jsonData)),
			}

			// Create an empty struct to decode into
			var got TestStruct

			// Call JSONDecode
			result, err := JSONDecode(resp, got)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check result
			assert.Equal(t, tt.want, result)
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
			want:    "{\"Name\":\"test name\"}\n", // Note: json.Encoder adds a newline
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture the JSON output
			var buf bytes.Buffer
			mockWriter := mockResponseWriter{
				header: http.Header{},
				WriteMockFunc: func(b []byte) (int, error) {
					return buf.Write(b)
				},
				WriteHeaderMockFunc: func(statusCode int) {},
			}

			// Call JSONEncode
			result, err := JSONEncode(mockWriter, tt.input)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check the encoded output
			assert.Equal(t, tt.want, buf.String())

			// Check that the returned value matches the input
			assert.Equal(t, tt.input, result)
		})
	}
}
