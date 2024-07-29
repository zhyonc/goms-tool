package controller

import (
	"encoding/json"
	"launcher/api"
	"launcher/config"
	"launcher/domain"
	"net/http"
)

type authController struct {
	conf *config.Config
}

func NewAuthController(conf *config.Config) AuthController {
	c := &authController{
		conf: conf,
	}
	return c
}

func (c *authController) saveToken(accessToken string, refreshToken string) error {
	c.conf.AccessToken = accessToken
	c.conf.RefreshToken = refreshToken
	return config.Save(c.conf)
}

// Renew implements AuthController.
func (c *authController) Renew() bool {
	if c.conf.RefreshToken == "" {
		return false
	}
	req := &domain.RenewTokenRequest{
		RefreshToken: c.conf.RefreshToken,
	}
	bytes, err := api.SendRequest(http.MethodPatch, api.PathJWT, req, "")
	if err != nil {
		return false
	}
	resp := &domain.NewTokenResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return false
	}
	return c.saveToken(resp.AccessToken, resp.RefreshToken) == nil
}

// Login implements AuthController.
func (c *authController) Login(username string, password string) error {
	req := &domain.LoginRequest{
		Username: username,
		Password: password,
	}
	bytes, err := api.SendRequest(http.MethodPost, api.PathJWT, req, "")
	if err != nil {
		return err
	}
	resp := &domain.NewTokenResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return err
	}
	return c.saveToken(resp.AccessToken, resp.RefreshToken)
}

// Signup implements AuthController.
func (c *authController) Signup(username string, password string, secondPassword string) error {
	req := &domain.SignupRequest{
		Username:       username,
		Password:       password,
		SecondPassword: secondPassword,
	}
	_, err := api.SendRequest(http.MethodPost, api.PathAccount, req, "")
	if err != nil {
		return err
	}
	return nil
}

// Retrieve implements AuthController.
func (c *authController) Retrieve(username string, secondPassword string, newPassword string) error {
	req := &domain.RetrieveRequest{
		Username:       username,
		SecondPassword: secondPassword,
		NewPassword:    newPassword,
	}
	_, err := api.SendRequest(http.MethodPatch, api.PathAccountPassword, req, "")
	if err != nil {
		return err
	}
	return nil
}
