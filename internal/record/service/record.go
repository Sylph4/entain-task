package service

import (
	"errors"
	"log"
	"time"

	"github.com/sylph4/entain-task/internal/record/model"
	"github.com/sylph4/entain-task/internal/record/repository"
	"github.com/sylph4/entain-task/storage/postgres"
	"gopkg.in/gorp.v1"
)

type IProcessRecordService interface {
	CancelOddRecords()
	ProcessRecord(request *model.ProcessRecordRequest) error
}

type ProcessRecordService struct {
	transactionRepository repository.ITransactionRepository
	userRepository        *repository.UserRepository
	db                    *gorp.DbMap
}

func NewProcessRecordService(
	recordsRepository repository.ITransactionRepository,
	userRepository *repository.UserRepository, db *gorp.DbMap) *ProcessRecordService {
	processRecordService := &ProcessRecordService{
		transactionRepository: recordsRepository, userRepository: userRepository, db: db}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			processRecordService.CancelOddRecords()
		}
	}()

	return processRecordService
}

func (s *ProcessRecordService) CancelOddRecords() {
	log.Println("Odd records process started at: ", time.Now().Format("2006-01-02T15:04:05"))
	transactions, err := s.transactionRepository.SelectLastTenOddTransactions()
	if err != nil {
		log.Println("selecting odd transactions failed")
	}

	if len(transactions) == 0 {
		log.Printf("no odd transactions to process")

		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("could not start a transaction")

		return
	}

	for i := range transactions {
		user, err := s.userRepository.SelectUserByID(tx, transactions[i].UserID)
		if err != nil {
			log.Printf("selecting user by id %s failed \n", transactions[i].UserID)
		}

		if user == nil {
			log.Printf("user %s does not exist \n", transactions[i].UserID)

			continue
		}

		if transactions[i].State == model.StateWin {
			if (user.Balance - transactions[i].Amount) < 0 {
				log.Printf("account balance can't be in a negative value, skipping balance correction")

				continue
			}

			user.Balance -= transactions[i].Amount
		} else {
			user.Balance += transactions[i].Amount
		}

		err = s.userRepository.Update(tx, *user)
		if err != nil {
			log.Printf("could not update user balance" + err.Error())

			continue
		} else {
			transactions[i].IsCanceled = true
		}
	}

	err = s.transactionRepository.UpdateBulk(tx, transactions)
	if err != nil {
		log.Printf("could not update transactions" + err.Error())

		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("could not commit transaction")
	}

	log.Println("Odd records process ended at: ", time.Now().Format("2006-01-02T15:04:05"))
}

func (s *ProcessRecordService) ProcessRecord(request *model.ProcessRecordRequest) error {
	existingTransaction, err := s.transactionRepository.SelectTransactionByID(request.TransactionID)
	if err != nil {
		return err
	}

	if existingTransaction != nil {
		return errors.New("transaction already processed")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	user, err := s.userRepository.SelectUserByID(tx, request.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user does not exist")
	}

	transaction := &postgres.Transaction{
		ID:          request.TransactionID,
		UserID:      request.UserID,
		State:       request.State,
		Amount:      request.Amount,
		ProcessedAt: time.Now().UTC(),
	}

	if transaction.State == model.StateWin {
		user.Balance += request.Amount
	} else {
		if (user.Balance - request.Amount) < 0 {
			return errors.New("account balance can't be in a negative value")
		}

		user.Balance -= request.Amount
	}

	err = s.transactionRepository.Insert(tx, *transaction)
	if err != nil {
		return err
	}

	err = s.userRepository.Update(tx, *user)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.New("committing transaction failed")
	}

	return nil
}
