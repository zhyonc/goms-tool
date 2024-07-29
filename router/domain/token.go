package domain

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func NewAccessClaims(id string, issuer string, audience string, expire uint16) CustomClaims {
	exp := time.Now().Add(time.Hour * time.Duration(expire))
	claims := CustomClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
		},
	}
	return claims
}

func NewRefreshClaims(id string, expire uint16) CustomClaims {
	exp := time.Now().Add(time.Hour * time.Duration(expire))
	claims := CustomClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	return claims
}

func CreateToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func ParseToken(requestToken string, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(requestToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: <%v>, expected: HS256", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func GetIDFromToken(c *gin.Context) (uint32, bool) {
	temp, ok := c.Get("id") // from middleware jwt_auth
	if !ok {
		slog.Error("Can't find id from access token")
		return 0, false
	}
	idStr, ok := temp.(string)
	if !ok {
		slog.Error("Can't assert string type")
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Can't convert id type")
		return 0, false
	}
	return uint32(id), true
}
