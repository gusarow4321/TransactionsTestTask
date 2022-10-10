package service

import (
	"TransactionsTestTask/internal/data"
	"TransactionsTestTask/internal/pkg/queue"
	"TransactionsTestTask/internal/server"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Service struct {
	repo  data.UserRepo
	queue *queue.Queue
}

type transaction struct {
	UserId int64 `json:"user_id"`
	Amount int64 `json:"amount"`
}

func NewService(repo data.UserRepo, queue *queue.Queue) server.Service {
	return &Service{
		repo:  repo,
		queue: queue,
	}
}

func (s *Service) AddTx(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var tx transaction
	err := json.Unmarshal(reqBody, &tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("New tx received: user %d, amount %d", tx.UserId, tx.Amount)

	balance, ok := s.queue.GetBalance(tx.UserId)
	if !ok {
		dbBalance, err := s.repo.GetBalance(r.Context(), tx.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.queue.AddToBalance(tx.UserId, dbBalance)
		balance = dbBalance
		log.Printf("Balance updated from db")
	}

	if balance+tx.Amount < 0 {
		http.Error(w, "Low amount", http.StatusBadRequest)
		return
	}
	s.queue.AddToBalance(tx.UserId, tx.Amount)
	s.queue.AddTx(tx.UserId, tx.Amount)
	tx.Amount += balance

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
