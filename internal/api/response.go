package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

func writeResponseJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func writeResponsePrettyJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}
