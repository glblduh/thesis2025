package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func encodeRes(w http.ResponseWriter, v any) error {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		errorRes(w, "JSON Encoder error", http.StatusInternalServerError)
	}
	return err
}

func errorRes(w http.ResponseWriter, errorResponse string, code int) {
	w.WriteHeader(code)
	err := encodeRes(w, &jsonErrorRes{
		Error: errorResponse,
	})
	if err != nil {
		w.Write([]byte(errorResponse))
	}
}

func decodeBody(w http.ResponseWriter, body io.Reader, v any) error {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		errorRes(w, "JSON Encoder error", http.StatusInternalServerError)
	}
	return err
}
