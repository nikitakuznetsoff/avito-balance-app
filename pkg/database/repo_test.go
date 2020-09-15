package database

import (
	"balanceapp/pkg/models"
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Тест на получение баланса
func TestGetBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &Repository{
		DB: db,
		Ctx: context.Background(),
	}
	user := &models.User{ ID: 1 }
	var userBalance float64 = 44

	rows := sqlmock.NewRows([]string{"balance"}).AddRow(userBalance)
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE id = ?").
		WithArgs(user.ID).
		WillReturnRows(rows)
	mock.ExpectCommit()

	balance, err := repo.GetBalance(user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
	if !reflect.DeepEqual(balance, userBalance) {
		t.Errorf("results not match, want %v, have %v", userBalance, balance)
		return
	}
}

// Тест на снятие средств
func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &Repository{
		DB: db,
		Ctx: context.Background(),
	}
	op := &models.Operation{
		UserID: 1,
		Value: 100,
	}

	row := sqlmock.NewRows([]string{"balance"}).AddRow(200)
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(row)
	mock.
		ExpectExec("UPDATE users").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err = repo.Get(op)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
}

// Тест на недостаток средств на счете
func TestGet_MoneyLackError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &Repository{
		DB: db,
		Ctx: context.Background(),
	}
	op := &models.Operation{
		UserID: 1,
		Value: 100,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnError(fmt.Errorf("user (ID:1) doesn't have much money"))
	mock.ExpectRollback()

	_, err = repo.Get(op)
	if err == nil {
		t.Errorf("was expecting an error, but there was none")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
}

// Тест на зачисление баланса
func TestSet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &Repository{
		DB: db,
		Ctx: context.Background(),
	}
	op := &models.Operation{
		UserID: 1,
		Value: 100,
	}

	row := sqlmock.NewRows([]string{"balance"}).AddRow(20)
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(row)
	mock.
		ExpectExec("UPDATE users").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err = repo.Set(op)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
}

// Тест на перевод денег
func TestTransfer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &Repository{
		DB: db,
		Ctx: context.Background(),
	}
	tr := &models.Transaction{
		SenderID: 1,
		ReceiverID: 2,
		Value: 100,
	}

	rowFirst := sqlmock.NewRows([]string{"balance"}).AddRow(200)
	rowSecond := sqlmock.NewRows([]string{"balance"}).AddRow(50)
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rowFirst)
	mock.
		ExpectExec("UPDATE users").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectQuery("SELECT balance FROM users WHERE id = ?").
		WithArgs(2).
		WillReturnRows(rowSecond)
	mock.
		ExpectExec("UPDATE users").
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err = repo.Transfer(tr)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}

}