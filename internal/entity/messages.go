package entity

import "time"

// небольшая бизнес-логика сообщений
const (
	StatusSent   = "done"
	StatusFailed = "fail"
)

// SentNow помечает, что сообщение было отправлено вот только что
func (m *Message) SentNow() {
	m.Status = StatusSent
	t := time.Now()
	m.Sent = &t
}

// Failed помечает, что сообщение не удалось отправить
func (m *Message) Failed() {
	m.Status = StatusFailed
}
