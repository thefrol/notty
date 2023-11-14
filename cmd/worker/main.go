// Воркер занимается отправкой сообщений
//
// Раз во сколько-то секунд он просыпается, и смотрит
// Есть ли активные рассылки, и если нет, то засыпает.
// Если активные рассылки есть, то создает воркеров,
// которые будут отправлять это все на эндпоинт
//
// Как бы он на кажой итерации создает такую сеть горутин.
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/postman"
	"gitlab.com/thefrol/notty/internal/storage/postgres"
	"gitlab.com/thefrol/notty/internal/storage/subscriptions"
	"gitlab.com/thefrol/notty/internal/stream"
)

const batchSize = 50
const timeout = 3 * time.Second
const WorkerCount = 50

func main() {
	// конфигурируем
	dsn, ok := os.LookupEnv("NOTTY_DSN")
	if !ok {
		fmt.Println("Нужно передать строку подключения в переменной NOTTY_DSN")
		os.Exit(3)
	}

	// соединяемся с БД
	db := postgres.MustConnect(dsn)

	//создаем сервисы
	mes := stream.Messages(db)

	subs := subscriptions.New(db)

	post := postman.Poster{}

	//и за дело

	for {
		time.Sleep(timeout)

		// проверяем сколько рассылок активно
		active, err := subs.Active()
		if err != nil {
			log.Fatal(err)
		}

		if len(active) == 0 {
			// если никого, то спим
			continue
		}

		// если же есть, то плодим горутины
		// у нас работает конвеер
		//
		// Некоторые этапы можно посерьезней распараллелить,
		// вопрос просто в желании, мне кажется параллелить
		// надо только отправку

		protos, err := mes.Generate(batchSize)
		if err != nil {
			log.Fatal(err)
		}

		msgs, err := mes.Create(protos)
		if err != nil {
			log.Fatal(err)
		}

		// отправляем в куче горутин
		var dones []chan entity.Message
		for i := 0; i < WorkerCount; i++ {
			done, err := post.Work(msgs)
			if err != nil {
				log.Fatal(err)
			}
			dones = append(dones, done)
		}

		/// собираем все в кучу и обновляем
		end, err := mes.Update(FanIn(dones...))
		if err != nil {
			log.Fatal(err)
		}

		terminate(end)

	}

}

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
	for _ = range end {
	}
}
