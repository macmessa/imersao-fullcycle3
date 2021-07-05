package repository

import (
	"database/sql"
	"github.com/macmessa/imersao-fullcycle3/codebank/domain"
)


type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepositoryDb(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (t *TransactionRepositoryDb) SaveTransaction(transaction domain.Transaction, creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(query:`
		insert into transactions(id, credit_card_id, amount, status, description, store, created_at)
		values($1, $2, $3, $4, $5, $6, $7)
	`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)

	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		// transaction logic here
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}
