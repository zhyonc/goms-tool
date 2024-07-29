package config

type TokenConfig struct {
	Issuer                 string
	Audience               string
	AccessTokenExpireHour  uint16
	RefreshTokenExpireHour uint16
	AccessTokenSignKey     string
	RefreshTokenSignKey    string
}

func defaultTokenConfig() TokenConfig {
	return TokenConfig{
		Issuer:                 "anyone",
		Audience:               "goms",
		AccessTokenExpireHour:  24,
		RefreshTokenExpireHour: 72,
		AccessTokenSignKey:     "123456",
		RefreshTokenSignKey:    "654321",
	}
}
