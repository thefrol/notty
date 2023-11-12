# Сервис уведомлений Notty

+ Постгрес с `pgx`
+ миграции с `goose`
+ докер для запуска инфраструктуры
+ интегральные тесты по репозиториям с постгрёй

## АПИ сервер

+ Мигрирует в базу сам, базу нужно указать
+ Запускается на порту `8080`
+ По маршруту `/docs/` покажет сваггер

Конфигурируется переменными окружения:

+ `NOTTY_DSN` - строка подключения к постгре

Есть скрипт, чтобы быстро запустить базу, см. раздел разработка

## Допзадания

### Сваггер интерфейс

Доступен по адресу `/docs/`

### Тестирование на Гитлабе

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
