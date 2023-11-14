package postman

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"gitlab.com/thefrol/notty/internal/entity"
)

type Poster struct {
	EndPoint string
	Token    string
	client   *resty.Client
}

func New(endpoint string, retryWaitSeconds int, retryCount int, token string) Poster {
	return Poster{
		EndPoint: endpoint,
		client: resty.New().
			SetRetryWaitTime(time.Second * time.Duration(retryWaitSeconds)).
			SetRetryCount(retryCount).
			SetAuthToken(token),
	}
}

type NotifyRequest struct {
	ID    int64
	phone int
	text  string
}

var (
	ErrorInvalidData = errors.New("bad request")
)

func (p Poster) Send(m entity.Message) error {
	ph, err := strconv.Atoi(m.Phone[1:])
	if err != nil {
		return ErrorInvalidData
	}

	r := NotifyRequest{
		ID:    rand.Int63(),
		phone: ph,
		text:  m.Text,
	}

	fmt.Printf("Отправка %+v\n", r)

	resp, err := p.client.R().
		SetBody(r).
		Post(p.EndPoint)

	if err != nil {
		return err
	}

	if resp.StatusCode() == 400 {
		return ErrorInvalidData
	}

	return nil
}

func (p Poster) Work(in chan entity.Message) (chan entity.Message, error) {
	done := make(chan entity.Message)
	go func() {

		for m := range in {

			err := p.Send(m)

			// обрабатываем ошибки
			// и помечаем сообщение
			if err != nil {
				if errors.Is(err, ErrorInvalidData) {
					m.Invalid()
				} else {
					m.Failed()
				}
			} else {
				m.SentNow()
			}
			done <- m
		}
		close(done)
	}()

	return done, nil
}
