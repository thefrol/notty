package dto

import (
	"encoding/json"
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/respond"
)

type Statistics map[string]int

func (s Statistics) ToResponseWriter(w http.ResponseWriter) { // todo не нравится!
	// не используем easyjson, потому что он не умеет в мапы

	// todo а что если статистику кешировать в редис?
	err := json.NewEncoder(w).Encode(&s)
	if err != nil {
		respond.InternalServerError(w, "cant Marshall Stats to json: %v", err) // todo MarshallingError
	}
} //todo нормально я сюда зависимость впилил?)
