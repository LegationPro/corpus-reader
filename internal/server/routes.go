package server

/*
Setup and initialize the routes for the server
*/
func (s *Server) SetupRoutes(handler *Handler) {
	s.handler.HandleFunc("POST /counter", handler.HandleCounter)
}
