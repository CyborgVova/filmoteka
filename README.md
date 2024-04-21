### Фильмотека(сервис поиска фильмов, адмика для управления данными)

- Запуск сервиса: `docker-compose up` (Ожидаем старт сервера: `Start server on port :8080 ...`)

- Данные для авторизации(Basic):

  - login: `admin`
  - password: `secret`

- Тест с использованием утилиты `curl`:
  1. Поиск фильма:
     - Точное совпадение: `curl -X GET "Content-Type: application/json" http://localhost:8080/get_film?film=Robocop 2`
     - Поиск по части названия: `curl -X GET "Content-Type: application/json" http://localhost:8080/get_film?film=Rob`
  2. Поиск актера:
  - Точное совпадение: `curl -X GET "Content-Type: application/json" http://localhost:8080/get_actor?actor=Anita Tsoy`
  - Поиск по части имени: `curl -X GET "Content-Type: application/json" http://localhost:8080/get_actor?actor=Ani`
  3. Удаление фильма:
     - `curl -X POST -H "Content-Type: application/json" -d "{\"title\": \"Robocop 3\", \"release\": 1994}" http://localhost:8080/delete_film`
  4. Удаление актера:
  - `curl -X POST -H "Content-Type: application/json" -d "{\"fullname\": \"Anita Tsoy\", \"dateofbirth\": \"1982-07-07\"}" http://localhost:8080/delete_actor`
  5. Изменить информацию о фильме:
  - `curl -X PATCH -H "Content-Type: application/json" -d {\"title\": \"Robocop 2(Next Step)\", \"rating\": 9} http://localhost:8080/set_film/2`
  6. Изменить информацию о актере:
  - `curl -X PATCH -H "Content-Type: application/json" -d {\"fullname\": \"Anita Tsoy-Ivanova\"} http://localhost:8080/set_actor/1`
