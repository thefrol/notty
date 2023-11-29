package chans

import "sync"

// FainIn собирает сообщения из нескольких каналов в один,
// после того как входящие каналы закроются, этот закроется сам собой
func FanIn[T any](ins ...<-chan T) chan T {
	out := make(chan T)

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

	// эта горутина дожидается,
	// когда закроются все входные каналы
	// и закрывает свой выходной
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
