package respond

import (
	"net/http"

	"github.com/mailru/easyjson"
	"gitlab.com/thefrol/notty/internal/entity"
)

// Customer в любом случае отвечает клиентом или сообщает об ошибке
func Customer(w http.ResponseWriter, c entity.Customer) {
	_, _, err := easyjson.MarshalToHTTPResponseWriter(c, w)
	if err != nil {
		InternalServerError(w, "не удалось демаршалить клиента %v", err)
		return
	}
}
