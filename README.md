## Тестовое задание messaggio 

тг для быстрой связи: https://t.me/Vilin0

Условия задания лежат [в файле](task.md)

Решение представляет собой документацию swagger, доступную в интернете по ссылке:
https://goods-galaxy.ru:7443/api/v1/swagger/index.html

По условию задания составил схему проекта:
![schema](img/schema.png)

На сайте можно сделать запрос статистики /message/statistic - получить информацию кол-во 
всего и обработанных сообщений. Добавить сообщение в систему по эндпоинту /message.
Для проверки можно отправить сообщение в систему по /message, проверить, что всего сообщений
увеличилось, а обработанных нет. Обработчик внутри кода обрабатывает сообщения последовательно
(в качестве симуляции спит 5сек), и секунд через 15 можно проверить статистику по /message/statistic и заметить, 
что кол-во обработанных сообщений сравнялось с общим кол-вом.

## Для удаленного развертывания 
Здесь используется https для доступа из интернета, поэтому понадобится сертификаты, для 
доменного имени, можно сделать, например, с помощью certbot. Если нет доменного имени,
то можно развернуть без https [локально](#для-локального-развертывания)

1. Скачать репозиторий с кодом
```bash
git clone https://github.com/Vilinvil/task_messaggio.git
```
2. Скопировать шаблоны env var в рабочую папку
```shell
make cp-env
```
3. Далее требуется заполнить env переменные 

| Файлы                            | Переменная                           | Назначение и причина                                                                                          |
|----------------------------------|--------------------------------------|---------------------------------------------------------------------------------------------------------------|
| .env.postgres                    | POSTGRES_PASSWORD                    | Пароль к постгресу - исходя из требований безопасности                                                        |
| .env.message, .env.messageworker | URL_DATABASE                         | URL подключения к постгресу с паролем - исходя из требований безопасности                                     |
| .env.message                     | DOMAIN                               | Домен сайта - для генерации документации и последующей возможности, выполнять запросы в интерфейсе в браузере |
| .env.nginx                       | SSL_CERTIFICATE, SSL_CERTIFICATE_KEY | Пути в системе до сертификатов - для работы TLS шифрования                                                    |

Остальные переменные можно оставить без изменения.
4. Экспортировать некоторые переменные из файла nginx в терминал, для возможности
использования их в docker compose файле. DOMAIN задать в соответствии с env.
```bash
export $(grep -v '^#' .env/.env.nginx | xargs -d '\n') &&
export DOMAIN=my.domain.ru
```

Проверить можно с помощью 
```bash
echo $DOMAIN $SSL_CERTIFICATE $SSL_CERTIFICATE_KEY
```

5. Запустить все через docker compose 
```bash
docker compose -f deployments/docker-compose.yaml up --build -d 
```
6. Накатить миграции в базе данных (не забыть поменять пароль в урле
на соответствующий env POSTGRES_PASSWORD)
```bash
docker exec deployments-message-1 ./migrate -database postgres://postgres:PasswordNeedChange@postgres:5432/postgres?sslmode=disable -path migrations up
```

Система готова. Можно увидеть в логах с помощью команды примерно следующее:
```bash
docker compose -f deployments/local-docker-compose.yaml logs -f message messageworker
```

```
message-1        | {"level":"info","ts":1722542703.6403918,"caller":"server/server.go:41","msg":"Start server:8090"}
messageworker-1  | {"level":"info","ts":1722542762.258721,"caller":"server/server.go:64","msg":"ServerWorker start job messages"}
```

## Для локального развертывания
Здесь не используется nginx и https, не нужны сертификаты и система будет работать
по ссылке http://localhost:7443

1. Скачать репозиторий с кодом
```bash
git clone https://github.com/Vilinvil/task_messaggio.git
```
2. Скопировать шаблоны env var в рабочую папку
```shell
make cp-env
```
3. Далее требуется заполнить env переменные

| Файлы                            | Переменная        | Назначение и причина                                                                                                    |
|----------------------------------|-------------------|-------------------------------------------------------------------------------------------------------------------------|
| .env.postgres                    | POSTGRES_PASSWORD | Пароль к постгресу - исходя из требований безопасности                                                                  |
| .env.message, .env.messageworker | URL_DATABASE      | URL подключения к постгресу с паролем - исходя из требований безопасности                                               |
| .env.message                     | DOMAIN=localhost  | Домен сайта localhost - для генерации документации и последующей возможности, выполнять запросы в интерфейсе в браузере |

Остальные переменные можно оставить без изменения.
4. Экспортировать DOMAIN=localhost для возможности
   использования их в docker compose файле. 
```bash
export DOMAIN=localhost
```

Проверить можно с помощью
```bash
echo $DOMAIN
```

5. Запустить все через docker compose
```bash
docker compose -f deployments/local-docker-compose.yaml up --build -d 
```
6. Накатить миграции в базе данных (не забыть поменять пароль в урле
   на соответствующий env POSTGRES_PASSWORD)
```bash
docker exec deployments-message-1 ./migrate -database postgres://postgres:PasswordNeedChange@postgres:5432/postgres?sslmode=disable -path migrations up
```

Система готова. Можно увидеть в логах с помощью команды примерно следующее:
```bash
docker compose -f deployments/local-docker-compose.yaml logs -f message messageworker
```

```
message-1        | {"level":"info","ts":1722542703.6403918,"caller":"server/server.go:41","msg":"Start server:8090"}
messageworker-1  | {"level":"info","ts":1722542762.258721,"caller":"server/server.go:64","msg":"ServerWorker start job messages"}
```

## Команды для разработки

### Установка утилиты для накатки миграций
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
```

### Установка утилиты для генерации кода сериализации easyjson
```bash
go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@v0.7.7
```

### Пример команды для накатки миграции
```bash
 migrate -database postgres://postgres:postgres:5432/postgres?sslmode=disable -path migrations/ up
```

### Установка утилиты для генерации документации к api
```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.3
```
