// Тут лежат помошники в размаршаливании запросов.
package decode

import (
	"fmt"
	"net/http"

	"github.com/mailru/easyjson"
	"gitlab.com/thefrol/notty/internal/entity"
)

// Customer размаршаливает и валидирует клиента
func Customer(r *http.Request) (entity.Customer, error) {
	c := entity.Customer{}

	err := easyjson.UnmarshalFromReader(r.Body, &c)
	if err != nil {
		return entity.Customer{}, fmt.Errorf("не удалось размаршалить запрос %w", err)
	}
	defer r.Body.Close()

	return c, nil
}

func Subscription(r *http.Request) (entity.Subscription, error) {
	c := entity.Subscription{}

	err := easyjson.UnmarshalFromReader(r.Body, &c)
	if err != nil {
		return entity.Subscription{}, fmt.Errorf("не удалось размаршалить рассылку из запроса %w", err)
	}
	defer r.Body.Close()

	return c, nil
}
