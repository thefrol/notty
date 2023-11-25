// Одна из задач это уведомление клиентом,
// причем тут мы обрабатываем тот случай, когда
// мы создали сообщение и отправили его
// если ошибка, то пометили сообщение
package usecase

import (
	"errors"
	"log"

	"gitlab.com/thefrol/notty/internal/app"
)

func FindAndSend(app app.App, notty app.Notifyer, batch int, workers int) {

	// // если же есть, то плодим горутины
	// msgs, err := app.CreateMessages(batch)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // отправляем в кучу горутин
	// var dones []<-chan entity.Message
	// for i := 0; i < workers; i++ {
	// 	err := notty.SendNotification(msgs)
	// 	if err != nil {
	// 		//todo
	// 		m
	// 	}
	// 	dones = append(dones, done)
	// }

	// /// собираем все в кучу и обновляем
	// end, err := app.Messages.Update(FanIn(dones...))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // и все выбрасываем
	// terminate(end)
	log.Fatal(errors.ErrUnsupported)
}
