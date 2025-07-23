package http_utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// SetJSONHeader Sets the JSON header
func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

}

func marshallInterface(data interface{}) ([]byte, error) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("error marshalling interface",
			slog.Any("data", data),
		)
		return nil, err
	}
	return respBytes, nil
}

// WriteJSONFromAny will write the desired struct to the response writer and set the content-type header to application/json
func WriteJSONFromAny(w http.ResponseWriter, body interface{}, httpStatus int) (int, error) {
	SetJSONHeader(w)
	w.WriteHeader(httpStatus)
	respBytes, err := marshallInterface(body)
	if err != nil {
		return 0, err
	}
	bytesWritten, err := w.Write(respBytes)

	return bytesWritten, nil
}

// MustWriteJSONFromAny will write the desired struct to the response writer and set the content-type header to application/json
func MustWriteJSONFromAny(w http.ResponseWriter, body interface{}, httpStatus int) {
	_, err := WriteJSONFromAny(w, body, httpStatus)
	if err != nil {
		panic(err)
	}

}

// WriteJSONFromBytes will write the byte slice to the response writer and set the content-type header to application/json
func WriteJSONFromBytes(w http.ResponseWriter, respBytes []byte, httpStatus int) (int, error) {
	SetJSONHeader(w)
	w.WriteHeader(httpStatus)

	bytesWritten, err := w.Write(respBytes)
	if err != nil {
		return 0, err
	}

	return bytesWritten, nil
}

// MustWriteJSONFromBytes will write the byte slice to the response writer and set the content-type header to application/json
func MustWriteJSONFromBytes(w http.ResponseWriter, respBytes []byte, httpStatus int) {
	_, err := WriteJSONFromBytes(w, respBytes, httpStatus)
	if err != nil {
		panic(err)
	}
}

// JSONDecode uses generics to decode the data into the provided type
func JSONDecode[T any](req *http.Response, v T) (T, error) {
	var data T
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return v, err
	}

	return data, nil

}

// JSONEncode uses generics to encode the data into the provided type
func JSONEncode[T any](w http.ResponseWriter, v T) (T, error) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return v, err
	}

	return v, nil
}
