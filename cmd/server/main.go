package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/sylph4/entain-task/internal/record/model"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	internalHttp "github.com/sylph4/entain-task/internal/record/http"
	"github.com/sylph4/entain-task/internal/record/middleware"
	"github.com/sylph4/entain-task/internal/record/repository"
	"github.com/sylph4/entain-task/internal/record/service"
	"github.com/sylph4/entain-task/storage"
)

func main() {
	db, err := storage.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	validate := validator.New()
	model.RegisterValidators(validate)

	transactionRepository := repository.NewTransactionRepository(db)
	userRepository := repository.NewUserRepository(db)
	processRecordService := service.NewProcessRecordService(transactionRepository, userRepository, db)
	recordHandler := internalHttp.NewProcessRecordHandler(processRecordService, validate)

	stm := middleware.SourceType{}
	stm.Populate([]string{"game", "server", "payment"})

	r := mux.NewRouter()
	r.HandleFunc("/process-record", recordHandler.ProcessRecord).Methods("POST").
		Headers("Content-Type", "application/json")

	r.Use(stm.Middleware)
	r.MethodNotAllowedHandler = middleware.MethodNotAllowedHandler()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("API listening at port :8080")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("shutting down the http server")
		}
	}()

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
