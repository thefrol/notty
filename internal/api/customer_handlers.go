package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/thefrol/notty/internal/api/decode"
	"gitlab.com/thefrol/notty/internal/api/respond"
	"gitlab.com/thefrol/notty/internal/api/validate"
)

// CreateClient implements generated.ServerInterface.
func (a *Api) CreateClient(w http.ResponseWriter, r *http.Request) {
	c, err := decode.Customer(r)
	if err != nil {
		respond.BadRequest(w, "%v", err) // может тут оставить место только для ошибки?
		return
	}

	// todo во эта часть -логика, которая должна быть более высокого уровня, в App
	// func (a App) CreateCustomer(customer) Customer

	// если айдишник не указан - создадим сами
	if c.Id == "" {
		c.Id = uuid.New().String()
	}

	res, err := a.app.CustomerRepository.Create(c) // todo а что если такой клиент существует??
	if err != nil {
		http.Error(w, "не удалось создать клиента"+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
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

	c, err := a.app.CustomerRepository.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	err := a.app.CustomerRepository.Delete(id)
	if err != nil {
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

	res, err := a.app.CustomerRepository.Update(c) // todo а что если такой клиент существует??
	if err != nil {
		http.Error(w, "не удалось обновить клиента"+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
		return
	}

	respond.Customer(w, res)
}

// CustomerStats implements generated.ServerInterface.
func (a *Api) CustomerStats(w http.ResponseWriter, r *http.Request, id string) {
	a.StatsByCustomerId(w, r, id)
}
