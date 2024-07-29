package config

type LoginServerConfig struct {
	UDPAddr              string
	UDPXORKey            string
	EnableBcryptPassword bool
}

func defaultLoginServerConfig() LoginServerConfig {
	return LoginServerConfig{
		UDPAddr:              "127.0.0.1:8484",
		UDPXORKey:            "123456",
		EnableBcryptPassword: true,
	}
}
