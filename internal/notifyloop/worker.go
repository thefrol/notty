package notifyloop

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"gitlab.com/thefrol/notty/internal/app"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/pkg/chans"
)

// Worker пользуется слоем приложения, чтобы
// генерировать и отправлять сообщения
type Worker struct {
	// Сервисное приложение, которое будет использоваться
	Notifyer *app.Notifyerrrr

	// таймаут между обращением на поиск новых сообщений в базе,
	// которые можно отправить. Например, если не отправлено
	// ни одного сообщения, то сервис засываем на это время
	Timeout time.Duration

	// размер батча, котрый мы достаем каждый раз
	BatchSize int

	//логгер
	Logger zerolog.Logger

	// базовый контекст на манер http.Server{}
	BaseContext context.Context
}

// FetchAndSend основная функция для воркера, тут он создает и
// отправляет сообщения
//
// После завершения контеста ctx, функция перейдет к остановке, и
// будет ждать завершения всех горутин, после чего отдаст
// управление
//
// Очень важно не путать контекст завершения ctx и базовый контекст
// Worker.BaseContext. Базовый контекст используется для доступа к базам
// данных и прочему. Он не будет завершен при остановке воркера: при
// остановке воркера мы дождемся завершения цикла генерации и отправки
// сообщений и после этго уже прекратим работу.
func (w Worker) FetchAndSend(ctx context.Context) {
	// создадим вейтгруппу, которая будет ждать завершения всех работ
	// когда остановится контекст ctx
	wg := sync.WaitGroup{}

	// наследуем базовый контекст, который будет использоваться для
	// создания запросов на внешние сервиси и к инфраструктуре.
	// Нам нужно просто проверить, что он не nil, и в таком случае
	// создать новый контекст
	var workContext context.Context
	if w.BaseContext != nil {
		var cancel context.CancelFunc
		workContext, cancel = context.WithCancel(w.BaseContext)
		defer cancel()
	} else {
		workContext = context.Background()
	}

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
			ms, err := w.Notifyer.CreateMessages(workContext, w.BatchSize)
			if err != nil {
				// продолждаем работу. Если ошибка, просто опять уходим в таймаут
				w.Logger.Error().
					Err(err).
					Msg("Ошибка при генерации сообщений из базы")
			}
			return ms
		})

		// Также резервируем сообщения из тех, что ранее не получилось
		// отправить. Они лежат в базе со статусов `fail`
		resendsMessages := chans.GeneratorFunc(func() []entity.Message {
			ts, err := w.Notifyer.ReserveMessages(workContext, w.BatchSize, entity.StatusFailed)
			if err != nil {
				// продолждаем работу если ошибка, просто опять уходим в таймаут
				w.Logger.Error().
					Err(err).
					Msg("Ошибка при получении сообщений из базы")

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
				err := w.Notifyer.SendSMS(workContext, m)
				if err != nil {
					w.Logger.Error().
						Err(err).
						Msg("Ошибка при отправке смс")

					// если ошибка отправки, то ставим статус -ошибка
					m.Failed()
				}
				// если нет ошибок, то ставим нужный статус и отправляем
				m.SentNow()

				// теперь у нас есть сообщение с нужным статусом,
				// мы просто его обновляем в базе
				err = w.Notifyer.UpdateMessage(workContext, m) // todo можно сделать updatefast без возвращения значения)
				if err != nil {
					w.Logger.Error().
						Err(err).
						Msg("Сообщение было отправлено, но не удалось обновить статус в хранилище")
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
	w.Logger.Info().
		Msg("Воркер ожидает остановку горутин")
	wg.Wait()

	w.Logger.Info().
		Msg("Воркер завершает работу")
}
