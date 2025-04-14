package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"product-svc/config"
	handlers "product-svc/handlers/product"
	repository "product-svc/repository/product"
	"product-svc/routes"
	usecases "product-svc/usecases/product"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	db, err := config.InitPostgreSQL(cfg.Postgres)
	if err != nil {
		return
	}
	defer db.Close()

	routes := initDepedencies(db)

	var wg sync.WaitGroup
	wg.Add(1)
	go routes.RunGRPC(cfg.GrpcPort, &wg)

	wg.Wait()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	<-stopSignal
	log.Default().Println("Received shutdown signal, shutting down servers...")

	routes.GrpcServer.GracefulStop()
	log.Default().Println("gRPC server stopped")

	if err := routes.Listener.Close(); err != nil {
		log.Fatalf("Error closing listener: %v", err)
	}
	log.Default().Println("HTTP server stopped")

	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
	log.Default().Println("Database connection closed")
	log.Default().Println("Service stopped gracefully")
	os.Exit(0)
}

func initDepedencies(db *sql.DB) *routes.Routes {
	productRepo := repository.NewStore(db)
	productUsecase := usecases.NewProductUsecase(productRepo)
	productHandler := handlers.NewHandler(productUsecase)

	// add handlers here
	return &routes.Routes{
		ProductHandler: productHandler,
	}
}
