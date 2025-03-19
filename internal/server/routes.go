package server

import "net/http"

/*
Setup and initialize the routes for the server
*/
func (s *Server) SetupRoutes(handler *Handler) {
	s.handler.HandleFunc("POST /counter", handler.handleCounter)

	s.handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
