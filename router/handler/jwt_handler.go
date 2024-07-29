package handler

import (
	"log/slog"
	"net/http"
	"router/api"
	"router/config"
	"router/domain"
	"router/mongodb"
	"router/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type jwtHandler struct {
	conf     *config.Config
	dbClient *mongodb.DBClient
}

func NewJWTHandler(conf *config.Config, dbClient *mongodb.DBClient) api.JWTAPI {
	h := &jwtHandler{
		conf:     conf,
		dbClient: dbClient,
	}
	return h
}

func (h *jwtHandler) createAccessToken(id string) (string, error) {
	accessClaims := domain.NewAccessClaims(
		id,
		h.conf.Token.Issuer,
		h.conf.Token.Audience,
		h.conf.Token.AccessTokenExpireHour,
	)
	return domain.CreateToken(accessClaims, h.conf.Token.AccessTokenSignKey)
}

func (h *jwtHandler) createRefreshToken(id string) (string, error) {
	refreshClaims := domain.NewRefreshClaims(
		id,
		h.conf.Token.RefreshTokenExpireHour,
	)
	return domain.CreateToken(refreshClaims, h.conf.Token.RefreshTokenSignKey)
}

func (h *jwtHandler) renewToken(id string) *domain.NewTokenResponse {
	accessToken, err := h.createAccessToken(id)
	if err != nil {
		slog.Error("Failed to create access token", "err", err)
		return nil
	}
	refreshToken, err := h.createRefreshToken(id)
	if err != nil {
		slog.Error("Failed to create refresh token", "err", err)
		return nil
	}
	response := &domain.NewTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response
}

// Login implements api.JWTAPI.
func (h *jwtHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByUsername(c, req.Username)
	if account == nil {
		c.String(http.StatusNotFound, domain.ErrNotFound, req.Username)
		return
	}
	// Check password
	ok := util.ComparePassword(h.conf.LoginServer.EnableBcryptPassword,
		account.Password, req.Password)
	if !ok {
		c.String(http.StatusUnauthorized, domain.ErrIncorrectPassword)
		return
	}
	response := h.renewToken(strconv.Itoa(int(account.ID)))
	if response == nil {
		c.String(http.StatusInternalServerError, domain.ErrServerBusy)
		return
	}
	c.JSON(http.StatusOK, response)
}

// Renew implements api.JWTAPI.
func (h *jwtHandler) Renew(c *gin.Context) {
	req := &domain.RenewTokenRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	claims, err := domain.ParseToken(req.RefreshToken, h.conf.Token.RefreshTokenSignKey)
	if err != nil {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	accountID, err := strconv.Atoi(claims.ID)
	if err != nil {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	account := h.dbClient.AccountUsecase.FindAccountByID(c, uint32(accountID))
	if account == nil {
		c.String(http.StatusNotFound, domain.ErrNotFound, "Account")
		return
	}
	// reply token
	response := h.renewToken(strconv.Itoa(int(account.ID)))
	if response == nil {
		c.String(http.StatusInternalServerError, domain.ErrServerBusy)
		return
	}
	c.JSON(http.StatusOK, response)
}
