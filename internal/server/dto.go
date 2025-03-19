package server

// DTO (Data-Transfer-Object) for counter request
// In a production environment, it could be enhanced using the validator package.
type CounterRequest struct {
	Directory string `json:"directory"`
	Word      string `json:"word"`
}

// DTO (Data-Transfer-Object) for counter response
type CounterResponse struct {
	Count int `json:"count"`
}

// DTO (Data-Transfer-Object) for error response
type ErrorResponse struct {
	Error string `json:"error"`
}
