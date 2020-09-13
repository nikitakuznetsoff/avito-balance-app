package database

import (
	"balanceapp/pkg/models"
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	DB	*sql.DB
	ctx context.Context
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// Метод для получения баланса пользователя с созданием транзакции
func (repo *Repository) GetBalance(user *models.User) (float64, error){
	tx, err := repo.DB.BeginTx(repo.ctx, nil)
	if err != nil {
		return -1, err
	}

	balance, err := repo.getBalance(user.ID, tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return -1, rollbackErr
		}
		return -1, err
	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return balance, nil
}

// Метод для запроса баланса
func (repo *Repository) getBalance(userID int64, tx *sql.Tx) (float64, error){
	var balance float64 = 0
	err := tx.
		QueryRow("SELECT balance FROM users WHERE id = ?", userID).
		Scan(&balance)
	if err != nil {
		return -1, err
	}
	return balance, nil
}

// Метод снятия средств во счета пользователя с созданием транзакции
func (repo *Repository) Get(op *models.Operation) (int64, error) {
	tx, err := repo.DB.BeginTx(repo.ctx, nil)

	rowsAffected, err := repo.get(op, tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return -1, rollbackErr
		}
		return -1, err
	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return rowsAffected, nil
}

// Метод для запроса на снятия средств
func (repo *Repository) get(op *models.Operation, tx *sql.Tx) (int64, error) {
	balance, err := repo.getBalance(op.UserID, tx)
	if err != nil {
		return -1, err
	}

	if balance < op.Value {
		return -1, errors.New("user doesn't have much money")
	}

	result, err := tx.Exec(
		"UPDATE users SET balance = ? WHERE id = ?",
		balance - op.Value, op.UserID,
	)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Метод для начисления средств пользователю с созданием транзакции
func (repo *Repository) Set(op *models.Operation) (int64, error) {
	tx, err := repo.DB.BeginTx(repo.ctx, nil)
	if err != nil {
		return -1, err
	}

	result, err := repo.set(op, tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return -1, rollbackErr
		}
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return result, err
}


// Метод запроса на зачисления средств
func (repo *Repository) set(op *models.Operation, tx *sql.Tx) (int64, error) {
	balance, err := repo.getBalance(op.UserID, tx)
	// Создание записи о пользователе, если она отсутствует
	if err == sql.ErrNoRows {
		result, err := tx.Exec(
			"INSERT INTO users (`id`, `balance`) VALUES (?, ?)",
			op.UserID, op.Value,
		)
		if err != nil {
			return -1, err
		}
		return result.RowsAffected()
	} else if err != nil {
		return -1, err
	}

	result, err := tx.Exec(
		"UPDATE users SET balance = ? WHERE id = ?",
		balance + op.Value, op.UserID,
	)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Метод для перевода средств от одного пользователя к другому
func (repo *Repository) Transfer(tr *models.Transaction) error {
	tx, err := repo.DB.BeginTx(repo.ctx, nil)
	// Снятие средств у отправителя
	_, err = repo.get(&models.Operation{ UserID: tr.ReceiverID, Value: tr.Value }, tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	// Начисление средств получателю
	_, err = repo.set(&models.Operation{ UserID: tr.SenderID, Value: tr.Value }, tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}