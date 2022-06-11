package response

import (
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/json"
)

var contentType = "application/json"

type Response struct {
	rw         http.ResponseWriter
	r          *http.Request
	statusCode int
	headers    map[string]string
	body       []byte
}

type ErrorResponse struct {
	Err string `json:"error"`
}

type ValidationErrResponse struct {
	Err []*teal.ValidationError `json:"errors"`
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

func Error(rw http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	res := New(rw, r)
	res.statusCode = http.StatusBadRequest

	res.body, err = json.ToJSON(&ErrorResponse{
		Err: err.Error(),
	})
	if err != nil {
		// TODO log marshal error
		res.body = []byte("")
	}
	res.Write()
}

func ValidationError(rw http.ResponseWriter, r *http.Request, verrs []*teal.ValidationError) {
	res := New(rw, r)
	res.statusCode = http.StatusBadRequest

	var err error
	res.body, err = json.ToJSON(&ValidationErrResponse{
		Err: verrs,
	})
	if err != nil {
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
