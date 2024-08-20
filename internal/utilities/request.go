package utilities

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// GetBodyContentAsBytes reads the whole request body properly into a byte array.
func GetBodyContentAsBytes(req *http.Request) ([]byte, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return nil, nil
	}

	originalBody := req.Body
	defer func() {
		if err := originalBody.Close(); err != nil {
			logrus.WithError(err).Warn("Close operation failed")
		}
	}()

	buf, err := io.ReadAll(originalBody)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(buf))

	return buf, nil
}

func GetBodyJson[T any](req *http.Request) (*T, error) {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	var value T
	err := decoder.Decode(&value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
