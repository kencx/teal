package response

import (
	"fmt"
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/util"
)

var contentType = "application/json"

type Response struct {
	rw         http.ResponseWriter
	r          *http.Request
	statusCode int
	headers    map[string]string
	body       []byte
}

func New(rw http.ResponseWriter, r *http.Request) *Response {
	return &Response{
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

type ErrorResponse struct {
	Err string `json:"error"`
}

func NewError(rw http.ResponseWriter, r *http.Request, err error) *Response {
	res := New(rw, r)
	res.statusCode = http.StatusBadRequest

	res.body, err = util.ToJSON(&ErrorResponse{
		Err: err.Error(),
	})
	if err != nil {
		// TODO log marshal error
		res.body = []byte("")
		res.statusCode = http.StatusInternalServerError
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

func Unauthorized(rw http.ResponseWriter, r *http.Request, body []byte) {
	res := NewError(rw, r, fmt.Errorf(string(body)))
	res.statusCode = http.StatusUnauthorized
	res.headers["WWW-Authenticate"] = `Basic realm="Restricted"`
	res.Write()
}

// TODO handlePanic

type ValidationErrResponse struct {
	Err []*teal.ValidationError `json:"errors"`
}

func ValidationError(rw http.ResponseWriter, r *http.Request, verrs []*teal.ValidationError) {
	res := New(rw, r)
	res.statusCode = http.StatusBadRequest

	var err error
	res.body, err = util.ToJSON(&ValidationErrResponse{
		Err: verrs,
	})
	if err != nil {
		res.statusCode = http.StatusInternalServerError
		res.body = []byte("")
	}
	res.Write()
}

func (r *Response) Write() {
	for k, v := range r.headers {
		r.rw.Header().Set(k, v)
	}

	r.rw.WriteHeader(r.statusCode)
	r.rw.Write(r.body)
}
