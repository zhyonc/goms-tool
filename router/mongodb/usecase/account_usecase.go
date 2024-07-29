package usecase

import (
	"context"
	"router/mongodb/model"
	"router/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountUsecase struct {
	baseUsecase
}

func NewAccountUsecase(db *mongo.Database) AccountUsecase {
	coll := db.Collection(accountColl)
	u := &accountUsecase{
		baseUsecase: NewBaseUsecase("AccountUsecase", coll),
	}
	return u
}

// FindAccountByUsername implements AccountUsecase.
func (u *accountUsecase) FindAccountByUsername(ctx context.Context, username string) *model.Account {
	filter := bson.D{{Key: "username", Value: username}}
	account := &model.Account{}
	ok := u.baseUsecase.FindOne(ctx, filter, account)
	if !ok {
		return nil
	}
	return account
}

// CreateNewAccount implements AccountUsecase.
func (u *accountUsecase) CreateNewAccount(ctx context.Context, account *model.Account) bool {
	account.RegisterDate = util.DBTime2Local(time.Now())
	return u.baseUsecase.InsertOne(ctx, account)
}

// UpdatePassword implements AccountUsecase.
func (u *accountUsecase) UpdatePassword(ctx context.Context, accountID uint32, isSecondPassword bool, passwrod string) bool {
	filter := bson.D{{Key: "_id", Value: accountID}}
	var update bson.M
	if isSecondPassword {
		update = bson.M{
			"second_password": passwrod,
			"update_date":     util.DBTime2Local(time.Now()),
		}
	} else {
		update = bson.M{
			"password":    passwrod,
			"update_date": util.DBTime2Local(time.Now()),
		}
	}
	return u.baseUsecase.UpdateOne(ctx, filter, update)
}

// UpdateLoginDate implements AccountUsecase.
func (u *accountUsecase) UpdateLoginDate(ctx context.Context, accountID uint32) {
	filter := bson.D{{Key: "_id", Value: accountID}}
	update := bson.M{"login_date": util.DBTime2Local(time.Now())}
	_ = u.baseUsecase.UpdateOne(ctx, filter, update)
}

// FindAccountByID implements AccountUsecase.
func (u *accountUsecase) FindAccountByID(ctx context.Context, accountID uint32) *model.Account {
	filter := bson.D{{Key: "_id", Value: accountID}}
	account := &model.Account{}
	ok := u.baseUsecase.FindOne(ctx, filter, account)
	if !ok {
		return nil
	}
	return account
}
