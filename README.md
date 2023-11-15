# Сервис уведомлений Notty

Отправит уведомления с фильтрацией по тегам и операторам. Не пропускает мат. За счет конкурентности готов отправлять по миллиону сообщений в минуту.

Через [веб-апи](./api/openapi.yml) добавляются клиенты и рассылки. У каждой рассылки можно указать время начала и окончания. Если рассылка сейчас активна, то из клиентов отфильтруются те, что подходят по тегам и операторам, и отправит им сообщение.

Сервис [развертывается за 10 сек](#быстрое-развертывание-без-скачивания-репозитория)

Под капотом база данных с двумя микросервисами. Один держит веб-интерфейс, другой проверяет базу и отправляет сообщения. Подробнее [чуть ниже](#описание-сервиса)

Все [допазадания отмечены тут](#допзадания)

## Развертывание

Можно воспользоваться компоновщиком докера в [быстром деплое](./deploy/fast-compose/) или [когда репозиторий скачан](./deploy/docker-compose.yml)

При любом деплое можно сразу воспользоваться [пользовательским сценариаем](./test/sample.http), чтобы опробовать деплой. Там будет созданы получатели , будет создана рассылка. И сообщение тут же уйдет пользователю, но только одному

### Быстрое развертывание без скачивания репозитория

Чтобы запустить сервис за десять секунд, достаточно старого советского...

#### Шаг 1. Установка токена

```bash
export TOKEN=<вставь_токен_сюда> # ваш токен
```

#### Шаг 2. Запуск

```bash
# скачаем компоуз
curl https://gitlab.com/api/v4/projects/52099976/repository/files/deploy%2Ffast-compose%2Fdocker-compose.yml/raw?ref=main > docker-compose.yml

#запустим
docker-compose up
```

> если для докера требуются привилегии, стоит использовать `sudo -E`

#### Шаг 3. Успех

Доки станут доступны [тут](http://localhost/docs). Возпользуйтесь [подборкой запросов](./test/sample.http), чтобы сразу отправить первое сообщение.

### Если репозиторий скачан

То достаточно запустить скрипт:

```bash
TOKEN=<your_token> ./scripts/run-local
```

### В кубернететес

Так же легко соскальзывает в кубернетес

```bash
kubectl apply -f ./deploy/k8s
```

## Описание сервиса

### Стек

+ Постгрес с `pgx`
+ миграции с `goose`
+ докер для запуска инфраструктуры
+ интегральные тесты по репозиториям с постгрёй
+ быстрая серилизация `easyjson`
+ `chi` и `resty` для http
+ развертывание в `докере` при помощи `докер компоуз`. И ещё деплойменты и сервисы в `кубернетес`

### АПИ сервер

+ Мигрирует в базу сам, базу нужно указать
+ Запускается на порту `8080`
+ По маршруту `/docs/` покажет сваггер
+ Очень быстрый, запросы на вытягивание **2-3 мс**

Конфигурируется переменными окружения:

+ `NOTTY_DSN` - строка подключения к постгре

Есть скрипт, чтобы [быстро запустить базу](#запустить-локальный-докер-с-базой), см. раздел разработка

## Допзадания

+ Доки по маршруту `/docs/`
+ Юнит тестирование, и ***интегральное тестирование** в пайплайне на гитлабе
+ Сделан [Dockerfile](./Dockerfile)
+ Сделан [docker-compose.yml](./deploy/docker-compose.yml)
+ [Deployments и sevice](./deploy/k8s/) для кубернетес

### Сваггер интерфейс

Доступен по адресу `/docs/`

### Тестирование на Гитлабе

Локально можно запустить командой `go run ./...`, а ещё ...

#### Интегральное тестирования хранилищ

Находится в `.gitlab-ci.yml` вот тут

```yaml
repository_test:
  stage: test
  services:
    - name: postgres
      alias: test_db
  script:
    - export NOTTY_TEST_DB="host=test_db  dbname=test user=test password=test sslmode=disable"
    - export NOTTY_TEST_NOSKIP=1
    - go test -v ./test/...
```

## Что мне удалось

+ Валидация запросов - в том числе на мат
+ Конкуретная отправка сообщений одновременно с поиском сообщений к отправке
+ хорошая семантика, например `validate.ClientRequest(c)`
+ скорость
+ `CI`, тесты в гитлабе вообще сказка. Я честно не знал, что гитлаб так хорош
+ Тулинг: запуск бд, очистка и наполенния тестовыми данными в `./tool/` и скриптах

## Что не удалось

+ `App` - этот класс должен содержать высокоуровневую логику, например, как именно обновляются подписки итд. Но там ничего нет. По сути особо ничего высокоуровневого у меня нет, но и как факт, все забито в хендлеры
+ Сущность рассылки называется `Subscription`, тогда как её суть полностью противоположна
+ Логгирование
+ ДДД. Высокие уровки взаимодействуют с нижними без интерфейсов, перескоки через слои немношт. Это все печально
+ Во всех деплойментах у меня нет постоянных хранилищ
+ Забыл часовой пояс у клиента

## Разработка

### Запустить локальный докер с базой

```bash
./scripts/start-postgres
```

ПОсле успешного запуска скрипт выведет строку для подключения

```bash
База данных для нотти доступна по такой строке:
host=localhost dbname=notty_dev_local user=notty_dev_local password=local.dev.1 port=15432 sslmode=disable
```

И далее можно создать `launch.json` с такой вот конфигой

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Апи сервер",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "./cmd/server",
            "env":{
                "NOTTY_DSN":"host=localhost dbname=notty_dev_local user=notty_dev_local password=local.dev.1 port=15432 sslmode=disable"
            }
        }
    ]
}
```

## Тайминги

+ 1,5 часа Написать опенапи спеку и сгенерировать сервер
+ 1,5-3 час репозитории и подключение бд
+ 3-4 час хендлеры
+ 4-7 час отправщик и новые сервисы
+ 7-8 час синхронизация с ендпоинтом
