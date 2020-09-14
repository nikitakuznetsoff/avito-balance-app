FROM golang:1.14

WORKDIR /go/src/balanceapp
COPY . /go/src/balanceapp

RUN go build -o ./bin/balanceapp ./cmd/balanceapp/
# Для возможности запуска скрипта
RUN chmod +x /go/src/balanceapp/scripts/*

EXPOSE 9000/tcp

CMD [ "/go/src/balanceapp/bin/balanceapp" ]



