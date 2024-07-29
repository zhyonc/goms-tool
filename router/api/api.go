package api

import "github.com/gin-gonic/gin"

const (
	PathJWT                       string = "/jwt"
	PathAccount                   string = "/account"
	PathAccountPassword           string = PathAccount + "/password"
	PathAuthGroup                 string = "/auth"
	PathAuthAccount               string = "/account"
	PathAuthAccountPassword       string = "/password"
	PathAuthAccountSecondPassword string = PathAuthAccount + "/second-password"
	PathAuthGame                  string = "/game"
	PathAuthGameSkipSDOAuth       string = PathAuthGame + "/skip-sdo-auth"
	PathAuthGameKick              string = PathAuthGame + "/kick"
)

type AccountAPI interface {
	Signup(c *gin.Context)               // POST PathAccount
	Retrieve(c *gin.Context)             // PATCH PathAccountPassword
	GetAccount(c *gin.Context)           // GET PathAuthAccount
	UpdatePassword(c *gin.Context)       // PATCH PathAuthAccountPassword
	UpdateSecondPassword(c *gin.Context) // PATCH PathAuthAccountSecondPassword
}

type JWTAPI interface {
	Login(c *gin.Context) // POST PathJWT
	Renew(c *gin.Context) // PATCH PathJWT
}

type GameAPI interface {
	SkipSDOAuth(c *gin.Context)    // GET PathAuthGameSkipSDOAuth
	KickGameClient(c *gin.Context) // GET PathAuthGameKick
}
