package http

import (
	"io"
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

func (s Server) MiddlewareBookValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		book := teal.Book{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(rw, r, err)
			return
		}
		err = util.FromJSON(body, &book)
		if err != nil {
			response.Error(rw, r, err)
			return
		}

		verrs := book.Validate()
		if len(verrs) > 0 {
			// s.ErrLog.Println(verrs)
			response.ValidationError(rw, r, verrs)
			return
		}

		// add book to context
		// ctx := context.WithValue(r.Context(), KeyBook{}, book)
		// r = r.WithContext(ctx)

		// call next handler
		next.ServeHTTP(rw, r)
	})
}
