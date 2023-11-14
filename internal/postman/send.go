package postman

import (
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

func (p Poster) Send(m entity.Message) {

	p.client.R().
		SetBody(msg).
		Post(p.EndPoint)
}

func (p Poster) Work(in chan entity.Message) (chan entity.Message, error) {
	done := make(chan entity.Message)
	go func() {
		count := 0
		for m := range in {
			count += 1
			time.Sleep(5 * time.Second)

			if count%2 == 0 {
				m.Failed()
			}
			m.SentNow()
			done <- m
		}
		close(done)
	}()

	return done, nil
}
