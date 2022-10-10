package data

import (
	"TransactionsTestTask/internal/pkg/ent"
	"context"
	"log"

	_ "github.com/lib/pq"
)

type Data struct {
	db *ent.Client
}

func NewData(conf string) (*Data, func(), error) {
	client, err := ent.Open("postgres", conf)
	if err != nil {
		log.Printf("failed opening connection to database: %v", err)
		return nil, nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Printf("failed creating schema resources: %v", err)
		return nil, nil, err
	}
	_, err = client.User.Create().SetID(1).SetBalance(0).Save(context.Background())
	if err != nil {
		log.Printf("user creating failing: %v", err)
	}
	_, err = client.User.UpdateOneID(1).SetBalance(100).Save(context.Background())
	if err != nil {
		log.Printf("user updating failing: %v", err)
	}

	d := &Data{
		db: client,
	}

	cleanup := func() {
		if err := d.db.Close(); err != nil {
			log.Printf("failed closing the data resources: %v", err)
		}
	}
	return d, cleanup, nil
}
