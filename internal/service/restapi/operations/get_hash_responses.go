// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetHashFoundCode is the HTTP code returned for type GetHashFound
const GetHashFoundCode int = 302

/*GetHashFound Redirection to long URL.

swagger:response getHashFound
*/
type GetHashFound struct {
}

// NewGetHashFound creates GetHashFound with default headers values
func NewGetHashFound() *GetHashFound {

	return &GetHashFound{}
}

// WriteResponse to the client
func (o *GetHashFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(302)
}

// GetHashNotFoundCode is the HTTP code returned for type GetHashNotFound
const GetHashNotFoundCode int = 404

/*GetHashNotFound Shortened URL not found.

swagger:response getHashNotFound
*/
type GetHashNotFound struct {
}

// NewGetHashNotFound creates GetHashNotFound with default headers values
func NewGetHashNotFound() *GetHashNotFound {

	return &GetHashNotFound{}
}

// WriteResponse to the client
func (o *GetHashNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}
