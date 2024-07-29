package config

type RouterConfig struct {
	Addr           string
	EnableStaticFS bool
	StaticFilePath string
}

func defaultRouterConfig() RouterConfig {
	return RouterConfig{
		Addr:           "127.0.0.1:80",
		EnableStaticFS: true,
		StaticFilePath: "./assets",
	}
}
