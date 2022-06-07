package http

import "net/http"

func (s *Server) GetAllAuthors(rw http.ResponseWriter, r *http.Request) {
	b, err := s.Authors.GetAll()
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	res, err := ToJSON(b)
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if len(res) == 0 {
		rw.Write([]byte("No books added"))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
}
