package handler

import (
	"fmt"
	"net/http"
	"router/api"
	"router/config"
	"router/domain"
	"router/mongodb"
	"router/mongodb/model"
	"router/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountHadnler struct {
	conf     *config.Config
	dbClient *mongodb.DBClient
}

func NewAccountHandler(conf *config.Config, dbClient *mongodb.DBClient) api.AccountAPI {
	h := &accountHadnler{
		conf:     conf,
		dbClient: dbClient,
	}
	return h
}

// Signup implements api.AccountAPI.
func (h *accountHadnler) Signup(c *gin.Context) {
	var req domain.SignupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByUsername(c, req.Username)
	if account != nil {
		c.String(http.StatusConflict, domain.ErrDuplicateUsername)
		return
	}
	// Enable bcrypt
	if h.conf.LoginServer.EnableBcryptPassword {
		// Password
		bPassword, err := util.Bcrypt(req.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, domain.ErrServerBusy)
			return
		}
		req.Password = bPassword
	}
	// Create new account
	var res bool = false
	h.dbClient.WithTransaction(func(ctx mongo.SessionContext) (any, error) {
		accountID := h.dbClient.CounterUsecase.GetAccountID(ctx)
		account := model.NewAccount(accountID, req.Username, req.Password, req.SecondPassword, c.ClientIP(), "")
		res = h.dbClient.AccountUsecase.CreateNewAccount(ctx, account)
		if !res {
			c.String(http.StatusInternalServerError, domain.ErrServerBusy)
			return nil, fmt.Errorf("failed to create new account")
		}
		c.String(http.StatusOK, "OK")
		return nil, nil
	})
}

// Retrieve implements api.AccountAPI.
func (h *accountHadnler) Retrieve(c *gin.Context) {
	var req domain.RetrieveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByUsername(c, req.Username)
	if account == nil {
		c.String(http.StatusNotFound, domain.ErrNotFound, "Account")
		return
	}
	if account.SecondPassword != req.SecondPassword {
		c.String(http.StatusUnauthorized, domain.ErrIncorrectSecondPassword)
		return
	}
	h.updatePasswordConfirm(c, account.ID, req.NewPassword)
}

// GetAccount implements api.AccountAPI.
func (h *accountHadnler) GetAccount(c *gin.Context) {
	id, ok := domain.GetIDFromToken(c)
	if !ok {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByID(c, id)
	if account == nil {
		c.String(http.StatusNotFound, domain.ErrNotFound, "Account")
		return
	}
	resp := domain.GetAccountResponse{
		Username: account.Username,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdatePassword implements api.AccountAPI.
func (h *accountHadnler) UpdatePassword(c *gin.Context) {
	id, ok := domain.GetIDFromToken(c)
	if !ok {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	var req domain.UpdatePasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByID(c, id)
	if account == nil {
		c.String(http.StatusNotFound, domain.ErrNotFound, "Account")
		return
	}
	// Check password
	enableBcrypt := h.conf.LoginServer.EnableBcryptPassword
	ok = util.ComparePassword(enableBcrypt, account.Password, req.OldPassword)
	if !ok {
		c.String(http.StatusUnauthorized, domain.ErrIncorrectPassword)
		return
	}
	h.updatePasswordConfirm(c, account.ID, req.NewPassword)
}

func (h *accountHadnler) updatePasswordConfirm(c *gin.Context, accountID uint32, newPassword string) {
	// Bcrypt password
	if h.conf.LoginServer.EnableBcryptPassword {
		bPassword, err := util.Bcrypt(newPassword)
		if err != nil {
			c.String(http.StatusInternalServerError, domain.ErrServerBusy)
			return
		}
		newPassword = bPassword
	}
	// UpdatePassword
	ok := h.dbClient.AccountUsecase.UpdatePassword(c, accountID, false, newPassword)
	if !ok {
		c.String(http.StatusInternalServerError, domain.ErrServerBusy)
		return
	}
	c.String(http.StatusOK, "OK")
}

// UpdateSecondPassword implements api.AccountAPI.
func (h *accountHadnler) UpdateSecondPassword(c *gin.Context) {
	id, ok := domain.GetIDFromToken(c)
	if !ok {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	var req domain.UpdateSecondPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByID(c, id)
	if account == nil {
		c.String(http.StatusNotFound, domain.ErrNotFound, "Account")
		return
	}
	// Check password
	if account.SecondPassword != req.OldPassword {
		c.String(http.StatusUnauthorized, domain.ErrIncorrectPassword)
		return
	}
	// Update password
	ok = h.dbClient.AccountUsecase.UpdatePassword(c, account.ID, true, req.NewPassword)
	if !ok {
		c.String(http.StatusInternalServerError, domain.ErrServerBusy)
		return
	}
	c.String(http.StatusOK, "OK")
}
