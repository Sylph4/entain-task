package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/sylph4/entain-task/storage/postgres"
	"gopkg.in/gorp.v1"
)

type ITransactionRepository interface {
	Insert(tx *gorp.Transaction, transaction postgres.Transaction) error
	SelectTransactionByID(id strfmt.UUID4) (*postgres.Transaction, error)
	SelectLastTenOddTransactions() ([]postgres.Transaction, error)
	UpdateBulk(tx *gorp.Transaction, transactions []postgres.Transaction) error
}

type TransactionRepository struct {
	db *gorp.DbMap
}

func NewTransactionRepository(conn *gorp.DbMap) *TransactionRepository {
	return &TransactionRepository{db: conn}
}

func (s *TransactionRepository) Insert(tx *gorp.Transaction, transaction postgres.Transaction) error {
	transaction.ProcessedAt = time.Now().UTC()
	if err := tx.Insert(&transaction); err != nil {
		return errors.New("inserting transaction record failed")
	}

	return nil
}

func (s *TransactionRepository) SelectTransactionByID(id strfmt.UUID4) (*postgres.Transaction, error) {
	transaction := &postgres.Transaction{}
	if err := s.db.SelectOne(&transaction, `
		SELECT
			*
		FROM
			transaction
		WHERE
			id = $1
	`,
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			// transaction not found
			return nil, nil
		}

		return nil, err
	}

	return transaction, nil
}

func (s *TransactionRepository) SelectLastTenOddTransactions() ([]postgres.Transaction, error) {
	var transactions []postgres.Transaction
	if _, err := s.db.Select(&transactions, `
		SELECT
			*
		FROM
			transaction
		WHERE 
			mod(amount,2) <> 0
		AND
			is_canceled = FALSE
		ORDER BY 
			processed_at DESC
		LIMIT 10
	`,
	); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *TransactionRepository) UpdateBulk(tx *gorp.Transaction, transactions []postgres.Transaction) error {
	for i := range transactions {
		_, err := tx.Update(&transactions[i])
		if err != nil {
			return err
		}
	}

	return nil
}
