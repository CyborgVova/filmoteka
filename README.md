### Фильмотека

- Сервис поиска фильмов, админка для управления данными

- Запуск сервиса: `docker-compose up` (Ожидаем старт сервера: `Start server on port :8080 ...`)

- Данные для авторизации(Basic):

  - login: `admin`
  - password: `secret`

- Тест с использованием утилиты `curl`:
  1. Поиск фильма:
     - Точное совпадение: `curl -X GET -H "Content-Type: application/json" http://localhost:8080/get_film?film=Robocop\%202`
     - Поиск по части названия: `curl -X GET -H "Content-Type: application/json" http://localhost:8080/get_film?film=Rob`
  2. Поиск актера:
     - Точное совпадение: `curl -X GET -H "Content-Type: application/json" http://localhost:8080/get_actor?actor=Anita\%20Tsoy`
     - Поиск по части имени: `curl -X GET -H "Content-Type: application/json" http://localhost:8080/get_actor?actor=Ani`
  3. Удаление фильма:
     - `curl -X DELETE -u admin:secret -H "Content-Type: application/json" -d "{\"title\": \"Robocop 3\", \"release\": 1994}" http://localhost:8080/delete_film`
  4. Удаление актера:
     - `curl -X DELETE -u admin:secret -H "Content-Type: application/json" -d "{\"fullname\": \"Anita Tsoy\", \"dateofbirth\": \"1982-07-11\"}" http://localhost:8080/delete_actor`
  5. Изменить информацию о фильме:
     - `curl -X PATCH -u admin:secret -H "Content-Type: application/json" -d "{\"title\": \"Robocop 2(Next Step)\", \"rating\": 9}" http://localhost:8080/set_film/2`
  6. Изменить информацию о актере:
     - `curl -X PATCH -u admin:secret -H "Content-Type: application/json" -d "{\"fullname\": \"Sergey Ivanovich Burunov\"}" http://localhost:8080/set_actor/3`
  7. Добавить фильм:
     - `curl -X POST -u admin:secret -H "Content-Type: application/json" -d "{\"title\": \"Robocop 4\", \"description\": \"New part about dead police officer\", \"release\": 1996, \"rating\": 10}" http://localhost:8080/add_film`
  8. Добавить актера:
     - `curl -X POST -u admin:secret -H "Content-Type: application/json" -d "{\"fullname\": \"Leonid Yarmolnik\", \"sex\": \"male\", \"dateofbirth\": \"07-22-1968\"}" http://localhost:8080/add_actor`
