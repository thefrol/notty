package api

import (
	"errors"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/decode"
	"gitlab.com/thefrol/notty/internal/api/respond"
	"gitlab.com/thefrol/notty/internal/api/validate"
	"gitlab.com/thefrol/notty/internal/app"
)

// CreateClient implements generated.ServerInterface.
func (a *Api) CreateClient(w http.ResponseWriter, r *http.Request) {
	c, err := decode.Customer(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	res, err := a.app.Customers.Create(c) // todo а что если такой клиент существует??
	if err != nil {
		if errors.Is(err, app.ErrorCustomerExists) {
			respond.Errorf(w, http.StatusConflict, "Клинт уже с id %s существует ", c.Id)
		}
		respond.InternalServerError(w, "Неизвестная ошибка %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respond.Customer(w, res)
}

// GetClient implements generated.ServerInterface.
func (a *Api) GetClient(w http.ResponseWriter, r *http.Request, id string) {
	if err := validate.Id(id); err != nil {
		respond.BadRequest(w, "%v", err)
		return
	}

	c, err := a.app.Customers.Get(id)
	if err != nil {
		if errors.Is(err, app.ErrorCustomerNotFound) {
			respond.NotFound(w, "Клиент с id %s не обнаружен", id)
			return
		}
		respond.InternalServerError(w, "Не удалось найти клиента с id=%s %v", id, err)
		return
	}
	respond.Customer(w, c)
}

// DeleteClient implements generated.ServerInterface.
func (a *Api) DeleteClient(w http.ResponseWriter, r *http.Request, id string) {
	if err := validate.Id(id); err != nil {
		respond.BadRequest(w, "%v", err)
		return
	}

	err := a.app.Customers.Delete(id)
	if err != nil {
		if errors.Is(err, app.ErrorCustomerNotFound) {
			respond.NotFound(w, "Клиент с id %s не обнаружен", id)
			return
		}
		respond.InternalServerError(w, "Не удалось удалить клиента %v", err)
		return
	}
}

// UpdateClient implements generated.ServerInterface.
func (a *Api) UpdateClient(w http.ResponseWriter, r *http.Request, id string) {
	c, err := decode.Customer(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	c.Id = id // заменяем айдишник на тот, что стоит в запросе

	res, err := a.app.Customers.Update(c) // todo а что если такой клиент существует??
	if err != nil {
		if errors.Is(err, app.ErrorCustomerNotFound) {
			respond.NotFound(w, "Клиент с id %s не обнаружен", id)
			return
		}
		http.Error(w, "не удалось обновить клиента"+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
		return
	}

	respond.Customer(w, res)
}

// CustomerStats implements generated.ServerInterface.
func (a *Api) CustomerStats(w http.ResponseWriter, r *http.Request, id string) {
	a.StatsByCustomerId(w, r, id)
}
