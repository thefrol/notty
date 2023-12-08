package api

import (
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/respond"
)

// GetMessage implements generated.ServerInterface.
func (a *Server) GetMessage(w http.ResponseWriter, r *http.Request, id string) {
	respond.InternalServerError(w, "Не реализовано :(")
}
