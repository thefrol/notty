// Этот пакет предоставляет интерфейс к внешнему сервису
// отправки смс
package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"gitlab.com/thefrol/notty/internal/entity"
)

type PosterService struct {
	EndPoint string
	Token    string
	client   *resty.Client
}

func NewSMSEndpoint(endpoint string, retryWaitSeconds int, retryCount int, token string) PosterService {
	return PosterService{
		EndPoint: endpoint,
		client: resty.New().
			SetRetryWaitTime(time.Second * time.Duration(retryWaitSeconds)).
			SetRetryCount(retryCount).
			SetAuthToken(token).
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}),
		// todo bug надо как-то из докера научиться проверять сертификаты
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

func (p PosterService) Send(m entity.Message) error {

	ph, err := strconv.Atoi(m.Phone[1:])
	if err != nil {
		return ErrorInvalidData
	}

	id := rand.Int63()

	r := NotifyRequest{
		ID:    id,
		phone: ph,
		text:  m.Text,
	}

	u, err := url.JoinPath(p.EndPoint, strconv.Itoa(int(id)))
	if err != nil {
		return err
	}

	fmt.Printf("Отправка %+v\n", r)

	resp, err := p.client.R().
		SetBody(r).
		Post(u)

	fmt.Println("RESPONSE:", resp.StatusCode())

	if err != nil {
		return err
	}

	if resp.StatusCode() == 400 {
		return ErrorInvalidData
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("Неизвестная ошибка %v", err)
	}

	return nil
}

func (p PosterService) Work(in <-chan entity.Message) (<-chan entity.Message, error) {
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
