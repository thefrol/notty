// Этот пакет предоставляет интерфейс к внешнему сервису
// отправки смс
package sms

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/entity"
)

// PosterService это надстройка над HTTP клиентом для сервиса
// отправки смсок
type PosterService struct {
	EndPoint string
	Token    string
	client   *resty.Client
}

// NewEndpoint позволяет отправлять смски, путем HTTP запросов
// на endpoint. Сервис требует авторизации через jwt, поэтому
// нужно указать Bearer, который передастся в заголовке Authorization
//
// При этом можно указать количество повторных попыток при запросе
// и вреся между ними
func NewEndpoint(endpoint string, retryWaitSeconds int, retryCount int, token string) PosterService {
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

func (p PosterService) Send(ctx context.Context, m entity.Message) error {

	ph, err := strconv.Atoi(m.Phone[1:])
	if err != nil {
		return ErrorInvalidData
	}

	// тут конечно вооще не понятно, что происходит.
	// так уж получилось, что спроектировав сервис,
	// я обнаружил что тут айдишник это int64
	//
	// что он означает я вообще не знаю,
	// поэтому чтобы сообщения нормально отправлялись
	// я просто генерирую случайный айдишник при отправке)
	//
	// наверное, его по-хорошему логгировать надо, чтобы была какая-то
	// хоть всязь при необходимости
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

	log.Info().
		Int64("Id", r.ID).
		Str("Text", r.text).
		Int("Phone", r.phone).
		Msg("Отправка сообщения")

	resp, err := p.client.R().
		SetBody(r).
		SetContext(ctx).
		Post(u)

	if err != nil {
		log.Info().
			Int64("Id", r.ID).
			Str("Text", r.text).
			Int("Phone", r.phone).
			AnErr("err", err).
			Msg("ошибка запроса к смс-серверу")

		return err
	}
	log.Info().
		Int64("Id", r.ID).
		Str("Text", r.text).
		Int("Phone", r.phone).
		Int("ResponseCode", resp.StatusCode()).
		Msg("Результаты отправки сообщения")

	if resp.StatusCode() == 400 {
		return ErrorInvalidData
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("неизвестная ошибка %v", err)
	}

	return nil
}
