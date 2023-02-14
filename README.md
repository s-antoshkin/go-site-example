# Golang site example
Пример представляет собой небольшой сайт - "приглашение на вечеринку": формы, макеты, взаимодействие с БД `Postgres(pgx)` + тесты.

## Requirements
Go 1.19 or above.

## Installation
```
git clone https://github.com/s-antoshkin/go-site-example
cd go-site-example
go get -u
psql -h localhost -p 5432 -d rsvp_db -U postgres -f rsvp.sql
```

## Start the application
```
go run .
```
## Start the application
```
go run .
```