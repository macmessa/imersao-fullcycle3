package usecase

import (
	"encoding/json"
	"time"

	"github.com/macmessa/imersao-fullcycle3/codebank/domain"
	"github.com/macmessa/imersao-fullcycle3/codebank/dto"
	"github.com/macmessa/imersao-fullcycle3/codebank/infrastructure/message_broker"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
	MessageProducer       message_broker.MessageProducer
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	creditCard := u.hydrateCreditCard(transactionDto)
	ccBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	creditCard.ID = ccBalanceAndLimit.ID
	creditCard.Limit = ccBalanceAndLimit.Limit
	creditCard.Balance = ccBalanceAndLimit.Balance
	transaction := u.newTransaction(transactionDto, ccBalanceAndLimit)
	transaction.ProcessAndValidate(creditCard)
	err = u.TransactionRepository.SaveTransaction(*transaction, *creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	transactionDto.ID = transaction.ID
	transactionDto.CreatedAt = transaction.CreatedAt
	transactionJson, err := json.Marshal(transactionDto)

	if err != nil {
		return domain.Transaction{}, err
	}

	err = u.MessageProducer.Publish(string(transactionJson), "payments")

	if err != nil {
		return domain.Transaction{}, err
	}

	return *transaction, nil
}

func (u UseCaseTransaction) hydrateCreditCard(transactionDto dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDto.Name
	creditCard.Number = transactionDto.Number
	creditCard.ExpirationMonth = transactionDto.ExpirationMonth
	creditCard.ExpirationYear = transactionDto.ExpirationYear
	creditCard.CVV = transactionDto.CVV
	return creditCard
}

func (u UseCaseTransaction) newTransaction(transaction dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	newTransaction := domain.NewTransaction()
	newTransaction.CreditCardId = cc.ID
	newTransaction.Amount = transaction.Amount
	newTransaction.Store = transaction.Store
	newTransaction.Description = transaction.Description
	newTransaction.CreatedAt = time.Now()
	return newTransaction
}
