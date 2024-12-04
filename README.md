Music Library API

Описание

Music Library API предоставляет функционал для управления музыкальной библиотекой. Он включает в себя:
Получение списка песен
Получение текста песни
Удаление песни
Добавление новой песни
API доступен по адресу: http://localhost:8080.
Как запустить проект

1. Убедитесь, что у вас установлены:
Docker
Docker Compose
2. Настройте .env файл
Переименуйте .env_example на .env файл в корне проекта и поменяйте в нем API_URL на url вашего API для получения 
дополнительной информации о песне.
   

4. Запустите приложение
Выполните команду:
docker-compose up --build
Docker Compose создаст и запустит два контейнера:
music_library_app — ваше приложение на Go.
postgres-container — база данных PostgreSQL.
5. Проверка API
После успешного запуска приложение будет доступно по адресу: http://localhost:8080.
Endpoints

1. Получить все песни
Метод: GET /songs
Описание: Возвращает список песен с фильтрацией и пагинацией.
Пример запроса:
GET /songs?page=1&limit=10
2. Получить текст песни
Метод: GET /lyrics
Описание: Возвращает текст песни.
Пример запроса:
GET /lyrics?id=123&page=1&limit=10
3. Удалить песню
Метод: POST /delete
Описание: Удаляет песню по её ID.
Пример запроса:
POST /delete?id=123
4. Добавить новую песню
Метод: POST /add
Описание: Добавляет новую песню в библиотеку.
Пример запроса:
POST /add?group=Queen&song=Bohemian Rhapsody
