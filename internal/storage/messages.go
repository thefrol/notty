package storage

import "gitlab.com/thefrol/notty/internal/entity"

type MessageService struct {
	repo MessageRepo
}

func NewMessages(repo MessageRepo) MessageService {
	return MessageService{
		repo: repo,
	}
}

// Spawn создает сообщения из тех, что ещё не существуют
// но которые на данный момент можно отправить, и устанавливает
// им статус status, вернет не более n сообщений
func (ms MessageService) Spawn(n int, status string) ([]entity.Message, error) {
	return ms.repo.Spawn(n, status)
}
