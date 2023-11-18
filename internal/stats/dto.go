// Этот пакет отвечает за сбор статистики, содержит репозиторий, сервис
// и обьект данных
//
// В этом модуле я просто пробую совершенно другую семантику работы с данными.
// все, что касается одной какой-то сущности или данных будет лежать в одной папке
package stats

import (
	"encoding/json"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/respond"
)

type Statistics map[string]int

func (s Statistics) ToResponseWriter(w http.ResponseWriter) {
	// не используем easyjson, потому что он не умеет в мапы

	// todo а что если статистику кешировать в редис?
	err := json.NewEncoder(w).Encode(&s)
	if err != nil {
		respond.InternalServerError(w, "cant Marshall Stats to json: %v", err) // todo MarshallingError
	}
} //todo нормально я сюда зависимость впилил?)
