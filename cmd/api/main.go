package main

import (
	"AccountManagementSystem/env_helper"
	"AccountManagementSystem/internal/handlers"
	"AccountManagementSystem/internal/queue_processor"
	"AccountManagementSystem/internal/queue_producer"
	"AccountManagementSystem/internal/repository"
	"AccountManagementSystem/internal/server"
	"AccountManagementSystem/internal/services"
	"AccountManagementSystem/log_color"
	"AccountManagementSystem/log_helper"
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "AccountManagementSystem/docs"
	_ "github.com/lib/pq" // <-- REQUIRED: registers "postgres" drive
	//_ "AccountManagementSystem/internal/database/migrations" // 👈 import so migrations register
	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

// @title AccountManagementSystem
// @version     1.0
// @description API for Account Management system
func main() {

	log_helper.LogServiceInitialized("main")
	server := server.New()

	server.RegisterFiberRoutes()
	initHandler(server)
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func initHandler(server *server.FiberServer) {

	slog.Info(log_helper.LogServiceInitializing("kafka producer"))
	kafkaBroker := env_helper.ReadString("KAFKA_BROKERS", "kafka:9092")
	kafkaTopic := env_helper.ReadString("KAFKA_TOPIC", "transactions")
	kafkaGroup := env_helper.ReadString("KAFKA_GROUP", "transaction-service")
	kafkaQueue, err := queue_producer.NewKafkaQueue(kafkaBroker, kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info(log_helper.LogServiceInitialized("kafka producer"))

	slog.Info(log_color.Black("===================================="))
	slog.Info(log_helper.LogServiceInitializing("account_System"))
	accountRepo := repository.NewAccountRepo(server.DB())
	accountService := services.NewAccountService(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountService)
	server.RegisterAPIRoutes(accountHandler)
	slog.Info(log_helper.LogServiceInitialized("account_System"))

	slog.Info(log_color.Black("===================================="))
	slog.Info(log_helper.LogServiceInitializing("transaction_System"))
	transactionRepo := repository.NewTransactionRepo(server.DB())
	transactionService := services.NewTransactionService(accountRepo, transactionRepo, kafkaQueue)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	server.RegisterAPIRoutes(transactionHandler)
	slog.Info(log_helper.LogServiceInitialized("transaction_System"))

	slog.Info(log_color.Black("===================================="))
	slog.Info(log_helper.LogServiceInitializing("kafka consumer"))
	queue_processor.StartConsumer(transactionService, kafkaBroker, kafkaTopic, kafkaGroup)
	slog.Info(log_helper.LogServiceInitialized("kafka consumer"))
}

func initGoose() {
	dbstring := os.Getenv("GOOSE_DBSTRING")

	db, err := sql.Open("postgres", dbstring)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := goose.UpContext(context.Background(), db, "internal/database/migrations"); err != nil {
		log.Fatal(err)
	}
}
