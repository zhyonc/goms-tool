package view

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

// Event BUS
const (
	eventMainViewShowMessage       string = "mainView:ShowMessage"
	eventMainViewShowMessageResult string = "mainView:ShowMessageResult"
	eventAuthViewSetVisible        string = "authView:SetVisible"
	eventHomeViewSetVisible        string = "homeView:SetVisible"
	eventHomeViewAccess            string = "homeView:Access"
	eventMenuViewSetVisible        string = "menuView:SetVisible"
	eventMenuViewUnpatchHosts      string = "menuView:UnpatchHosts"
)

type BaseView interface {
	CreateSubscribe()
}

type MainView interface {
	BaseView
	Run()
	ShowMessage(title string, message string, style walk.MsgBoxStyle)
	ShowMessageResult(title string, message string, eventTopic string)
}

type SubView interface {
	BaseView
	SetVisible(visible bool)
	CreateWidget() declarative.Widget
}

type AuthView interface {
	SubView
	AutoLogin()
}

type HomeView interface {
	SubView
	Access(ok chan bool)
}

type MenuView interface {
	SubView
	UnpatchHosts(errChan chan error)
}
