package main

import (
	"balanceapp/pkg/database"
	"balanceapp/pkg/handlers"
	"context"
	"log"

	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	dsn := "root:pass@tcp(localhost:3306)/balanceapp?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	repo := database.NewRepository(db, context.Background())
	handler := handlers.Handler{Repo: repo}

	router := mux.NewRouter()
	router.HandleFunc("/balance", handler.GetUserBalance).Methods("POST")
	router.HandleFunc("/balance/withdraw", handler.WithdrawMoney).Methods("POST")
	router.HandleFunc("/balance/deposit", handler.DepositMoney).Methods("POST")
	router.HandleFunc("/balance/transfer", handler.Transfer).Methods("POST")

	address := ":9000"
	fmt.Printf("Starting server on port %v\n", address)
	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err)
	}
}