package main

import (
	"TransactionsTestTask/internal/data"
	"TransactionsTestTask/internal/pkg/queue"
	"TransactionsTestTask/internal/server"
	"TransactionsTestTask/internal/service"
	"TransactionsTestTask/internal/worker"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dbConf string
var httpAddr string

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	flag.StringVar(&dbConf, "db", "host=postgres user=user password=password dbname=app sslmode=disable", "Postgres conf")
	flag.StringVar(&httpAddr, "http", ":3000", "Address to listen on")
}

func newApp() (*http.Server, *worker.Worker, func(), error) {
	d, cleanup, err := data.NewData(dbConf)
	if err != nil {
		return nil, nil, nil, err
	}
	userRepo := data.NewUserRepo(d)
	q := queue.NewQueue()
	srv := service.NewService(userRepo, q)
	if err != nil {
		return nil, nil, nil, err
	}

	return server.NewServer(httpAddr, srv), worker.NewWorker(userRepo, q), cleanup, nil
}

func main() {
	flag.Parse()

	app, w, cleanup, err := newApp()
	if err != nil {
		log.Panic(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = app.ListenAndServe(); err != nil {
			log.Printf("Server Failed: %v", err)
			done <- syscall.SIGINT
		}
	}()

	workerCtx, workerCancel := context.WithCancel(context.Background())

	go func() {
		if err = w.Start(workerCtx); err != nil {
			log.Printf("Worker Failed: %v", err)
			done <- syscall.SIGINT
		}
	}()

	log.Print("Server Started")

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		workerCancel()
		cleanup()
		cancel()
	}()

	if err = app.Shutdown(ctx); err != nil {
		log.Printf("Shutdown Failed: %v", err)
	}

	log.Print("Server Stopped")
}
