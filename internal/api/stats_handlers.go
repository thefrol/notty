package api

import (
	"errors"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/respond"
	"gitlab.com/thefrol/notty/internal/app"
)

// FullStats implements generated.ServerInterface.
func (a *Server) FullStats(w http.ResponseWriter, r *http.Request) {
	stats, err := a.app.FullStats(r.Context())
	if err != nil {
		respond.InternalServerError(w, "Unknown error %v", err)
	}
	stats.ToResponseWriter(w) // чет не оч todo
}

// StatsBySubscriptionId implements generated.ServerInterface.
func (a *Server) StatsBySubscriptionId(w http.ResponseWriter, r *http.Request, id string) {
	stats, err := a.app.StatsBySubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrorSubscriptionNotFound) {
			respond.NotFound(w, "рассылка с id %s не найдена: %v", id, err)
			return
		}
		respond.InternalServerError(w, "Unknown error %v", err)
	}

	stats.ToResponseWriter(w) // чет не оч todo
}

// StatsBySubscriptionId implements generated.ServerInterface.
func (a *Server) StatsByCustomerId(w http.ResponseWriter, r *http.Request, id string) {
	stats, err := a.app.StatsByClient(r.Context(), id)
	if err != nil {
		if errors.Is(err, app.ErrorCustomerNotFound) {
			respond.NotFound(w, "клиент с id %s не найден: %v", id, err)
			return
		}
		respond.InternalServerError(w, "Unknown error %v", err)
	}

	stats.ToResponseWriter(w) // чет не оч todo
}
