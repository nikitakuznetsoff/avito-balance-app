# avito-balance-app
## Запуск
Запуск осуществуляется с помощью `docker-compose up` из корня приложения
## Описание проекта
Основные моменты:
- Проект запускает на порту `:9000`
- В качестве HTTP роутера используется `https://github.com/gorilla/mux`
- В качестве хранилища используется база MySQL (v. 8.0.21)
  + База данных содержит 1 таблицу `users` c типом хранения данных InnoDB
  + Структуру таблиц можно посмотреть в файле инициализации `_sql/db.sql`
  
## Описание API
  
## Структура проекта
Пытался структурировать проект в соответствии с `https://github.com/golang-standards/project-layout`
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

- ```bin/chatsapp``` - бинарник для запуска проекта в контейнере
- ```_sql/db.sql``` - файл инициализации базы данных с созданием нужных таблиц
- ```cmd/chatsapp/main.go``` - файл для запуска приложения
- ```pkg/database``` - реализация работы с БД по паттерну "Репозиторий"
- ```pkg/handlers``` - HTTP обработчики для запросов
- ```pkg/models``` - описания объектов
- ```script/wait-for-it.sh``` - скрипт для ожидания доступности TCP хоста с портом ```https://github.com/vishnubob/wait-for-it```
  > Используется во время развертывания для ожидания запуска БД
