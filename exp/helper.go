package main

import (
	"net/http"

	"github.com/gorilla/schema"
)

// really should be in the controllers package
// cus (for now) it's used to keep code for parsing
// forms.

// parseForm takes in a request and a destination
// and parses the form into the destination object.
func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
