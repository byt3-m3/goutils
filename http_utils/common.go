package http_utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

}

func setJSONHalHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/hal+json")

}

func marshallInterface(data interface{}) ([]byte, error) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return respBytes, nil
}

// WriteJSONFromAny will write the desired struct to the response writer and set the content-type header to application/json
func WriteJSONFromAny(w http.ResponseWriter, body interface{}, httpStatus int) (int, error) {
	setJSONHeader(w)
	w.WriteHeader(httpStatus)
	respBytes, err := marshallInterface(body)
	if err != nil {
		return 0, err
	}
	bytesWritten, err := w.Write(respBytes)

	return bytesWritten, nil
}

// WriteJSONFromBytes will write the byte slice to the response writer and set the content-type header to application/json
func WriteJSONFromBytes(w http.ResponseWriter, respBytes []byte, httpStatus int) (int, error) {
	setJSONHeader(w)
	w.WriteHeader(httpStatus)

	bytesWritten, err := w.Write(respBytes)
	if err != nil {
		return 0, err
	}

	return bytesWritten, nil
}

// WriteJSONHalFromBytes will write the byte slice to the response writer and set the content-type header to application/hal+json
func WriteJSONHalFromBytes(w http.ResponseWriter, respBytes []byte, httpStatus int) (int, error) {
	setJSONHalHeader(w)
	w.WriteHeader(httpStatus)

	bytesWritten, err := w.Write(respBytes)
	if err != nil {
		return 0, err
	}

	return bytesWritten, nil
}
