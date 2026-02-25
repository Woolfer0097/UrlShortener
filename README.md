# UrlShortener

## Запуск

```
cp .env.example .env && docker compose up -d
```

## Покрыл важный функционал unit-тестами

```
go test ./...
```

## Тип хранилища
В .env

STORAGE_TYPE=postgresql/in_memory

## Docker Образ

Залил на свой dockerhub woolfer0097kek/url-shortener:1.0.1
