package queue

import (
	"log"
	"sync"
)

type Queue struct {
	sync.Mutex
	balances   map[int64]int64
	userQueues map[int64][]int64
}

func NewQueue() *Queue {
	return &Queue{
		balances:   make(map[int64]int64),
		userQueues: make(map[int64][]int64),
	}
}

func (q *Queue) AddTx(userId, amount int64) {
	q.Lock()
	defer q.Unlock()

	if userQueue, ok := q.userQueues[userId]; ok {
		q.userQueues[userId] = append(userQueue, amount)
	} else {
		q.userQueues[userId] = make([]int64, 1, 10)
		q.userQueues[userId][0] = amount
	}

	log.Printf("Added tx to queue: user %d amount %d", userId, amount)
}

func (q *Queue) GetTxs() []int64 {
	q.Lock()
	defer q.Unlock()

	res := make([]int64, 0, len(q.userQueues)*2)

	for k, v := range q.userQueues {
		res = append(res, k, v[0])

		log.Printf("Tx passed to worker: user %d amount %d", k, v[0])

		if len(v) == 1 {
			delete(q.userQueues, k)
		} else {
			q.userQueues[k] = q.userQueues[k][1:]
		}
	}

	return res
}

func (q *Queue) AddToBalance(userId int64, balance int64) {
	q.Lock()
	defer q.Unlock()

	q.balances[userId] += balance
}

func (q *Queue) GetBalance(userId int64) (int64, bool) {
	q.Lock()
	defer q.Unlock()

	b, ok := q.balances[userId]

	return b, ok
}
