package controller

import (
	"encoding/json"
	"fmt"
	"launcher/api"
	"launcher/config"
	"launcher/domain"
	"net/http"
)

type homeController struct {
	conf *config.Config
}

func NewHomeController(conf *config.Config) HomeController {
	c := &homeController{
		conf: conf,
	}
	return c
}

// GetAccount implements HomeController.
func (c *homeController) GetAccount() ([]string, error) {
	bytes, err := api.SendRequest(http.MethodGet, api.PathAuthAccount, nil, c.conf.AccessToken)
	if err != nil {
		return nil, err
	}
	resp := &domain.GetAccountResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	texts := make([]string, 0)
	texts = append(texts, resp.Username)
	texts = append(texts, fmt.Sprintf("%v", resp.IsForeverBanned))
	texts = append(texts, fmt.Sprintf("%v", resp.CashPoint))
	texts = append(texts, fmt.Sprintf("%v", resp.MaplePoint))
	return texts, nil
}

// UpdatePassword implements HomeController.
func (c *homeController) UpdatePassword(isSecondPassword bool, password string, newPassword string) error {
	var err error
	if isSecondPassword {
		req := &domain.UpdateSecondPasswordRequest{
			OldPassword: password,
			NewPassword: newPassword,
		}
		_, err = api.SendRequest(http.MethodPatch, api.PathAuthAccountSecondPassword, req, c.conf.AccessToken)
	} else {
		req := &domain.UpdatePasswordRequest{
			OldPassword: password,
			NewPassword: newPassword,
		}
		_, err = api.SendRequest(http.MethodPatch, api.PathAuthAccountPassword, req, c.conf.AccessToken)
	}
	return err
}

// Logout implements HomeController.
func (c *homeController) Logout() error {
	c.conf.AccessToken = ""
	c.conf.RefreshToken = ""
	return config.Save(c.conf)
}
