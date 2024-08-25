package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"surge/internal/utilities"
	"time"
)

var defaultTimeout = time.Second * 10

func chooseHost(base, defaultHost string) string {
	if base == "" {
		return "https://" + defaultHost
	}

	baseLen := len(base)
	if base[baseLen-1] == '/' {
		return base[:baseLen-1]
	}

	return base
}

func makeRequest(ctx context.Context, tok *oauth2.Token, g *oauth2.Config, url string, dst interface{}) error {
	client := g.Client(ctx, tok)
	client.Timeout = defaultTimeout
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer utilities.SafeClose(res.Body)

	bodyBytes, _ := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return newHttpError(res.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(res.Body).Decode(dst); err != nil {
		return err
	}

	return nil
}
