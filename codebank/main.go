package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/macmessa/imersao-fullcycle3/codebank/domain"
	"github.com/macmessa/imersao-fullcycle3/codebank/infrastructure/repository"
	"github.com/macmessa/imersao-fullcycle3/codebank/usecase"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Marco"
	cc.ExpirationYear = 2021
	cc.ExpirationMonth = 7
	cc.CVV = 321
	cc.Limit = 20000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"host.docker.internal",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	println(err)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}
