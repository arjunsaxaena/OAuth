package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"reflect"

	"oauth/Go/common/validation"
)


type JsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}


// ReadJsonAndValidate reads JSON from the request body and validates it.
func ReadJsonAndValidate (w http.ResponseWriter, r *http.Request, data any) error{
	if err := ReadJsonFromBody(w, r, data); err != nil {
	    return err
	}

	if err := validation.ValidateRequest(data); err != nil {
		return err
	}

	return nil
}


func ReadJsonFromBody(w http.ResponseWriter, r *http.Request, data any) error{
	maxBytes := 10 << 20 //  = 10 * 2^20 = 10 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) // Wraps the request body so that clients can’t send more than maxBytes

	dec := json.NewDecoder(r.Body) // Create a JSON decoder for the request body
	err := dec.Decode(data) // Decode JSON into your struct
	if err != nil {
		log.Println(err)
		log.Println(reflect.TypeOf(err))
		return err
	}
		
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value") // If the body has extra JSON objects, the second decode will succeed instead of hitting EOF (end of file)
	}

	return nil
}
