### Фильмотека(сервис поиска фильмов, адмика для управления данными)

- Запуск базы:  `docker-compose up -d`

- Миграция базы: `goose -dir db/migrations postgres "host=localhost dbname=docker user=docker password=docker sslmode=disable" up`

- Запуск сервиса: `go run cmd/app/app.go`

- Тест в браузере:
    1. Поиск фильма:
        - Точное совпадение: `http://localhost:8080/get_film?film=Robocop 2`
        - Поиск по части названия: `http://localhost:8080/get_film?film=Rob`
    2. Поиск актера: