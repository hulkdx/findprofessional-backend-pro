package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ReadJSON(r *http.Request, out interface{}) error {
	if r.Body == nil || r.ContentLength == 0 {
		return nil
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(out)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	if dec.More() {
		return errors.New("unexpected content after JSON")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

func WriteError(w http.ResponseWriter, code int, error string) {
	data := map[string]any{
		"error": error,
	}
	WriteJSON(w, code, data)
}

func WriteGeneralError(w http.ResponseWriter, error error) {
	WriteError(w, http.StatusInternalServerError, error.Error())
}
