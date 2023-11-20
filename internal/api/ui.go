package api

import (
	"net/http"

	"gitlab.com/thefrol/notty/internal/api/generated"
	"gitlab.com/thefrol/notty/pkg/swagger"
)

// ErrorEndpoint создает такую ручку, в которой все время при обращении
// будет выводиться ошибка. Это полезно, когда что-то не удалось запустить
// но мы заранее знаем, что по этому адресу при этом запуске сервера
// уже ничего не достать. При этого тут что-то может быть
func ErrorEndpoint(msg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
	}
}

// Swagger возвращает ручку, на которой держится сваггер текущей спецификации
// если доку не удалось создать, то там будет выводиться ошибка, но всегда
// эта ручка существует
func (a *Server) Swagger() http.HandlerFunc {
	sw, err := generated.GetSwagger()
	if err != nil {
		return ErrorEndpoint("Сваггер интерфейс не получается запустить, не возможно пропарсить спеку  вставленну в бинарник")
	}
	return swagger.Handler(sw, sw.Info.Title)
}
