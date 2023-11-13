// Приложение которое смотрит какие сообщения готовы к отправке
package app

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/storage/messages"
	"gitlab.com/thefrol/notty/internal/storage/subscriptions"
)

type Obs struct {
	MessagesRepo      messages.Messages
	SubscriptionsRepo subscriptions.Subscriptions
	wg                sync.WaitGroup
}

func (o *Obs) FetchAndWork(timeout time.Duration, workers int) {
	for {
		fmt.Println("NEW ITERATION")
		time.Sleep(timeout)
		active, err := o.SubscriptionsRepo.Active()
		if err != nil {
			log.Fatal(err)
		}
		if len(active) == 0 {
			log.Println("Нет активных подписок")
			continue
		}

		ch, err := o.MessagesRepo.NewMessages(300)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < workers; i++ {
			o.worker(ch)
		}

		// Обязательно дождаться завершения этой итерации
		o.wg.Wait()
	}
}

var counter = atomic.Int64{}

func (o *Obs) worker(ch <-chan entity.Message) {
	o.wg.Add(1)
	counter.Add(1)
	workerN := int(counter.Load())
	go func() {
		for m := range ch {
			//todo тут может быть не новое сообщение
			fmt.Println(workerN, "->", m)
			m.Status = "created"
			cm, err := o.MessagesRepo.Create(m)
			if err != nil {
				fmt.Println(err)
			}
			// попытка отправить
			time.Sleep(time.Second * 30)

			// справляеся с ошибкой
			cm.Status = "failed"
			_, err = o.MessagesRepo.Update(cm)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// если не получилось, то удалить или ещё чего как отмаркировать
		}
		o.wg.Done()
		fmt.Println("THREAD DONE", workerN)
	}()
}
