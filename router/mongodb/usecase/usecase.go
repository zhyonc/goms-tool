package usecase

import (
	"context"
	"router/mongodb/model"
)

// db collection name
const (
	counterColl string = "counter"
	accountColl string = "accounts"
)

type CounterUsecase interface {
	GetAccountID(ctx context.Context) uint32
}

type AccountUsecase interface {
	FindAccountByUsername(ctx context.Context, username string) *model.Account
	CreateNewAccount(ctx context.Context, account *model.Account) bool
	UpdatePassword(ctx context.Context, accountID uint32, isSecondPassword bool, passwrod string) bool
	FindAccountByID(ctx context.Context, accountID uint32) *model.Account
}
