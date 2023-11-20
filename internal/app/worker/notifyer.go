package worker

import (
	"fmt"
	"log"
	"sync"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/postman"
	"gitlab.com/thefrol/notty/internal/stream"
)

type Notifyer struct {
	Messages stream.MessageStream
	Poster   postman.Poster
}

// FindAndSend обнаруживает новые сообщения и отправляет
//
//	// у нас работает конвеер
//
// Некоторые этапы можно посерьезней распараллелить,
// вопрос просто в желании, мне кажется параллелить
// надо только отправку
func (app Notifyer) FindAndSend(batch int, workers int) {

	// если же есть, то плодим горутины
	protos, err := app.Messages.Generate(batch)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := app.Messages.Create(protos)
	if err != nil {
		log.Fatal(err)
	}

	// отправляем в кучу горутин
	var dones []chan entity.Message
	for i := 0; i < workers; i++ {
		done, err := app.Poster.Work(msgs)
		if err != nil {
			log.Fatal(err)
		}
		dones = append(dones, done)
	}

	/// собираем все в кучу и обновляем
	end, err := app.Messages.Update(FanIn(dones...))
	if err != nil {
		log.Fatal(err)
	}

	// и все выбрасываем
	terminate(end)

}

func (app Notifyer) TryToResend(batch, workers int) {
	// получим зафейленные ранее сообщения
	in, err := app.Messages.Failed(batch)
	if err != nil {
		return
	}

	// отправляем в кучу горутин
	var dones []chan entity.Message
	for i := 0; i < workers; i++ {
		done, err := app.Poster.Work(in)
		if err != nil {
			log.Fatal(err)
		}
		dones = append(dones, done)
	}

	// логика такая что тут бы заглушку поставить,
	// типа только один тред может работать

	/// собираем все в кучу и обновляем
	end, err := app.Messages.Update(FanIn(dones...))
	if err != nil {
		log.Fatal(err)
	}

	// и все выбрасываем
	terminate(end)
}

// помошники по конкуретной работе
func FanIn(ins ...chan entity.Message) chan entity.Message {
	out := make(chan entity.Message)

	wg := sync.WaitGroup{}
	for _, in := range ins {
		wg.Add(1)

		in := in
		go func() {
			defer wg.Done()

			for m := range in {
				out <- m
			}
		}()
	}
	go func() {
		wg.Wait()

		fmt.Println("каналы закрыты")
		close(out)

	}()

	return out
}

func terminate(ins ...chan entity.Message) {
	end := FanIn(ins...)
	for range end {
	}
}
