package main

import (
	"launcher/config"
	"launcher/controller"
	"launcher/view"

	"github.com/asaskevich/EventBus"
)

func main() {
	eventBus := EventBus.New()
	conf := config.Load()
	authController := controller.NewAuthController(conf)
	homeController := controller.NewHomeController(conf)
	menuController := controller.NewMenuController(conf)
	authView := view.NewAuthView(eventBus, authController)
	homeView := view.NewHomeView(eventBus, homeController)
	menuView := view.NewMenuView(eventBus, menuController)
	mainView := view.NewMainView(eventBus, authView, homeView, menuView)
	if mainView == nil {
		return
	}
	authView.AutoLogin()
	mainView.Run()
}
