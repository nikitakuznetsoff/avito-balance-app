# avito-balance-app
## Запуск
Запуск осуществуляется с помощью `docker-compose up` из корня приложения
## Описание проекта
Основные моменты:
- Проект запускается на порту `:9000`
- В качестве HTTP роутера используется `https://github.com/gorilla/mux`
- В качестве хранилища используется база MySQL (v. 8.0.21)
  + База данных содержит 1 таблицу `users` c типом хранения данных InnoDB
  + Структуру таблиц можно посмотреть в файле инициализации `_sql/db.sql`
  
## Описание API
- `localhost:9000/balance` - получение баланса пользователя
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"userID": 2}' \
http://localhost:9000/balance
```
- `localhost:9000/balance?currency=USD` - получение баланса пользователя в заданной валюте
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"userID": 2}' \
http://localhost:9000/balance?currency=USD
```
- `localhost:9000/balance/withdraw` - снятие денег с баланса пользователя
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"userID": 2, "value": 50}' \
http://localhost:9000/balance/withdraw
```
- `localhost:9000/balance/deposit` - зачисление денег на баланс пользователя
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"userID": 2, "value": 10.22}' \
http://localhost:9000/balance/deposit
```
- `localhost:9000/balance/transfer` - перевод денег от одного пользователя к другому
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"sender": 2, "receiver": 3, "value": 50}' \
http://localhost:9000/balance/transfer
```
## Формат данных
- Идентификатор пользователя - положительное целое число
- Денежная сумма - вещественное число, с не более чем двумя знаками после запятой
  
## Структура проекта
Пытался структурировать проект в соответствии с [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
```
BalanceApp
│   README.md
│   Dockerfile
|   docker-compose.yml
|
└───bin
│   | balanceapp
|
└───_sql
│   | db.sql
│
└───cmd
│   └───balanceapp
│       │   main.go
│   
└───pkg
|   |
│   └───database
│   |   │   repo.go
|   |
│   └───handlers
│   |   │   handler.go
|   |
│   └───models
│       │   currency.go
│       │   operation.go
│       │   transaction.go
│       │   user.go
|
└───script
    │   wait-for-it.sh
```

- `bin/chatsapp` - бинарник для запуска проекта в контейнере
- `_sql/db.sql` - файл инициализации базы данных с созданием нужных таблиц
- `cmd/chatsapp/main.go` - файл для запуска приложения
- `pkg/database` - реализация работы с БД по паттерну "Репозиторий"
- `pkg/handlers` - HTTP обработчики для запросов
- `pkg/models` - описания объектов
- `script/wait-for-it.sh` - скрипт для ожидания доступности TCP хоста с портом [wait-for-it](https://github.com/vishnubob/wait-for-it)
  > Используется во время развертывания для ожидания запуска БД
