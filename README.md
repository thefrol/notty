# Сервис уведомлений Notty

+ Постгрес с `pgx`
+ миграции с `goose`
+ докер для запуска инфраструктуры
+ интегральные тесты по репозиториям с постгрёй
+ быстрая серилизация `easyjson`
+ Валидация запросов - в том числе на мат

## АПИ сервер

+ Мигрирует в базу сам, базу нужно указать
+ Запускается на порту `8080`
+ По маршруту `/docs/` покажет сваггер
+ Очень быстрый, запросы на вытягивание **2-3 мс**

Конфигурируется переменными окружения:

+ `NOTTY_DSN` - строка подключения к постгре

Есть скрипт, чтобы быстро запустить базу, см. раздел разработка

## Допзадания

+ Доки по маршруту `/docs/`
+ Юнит тестирование, и ***интегральное тестирование** в пайплайне на гитлабе

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

+ хорошая семантика, например `validate.ClientRequest(c)`
+ скорость
+ `CI`, тесты в гитлабе вообще сказка. Я честно не знал, что гитлаб так хорош
+ Тулинг: запуск бд, очистка и наполенния тестовыми данными в `./tool/` и скриптах

## Что не удалось

+ `App` - этот класс должен содержать высокоуровневую логику, например, как именно обновляются подписки итд. Но там ничего нет. По сути особо ничего высокоуровневого у меня нет, но и как факт, все забито в хендлеры
+ Сущность рассылки называется `Subscription`, тогда как её суть полностью противоположна

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
