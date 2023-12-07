package http_utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func buildJSONResponse(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("content-type", "application/json")
	return w
}

func marshallInterface(data interface{}) ([]byte, error) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return respBytes, nil
}

func BuildJSONResponseWithBody(w http.ResponseWriter, body interface{}, httpStatus int) (int, error) {
	w = buildJSONResponse(w)
	w.WriteHeader(httpStatus)
	respBytes, err := marshallInterface(body)
	if err != nil {
		return 0, err
	}
	bytesWritten, err := w.Write(respBytes)

	return bytesWritten, nil
}
