package handler

import "net/http"

// CheckAlive is a handler for checkalive.
func CheckAlive() Handler {
	return func(w ResultWriter, _ *http.Request) error {
		w.Body().Set("result", "success")
		return nil
	}
}
