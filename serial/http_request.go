package serial

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Define the struct that captures all relevant fields from the http.Request
type RequestData struct {
	Method        string              `json:"method"`
	URL           string              `json:"url"`
	Proto         string              `json:"proto"`
	Header        map[string][]string `json:"header"`
	Body          string              `json:"body"`
	ContentLength int64               `json:"content_length"`
	Host          string              `json:"host"`
	RemoteAddr    string              `json:"remote_addr"`
}

// Convert an http.Request to a RequestData instance
func requestToJSONStruct(req *http.Request) (*RequestData, error) {
	// Read the body (assuming it's small and can fit in memory)
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		// Restore the io.ReadCloser to its original state
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	requestData := &RequestData{
		Method:        req.Method,
		URL:           req.URL.String(),
		Proto:         req.Proto,
		Header:        req.Header,
		Body:          string(bodyBytes),
		ContentLength: req.ContentLength,
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
	}

	return requestData, nil
}

// Encode the http.Request to a JSON byte slice
func EncodeRequestToJSON(req *http.Request) ([]byte, error) {
	requestData, err := requestToJSONStruct(req)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(requestData, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
