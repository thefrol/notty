package storage

import (
	"database/sql"
	"fmt"

	"gitlab.com/thefrol/notty/internal/dto"
	"gitlab.com/thefrol/notty/internal/entity"
	protootypes "gitlab.com/thefrol/notty/internal/storage/prototypes"
)

type MessageStream struct {
	MessageRepository MessageRepo
	Spawner           protootypes.Prototypes
}

type MessageRepo interface {
	ByStatus(string, int) ([]entity.Message, error)
	Update(entity.Message) (entity.Message, error) // todo это какая-то хорошая сигнатура
	Create(entity.Message) (entity.Message, error)
	Spawn(int, string) ([]entity.Message, error)
}

// todo это сервис и его надо в сервисы с интерфейсом от адаптеров

func Messages(db *sql.DB, messages MessageRepo) MessageStream {
	return MessageStream{
		MessageRepository: messages,
		Spawner:           protootypes.New(db),
	}
}

func (stream MessageStream) Generate(batch int) (chan dto.Prototype, error) {

	b, err := stream.Spawner.Spawn(batch)
	if err != nil {
		return nil, err
	}

	out := make(chan dto.Prototype)

	go func() {
		for _, msg := range b {
			out <- msg
		}
		close(out)
	}()

	return out, nil
}

func (stream MessageStream) Failed(batch int) (chan entity.Message, error) {

	b, err := stream.MessageRepository.ByStatus(entity.StatusFailed, batch)
	if err != nil {
		return nil, err
	}

	out := make(chan entity.Message)

	go func() {
		for _, msg := range b {
			out <- msg
		}
		close(out)
	}()

	return out, nil
}

func (stream MessageStream) SetStatus(in chan entity.Message, status string) (chan entity.Message, error) {
	out := make(chan entity.Message)
	go func() {
		for m := range in {
			m.Status = status
			u, err := stream.MessageRepository.Update(m)
			if err != nil {
				fmt.Println("marker:", err)
				continue
			}
			out <- u
		}
		close(out)

		fmt.Println("Marker closed")
	}()

	return out, nil
}

func (stream MessageStream) Update(in chan entity.Message) (chan entity.Message, error) {
	out := make(chan entity.Message)
	go func() {
		for m := range in {
			u, err := stream.MessageRepository.Update(m)
			if err != nil {
				fmt.Println("marker:", err)
				continue
			}
			out <- u
		}
		close(out)

		fmt.Println("Marker closed")
	}()

	return out, nil
}

func (stream MessageStream) Create(in chan dto.Prototype) (chan entity.Message, error) {
	out := make(chan entity.Message)
	go func() {
		for p := range in {
			msg := p.ToMessage()
			cr, err := stream.MessageRepository.Create(msg)
			if err != nil {
				fmt.Println("creator:", err)
				continue
			}
			out <- cr
		}
		close(out)

		fmt.Println("Creator closed")
	}()
	return out, nil
}

// todo
//
// в App должна лежать вся высокоуровневая логика,
// типа, какие рассылки мы должны обновлять и как
// может быть например, стоит проверить, что
// в рассылке уже есть
//
// Например, запретить обновлять время начала для рассылки
// Вообще у рассылки должно быть что-то вроде Put
// И тут тожно нужны какие-то Nullable типы как будето
//
// Но вернемся к логике, например, у нас есть условие,
// нельзя обновлять уже завершенные рассылки,
// или нельзя менять время начала на более раннее, если
// это время в прошлом
//
// ИЛи например, нельзя матерные слова в тексте, и фильтр
// или проверка текста на фейлы, странные символы. Например,
// запретить какие-нить кракозябры или ссылки, ну кто знает
// все это может быть тут
//
// Это вроде и валидация, просто высокоуровневая
