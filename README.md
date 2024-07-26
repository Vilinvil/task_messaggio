## Тестовое задание messaggio 

тг для быстрой связи: https://t.me/Vilin0

Условия задания лежат [в файле](task.md)

### Установка утилиты для накатки миграций
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
```

### Пример команды для накатки миграции
```bash
 migrate -database postgres://postgres:postgres:5432/postgres?sslmode=disable -path migrations/ up
```

### Установка утилиты для генерации документации к api
```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.3
```