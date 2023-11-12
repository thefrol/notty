package api

import (
	"log"
	"net/http"

	"gitlab.com/tanna.dev/openapi-doc-http-handler/elements"
	"gitlab.com/thefrol/notty/internal/api/generated"
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

// Docs создает ручку для сваггера, которую можно прицепить к роутеру
func Docs() http.HandlerFunc {
	h, err := elements.NewHandler(generated.GetSwagger())
	if err != nil {
		log.Println("Cant handle swagger")
		return ErrorEndpoint("Это правильный адрес, но сваггер страничку не удалось создать. Не пропарсилась спека.")
	}
	return h
}
