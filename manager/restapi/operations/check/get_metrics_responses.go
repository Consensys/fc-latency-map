// Code generated by go-swagger; DO NOT EDIT.

package check

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ConsenSys/fc-latency-map/manager/models"
)

// GetMetricsOKCode is the HTTP code returned for type GetMetricsOK
const GetMetricsOKCode int = 200

/*GetMetricsOK Manager metrics sent

swagger:response getMetricsOK
*/
type GetMetricsOK struct {

	/*
	  In: Body
	*/
	Payload *models.Metrics `json:"body,omitempty"`
}

// NewGetMetricsOK creates GetMetricsOK with default headers values
func NewGetMetricsOK() *GetMetricsOK {

	return &GetMetricsOK{}
}

// WithPayload adds the payload to the get metrics o k response
func (o *GetMetricsOK) WithPayload(payload *models.Metrics) *GetMetricsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get metrics o k response
func (o *GetMetricsOK) SetPayload(payload *models.Metrics) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMetricsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetMetricsDefault Internal error

swagger:response getMetricsDefault
*/
type GetMetricsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetMetricsDefault creates GetMetricsDefault with default headers values
func NewGetMetricsDefault(code int) *GetMetricsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetMetricsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get metrics default response
func (o *GetMetricsDefault) WithStatusCode(code int) *GetMetricsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get metrics default response
func (o *GetMetricsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get metrics default response
func (o *GetMetricsDefault) WithPayload(payload *models.Error) *GetMetricsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get metrics default response
func (o *GetMetricsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMetricsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
