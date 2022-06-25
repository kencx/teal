package http

import (
	"net/http"

	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

type health struct {
	Version string `json:"version"`
}

func NewHealth(version string) *health {
	return &health{version}
}

func (s *Server) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	h := NewHealth("v1.0")

	res, err := util.ToJSON(response.Envelope{"health": h})
	if err != nil {
		response.InternalServerError(rw, r, err)
	}

	response.OK(rw, r, res)
}
