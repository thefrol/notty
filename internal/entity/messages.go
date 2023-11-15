package entity

import "time"

// небольшая бизнес-логика сообщений
const (
	StatusSent    = "done"
	StatusFailed  = "fail"
	StatusInvalid = "invalid"
)

// SentNow помечает, что сообщение было отправлено вот только что
func (m *Message) SentNow() {
	// todo
	//
	// тут надо подумать, а что если я
	// я попытаюсь поставить там
	// для неправильного сообщения что оно
	// отправлено. возможно ли такое?
	// итд
	m.Status = StatusSent
	t := time.Now()
	m.Sent = &t
}

// Failed помечает, что сообщение не удалось отправить
func (m *Message) Failed() {
	m.Status = StatusFailed
}

// Invalid помечает, что сообщение не удалось отправить
func (m *Message) Invalid() {
	m.Status = StatusInvalid
}
