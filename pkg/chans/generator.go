package chans

// GeneratorFunc делает канал из функции, которая генерирует батч.
// Указанная функция gen будет запущена в отдельной горутине и весь вывод отправится
// в канал
func GeneratorFunc[T any](gen func() []T) <-chan T {
	ch := make(chan T)

	go func() {
		for _, v := range gen() {
			ch <- v
		}
		close(ch)
	}()
	return ch
}
