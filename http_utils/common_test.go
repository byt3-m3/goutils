package http_utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"testing"
)

var ()

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

func TestWriteJSONHalFromBytes(t *testing.T) {
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

	count, err := WriteJSONHalFromBytes(mockWriter, dataBytes, http.StatusOK)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Greater(t, count, 0)
	assert.Equal(t, mockWriter.Header().Get("Content-Type"), "application/hal+json")
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
