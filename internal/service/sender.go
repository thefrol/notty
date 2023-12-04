package service

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/pkg/chans"
)

// Worker пользуется слоем приложения, чтобы
// генерировать и отправлять сообщения
type Worker struct {
	Notifyer  *app.Notifyerrrr
	Timeout   time.Duration
	BatchSize int
}

// FetchAndSend основная функция для воркера, тут он создает и
// отправляет сообщения
//
// После завершения контеста, функция перейдет к остановке, и
// будет ждать завершения всех горутин, после чего отдаст
// управление
func (w Worker) FetchAndSend(ctx context.Context) {
	// создадим вейтгруппу, которая будет ждать завершения всех работ
	// когда остановится контекст ctx
	wg := sync.WaitGroup{}
loop:
	for {
		// ищем новые сообщения

		// todo стоит обратить внимание, что создавать в базе данных сообщения
		// мы можем независимо от отправки, и копить их как-то под отправку, и
		// потом отправлять
		//
		// Но  в данной ситуации я предпочитаю все каналы в рамках одной итерации открывать и закрывать

		// Создаем канал с новосгенерированными сообщениями
		newMessages := chans.GeneratorFunc(func() []entity.Message {
			ms, err := w.Notifyer.CreateMessages(w.BatchSize)
			if err != nil {
				// продолждаем работу. Если ошибка, просто опять уходим в таймаут
				log.Error().
					Str("Message", "Ошибка при генерации сообщений из базы").
					Err(err)
			}
			return ms
		})

		// Также резервируем сообщения из тех, что ранее не получилось
		// отправить. Они лежат в базе со статусов `fail`
		resendsMessages := chans.GeneratorFunc(func() []entity.Message {
			ts, err := w.Notifyer.ReserveMessages(w.BatchSize, entity.StatusFailed)
			if err != nil {
				// продолждаем работу если ошибка, просто опять уходим в таймаут
				log.Error().
					Str("Message", "Ошибка при получении сообщений из базы").
					Err(err)

			}
			return ts
		})

		//собираем оба канала в один
		MessagesToSend := chans.FanIn(newMessages, resendsMessages)

		// И теперь с каждым сообщеним работаем, пытаемся отправить.
		// если не получается, то ждем
		// По сути мы не двинемся дальше пока не раскидаем по воркерам
		// все задачки
		for m := range MessagesToSend {
			m := m // на всякий
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := w.Notifyer.SendSMS(m)
				if err != nil {
					log.Error().
						Str("Message", "Ошибка при отправке смс").
						Err(err)

					// если ошибка отправки, то ставим статус -ошибка
					m.Failed()
				}
				// если нет ошибок, то ставим нужный статус и отправляем
				m.SentNow()

				// теперь у нас есть сообщение с нужным статусом,
				// мы просто его обновляем в базе
				err = w.Notifyer.UpdateMessage(m) // todo можно сделать updatefast без возвращения значения)
				if err != nil {
					log.Error().
						Str("Message", "Сообщение было отправлено, но не удалось обновить статус в хранилище").
						Err(err)
				}
			}()
		}

		select {
		case <-ctx.Done():
			// выходим если контекст закончился
			break loop
		default:
			// если не закончился
			// спим какое-то время.
			time.Sleep(w.Timeout)

			// todo я вообще не понимаю нужен ли нам таймаут или нет)
			// потому что он нужен условно только если мы не создали ни одного
			// сообщения
			//
			// возможно как раз можно из вейтгруппы прочитать
		}

	}
	log.Info().
		Str("Message", "Воркер ожидает остановку горутин")
	wg.Wait()

	log.Info().
		Str("Message", "Воркер завершает работу")
}
