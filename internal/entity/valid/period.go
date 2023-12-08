package valid

import (
	"time"
)

// Period проверяет, что указанный интервал - валидный
// например время начала больше времени конца)
func Period(start, end time.Time) error {
	if end.Sub(start) < 0 {
		return ErrorInvalidPeriod
	}
	return nil
}
