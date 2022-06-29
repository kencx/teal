package http

import (
	"net/http"

	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

type health struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func NewHealth(version, env string) *health {
	return &health{
		Version:     version,
		Environment: env,
	}
}

func (s *Server) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	h := NewHealth("1.0", "dev")

	res, err := util.ToJSON(response.Envelope{"healthcheck": h})
	if err != nil {
		response.InternalServerError(rw, r, err)
	}

	response.OK(rw, r, res)
}
