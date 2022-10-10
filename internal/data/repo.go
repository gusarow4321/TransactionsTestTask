package data

import (
	"TransactionsTestTask/internal/pkg/ent/user"
	"context"
)

type UserRepo interface {
	UpdateBalance(context.Context, int64, int64) (int64, error)
	GetBalance(context.Context, int64) (int64, error)
}

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data) UserRepo {
	return &userRepo{
		data: data,
	}
}

func (r *userRepo) UpdateBalance(ctx context.Context, userId int64, inc int64) (int64, error) {
	u, err := r.data.db.User.UpdateOneID(userId).AddBalance(inc).Save(ctx)
	if err != nil {
		return 0, err
	}
	return u.Balance, nil
}

func (r *userRepo) GetBalance(ctx context.Context, userId int64) (int64, error) {
	u, err := r.data.db.User.Query().Where(user.ID(userId)).Only(ctx)
	if err != nil {
		return 0, err
	}
	return u.Balance, nil
}
