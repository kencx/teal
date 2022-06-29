package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kencx/teal/http/response"
)

func HandleInt64(key string, rw http.ResponseWriter, r *http.Request) int64 {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars[key])
	if err != nil {
		response.BadRequest(rw, r, fmt.Errorf("unable to process id: %v", err))
		return -1
	}
	return int64(id)
}

func HandleString(key string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[key]
}
