package handler

import (
	"encoding/json"
	"net/http"
	"router/api"
	"router/config"
	"router/domain"
	"router/mongodb"

	"github.com/gin-gonic/gin"
)

type gameHandler struct {
	conf              *config.Config
	dbClient          *mongodb.DBClient
	loginServerXORKey []byte
}

func NewGameHandler(conf *config.Config, dbClient *mongodb.DBClient) api.GameAPI {
	h := &gameHandler{
		conf:              conf,
		dbClient:          dbClient,
		loginServerXORKey: []byte(conf.LoginServer.UDPXORKey),
	}
	return h
}

// SkipSDOAuth implements api.GameAPI.
func (h *gameHandler) SkipSDOAuth(c *gin.Context) {
	id, ok := domain.GetIDFromToken(c)
	if !ok {
		c.String(http.StatusBadRequest, domain.ErrBadRequest)
		return
	}
	var req api.SkipSDOAuthRequest
	req.AccountID = id
	content, _ := json.Marshal(req)
	msg := api.NewMessage(c.ClientIP(), api.SkipSDOAuth, content)
	err := msg.Send(h.conf.LoginServer.UDPAddr, h.loginServerXORKey)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, msg.Status)
}

// KickGameClient implements api.GameAPI.
func (h *gameHandler) KickGameClient(c *gin.Context) {
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
	msg := api.NewMessage(c.ClientIP(), api.KickGameClient, nil)
	err := msg.Send(h.conf.LoginServer.UDPAddr, h.loginServerXORKey)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, msg.Status)
}
