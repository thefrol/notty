package service

import (
	"gitlab.com/thefrol/notty/internal/dto"
)

type StatisticsAdapter interface {
	All() (dto.Statistics, error)
	Filter(string, string, string) (dto.Statistics, error)
}

type StatisticsService struct {
	adapter StatisticsAdapter
}

func NewStatistics(a StatisticsAdapter) StatisticsService {
	return StatisticsService{
		adapter: a,
	}
}

func (s StatisticsService) All() (dto.Statistics, error) {
	return s.adapter.All()
}

// Filter возвращает статистику по сообщениям, где можно указать
// subId - имя рассылки, в формате SQL Like
// customerId - имя клиента, в таком же формате
// status - сообщения определенного статуса
func (s StatisticsService) Filter(subId, customerId, status string) (dto.Statistics, error) {
	return s.adapter.Filter(subId, customerId, status)
}
