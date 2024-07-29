package controller

type AuthController interface {
	Renew() bool
	Login(username string, password string) error
	Signup(username string, password string, secondPassword string) error
	Retrieve(username string, secondPassword string, newPassword string) error
}

type HomeController interface {
	GetAccount() ([]string, error)
	UpdatePassword(isSecondPassword bool, password string, newPassword string) error
	Logout() error
}

type MenuController interface {
	StartGame() error
	SkipSDOAuth() (string, error)
	KickGame() (string, error)
	SetGameDir(dir string) error
	GetCurrentGameDir() string
	PatchGame() error
	CheckHosts() bool
	PatchHosts() error
	UnpatchHosts() error
	ExitGmae() error
}
