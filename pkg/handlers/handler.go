package handlers

import (
	"balanceapp/pkg/database"
	"balanceapp/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

type Handler struct {
	Repo *database.Repository
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	currency := r.FormValue("currency")
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := &models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Проверка валидности значений
	if user.ID < 0 {
		http.Error(w, errors.New("incorrect params").Error(), http.StatusBadRequest)
		return
	}

	balance, err := h.Repo.GetBalance(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Если есть параметр другой валюты, то получаем текущий курс и изменяем значение баланса
	if currency != "" {
		url := "https://api.exchangeratesapi.io/latest?base=RUB"
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		curr := &models.Currency{}
		err = json.Unmarshal(respBody, &curr)

		currValue, ok := curr.Rates[currency]
		if !ok {
			http.Error(w, errors.New("incorrect currency").Error(), http.StatusBadRequest)
			return
		}
		balance = balance * currValue
	} else {
		currency = "RUB"
	}
	fmt.Fprintf(
		w,
		"\n[BALANCE] Balance of User(ID:%v) - %v %s\n",
		user.ID, balance, currency)
}

func (h *Handler) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	operation := &models.Operation{}
	err = json.Unmarshal(body, &operation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		Проверка валидности значений - значения положительны
		и кол-во знаков после запятой у денежной суммы меньше 3
	 */
	if 	operation.UserID < 0 ||
		operation.Value < 0  ||
		operation.Value - math.Round(operation.Value * 100) / 100 > 0 {
		http.Error(w, errors.New("incorrect params").Error(), http.StatusBadRequest)
		return
	}
	_, err = h.Repo.Get(operation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(
		w,
		"\n[WITHDRAW] %v %s was withdrawn from the account (UserID: %v)\n",
		operation.Value, "RUB", operation.UserID,
	)
}

func (h *Handler) DepositMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	operation := &models.Operation{}
	err = json.Unmarshal(body, &operation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		Проверка валидности значений - значения положительны
		и кол-во знаков после запятой у денежной суммы меньше 3
	*/
	if 	operation.UserID < 0 ||
		operation.Value < 0  ||
		operation.Value - math.Round(operation.Value * 100) / 100 > 0 {
		http.Error(w, errors.New("incorrect params").Error(), http.StatusBadRequest)
		return
	}

	_, err = h.Repo.Set(operation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(
		w,
		"\n[DEPOSIT] %v %s was deposited to the account (UserID: %v)\n",
		operation.Value, "RUB", operation.UserID,
	)
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transaction := &models.Transaction{}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		Проверка валидности значений - значения положительны
		и кол-во знаков после запятой у денежной суммы меньше 3
	*/
	if transaction.SenderID < 0 || transaction.ReceiverID < 0 ||
		transaction.Value - math.Round(transaction.Value * 100) / 100 > 0 {
		http.Error(w, errors.New("incorrect params").Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.Transfer(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(
		w,
		"\n[TRANSFER] %v %s was transferred from User(ID:%v) to User(ID:%v)\n",
		transaction.Value, "RUB", transaction.SenderID, transaction.ReceiverID,
	)
}