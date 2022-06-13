package http

import (
	"net/http"

	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

type health struct {
	Version string
}

func NewHealth(version string) *health {
	return &health{version}
}

func (s *Server) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	h := NewHealth("v1.0")

	res, err := util.ToJSON(h)
	if err != nil {
		response.Error(rw, r, err)
	}

	response.OK(rw, r, res)
}
