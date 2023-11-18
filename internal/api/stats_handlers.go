package api

import (
	"errors"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/respond"
	"gitlab.com/thefrol/notty/internal/app"
)

// FullStats implements generated.ServerInterface.
func (a *Api) FullStats(w http.ResponseWriter, r *http.Request) {
	stats, err := a.app.FullStats()
	if err != nil {
		respond.InternalServerError(w, "Unknown error %v", err)
	}
	stats.ToResponseWriter(w) // чет не оч todo
}

// StatsBySubscriptionId implements generated.ServerInterface.
func (a *Api) StatsBySubscriptionId(w http.ResponseWriter, r *http.Request, id string) {
	stats, err := a.app.StatsBySubscription(id)
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
func (a *Api) StatsByCustomerId(w http.ResponseWriter, r *http.Request, id string) {
	stats, err := a.app.StatsByClient(id)
	if err != nil {
		if errors.Is(err, app.ErrorClientNotFound) {
			respond.NotFound(w, "клиент с id %s не найден: %v", id, err)
			return
		}
		respond.InternalServerError(w, "Unknown error %v", err)
	}

	stats.ToResponseWriter(w) // чет не оч todo
}
