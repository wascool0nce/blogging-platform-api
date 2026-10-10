# Blogging Platform API

Запуск базы данных: 
```bash
docker compose up --build -d
```

Запуск приложения:
```bash
go run ./cmd/blogging-platform-api
```

**Адрес API: http://localhost:8080**

Проверка:
```bash
curl -i http://localhost:8080/health
```
