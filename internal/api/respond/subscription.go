// Упрощенные ответы для хендлеров
package respond

import (
	"net/http"

	"github.com/mailru/easyjson"
	"gitlab.com/thefrol/notty/internal/entity"
)

// Subscription в любом случае отвечает рассылкой или сообщает об ошибке
func Subscription(w http.ResponseWriter, c entity.Subscription) {
	_, _, err := easyjson.MarshalToHTTPResponseWriter(c, w)
	if err != nil {
		InternalServerError(w, "не удалось демаршалить рассылку %v", err)
		return
	}
}
