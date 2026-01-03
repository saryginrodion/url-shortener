# url shortener

- auth - email + password
- opaque tokens

- MW checks token and returns sets context key "user"

- POST /urls - создать ссылку
- DELETE /urls/{id} - удалить ссылку
- GET /urls/{id} - перейти по ссылке
- GET /urls/{id}/metrick - метрики ссылки

# libs
- cleanenv -  
- chi - routing
- sqlx - DB
- https://github.com/pressly/goose - migrations

# что должно быть еще
- openapi + swagger
- ELK + trace_id генерящийся в MW + передающийся в контексте повсюду
- prometheus + grafana (желательно обложиться метриками вплоть до отдельных запросов к БД)
- если получится - можно поднять кубер
