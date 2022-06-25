package response

import (
	"net/http"

	"github.com/kencx/teal/util"
)

var contentType = "application/json"

type Envelope map[string]interface{}

type response struct {
	rw         http.ResponseWriter
	r          *http.Request
	statusCode int
	headers    map[string]string
	body       []byte
}

func New(rw http.ResponseWriter, r *http.Request) *response {
	return &response{
		rw:         rw,
		r:          r,
		statusCode: http.StatusOK,
		headers:    map[string]string{"Content-Type": contentType},
	}
}

func OK(rw http.ResponseWriter, r *http.Request, body []byte) {
	res := New(rw, r)
	res.statusCode = http.StatusOK
	res.body = body
	res.Write()
}

func NoContent(rw http.ResponseWriter, r *http.Request) {
	res := New(rw, r)
	res.statusCode = http.StatusNoContent
	res.Write()
}

func Created(rw http.ResponseWriter, r *http.Request, body []byte) {
	res := New(rw, r)
	res.statusCode = http.StatusCreated
	res.body = body
	res.Write()
}

func NewError(rw http.ResponseWriter, r *http.Request, err interface{}) *response {
	res := New(rw, r)
	res.statusCode = http.StatusBadRequest

	switch t := err.(type) {
	case error:
		res.body, err = util.ToJSON(Envelope{"error": t.Error()})
		if err != nil {
			res.body, err = util.ToJSON(Envelope{"error": "something went wrong"})
			res.statusCode = http.StatusInternalServerError
		}
	default:
		res.body, err = util.ToJSON(Envelope{"error": err})
		if err != nil {
			res.body, err = util.ToJSON(Envelope{"error": "something went wrong"})
			res.statusCode = http.StatusInternalServerError
		}
	}
	return res
}

func BadRequest(rw http.ResponseWriter, r *http.Request, err error) {
	res := NewError(rw, r, err)
	res.Write()
}

func InternalServerError(rw http.ResponseWriter, r *http.Request, err error) {
	res := NewError(rw, r, err)
	res.statusCode = http.StatusInternalServerError
	res.Write()
}

func NotFound(rw http.ResponseWriter, r *http.Request, err error) {
	res := NewError(rw, r, err)
	res.statusCode = http.StatusNotFound
	res.Write()
}

func Unauthorized(rw http.ResponseWriter, r *http.Request, err error) {
	res := NewError(rw, r, err)
	res.statusCode = http.StatusUnauthorized
	res.headers["WWW-Authenticate"] = `Basic realm="Restricted"`
	res.Write()
}

func ValidationError(rw http.ResponseWriter, r *http.Request, err map[string]string) {
	res := NewError(rw, r, err)
	res.statusCode = http.StatusUnprocessableEntity
	res.Write()
}

// TODO handlePanic

func (r *response) Write() {
	for k, v := range r.headers {
		r.rw.Header().Set(k, v)
	}

	r.rw.WriteHeader(r.statusCode)
	r.rw.Write(r.body)
}
