# Запускаем базу
docker compose up -d

# Установите migrate и air (для hot-reload)
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/cosmtrek/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


#  Проводим миграции

```
make migrate-down && make migrate-up
```

# Запускаем проект

 ```
make run
 ```
