package view

import (
	"launcher/controller"

	"github.com/asaskevich/EventBus"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type menuView struct {
	name              string
	eventBus          EventBus.Bus
	controller        controller.MenuController
	skipSDOAuthButton *walk.PushButton
	kickGameButton    *walk.PushButton
}

func NewMenuView(eventBus EventBus.Bus, controller controller.MenuController) MenuView {
	v := &menuView{
		name:       "MenuView",
		eventBus:   eventBus,
		controller: controller,
	}
	v.CreateSubscribe()
	return v
}

// CreateSubscribe implements MenuView.
func (v *menuView) CreateSubscribe() {
	v.eventBus.Subscribe(eventMenuViewSetVisible, v.SetVisible)
	v.eventBus.Subscribe(eventMenuViewUnpatchHosts, v.UnpatchHosts)
}

// CreateWidget implements GameView.
func (v *menuView) CreateWidget() declarative.Widget {
	return declarative.GroupBox{
		Title:  menuGroupBoxTitlte,
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.PushButton{
				Text:      startGameButtonText,
				OnClicked: v.onStartGameButtonClicked,
			},
			declarative.PushButton{
				AssignTo:  &v.skipSDOAuthButton,
				Text:      skipSDOAuthButtonText,
				OnClicked: v.onSkipSDOAuthButtonClicked,
				Visible:   false,
			},

			declarative.PushButton{
				Text:      setGameDirButtonText,
				OnClicked: v.onSetGamePathButtonClicked,
			},
			declarative.PushButton{
				Text:      patchGameButtonText,
				OnClicked: v.onPatchGameButtonClicked,
			},
			declarative.PushButton{
				Text:      patchHostsButtonText,
				OnClicked: v.onPatchHostsButtonClicked,
			},
			declarative.PushButton{
				AssignTo:  &v.kickGameButton,
				Text:      kickGameButtonText,
				OnClicked: v.onKickGameButtonClicked,
				Visible:   false,
			},
			declarative.PushButton{
				Text:      exitGameButtonText,
				OnClicked: v.onExitGameButtonClicked,
			},
		},
	}
}

// SetVisible implements MenuView.
func (v *menuView) SetVisible(visible bool) {
	v.skipSDOAuthButton.SetVisible(visible)
	v.kickGameButton.SetVisible(visible)
}

func (v *menuView) onStartGameButtonClicked() {
	go func() {
		err := v.controller.StartGame()
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
	}()
	go v.eventBus.Publish(eventMainViewShowMessage, v.name, startingGame, walk.MsgBoxOK)
}

func (v *menuView) onSkipSDOAuthButtonClicked() {
	msg, err := v.controller.SkipSDOAuth()
	if err != nil {
		v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
		return
	}
	v.eventBus.Publish(eventMainViewShowMessage, v.name, msg, walk.MsgBoxOK)
}

func (v *menuView) onKickGameButtonClicked() {
	msg, err := v.controller.KickGame()
	if err != nil {
		v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
		return
	}
	v.eventBus.Publish(eventMainViewShowMessage, v.name, msg, walk.MsgBoxOK)
}

func (v *menuView) onSetGamePathButtonClicked() {
	dlg := new(walk.FileDialog)
	dlg.Title = "Currect directory is " + v.controller.GetCurrentGameDir()
	dlg.Filter = "All Files (*.*)|*.*"
	if ok, err := dlg.ShowBrowseFolder(nil); err == nil && ok {
		dir := dlg.FilePath
		if dir == "" {
			return
		}
		err := v.controller.SetGameDir(dir)
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
		v.eventBus.Publish(eventMainViewShowMessage, v.name, setGameDirOK, walk.MsgBoxOK)
	}
}

func (v *menuView) onPatchGameButtonClicked() {
	err := v.controller.PatchGame()
	if err != nil {
		v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
		return
	}
	v.eventBus.Publish(eventMainViewShowMessage, v.name, patchGameOK, walk.MsgBoxOK)
}

func (v *menuView) onPatchHostsButtonClicked() {
	ok := v.controller.CheckHosts()
	if ok {
		v.eventBus.Publish(eventMainViewShowMessageResult, "Confirm", errDuplicateHosts, eventMenuViewUnpatchHosts)
		return
	}
	err := v.controller.PatchHosts()
	if err != nil {
		v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
		return
	}
	v.eventBus.Publish(eventMainViewShowMessage, v.name, patchHostsOK, walk.MsgBoxOK)
}

// UnpatchHosts implements MenuView.
func (v *menuView) UnpatchHosts(errChan chan error) {
	errChan <- v.controller.UnpatchHosts()
}

func (v *menuView) onExitGameButtonClicked() {
	go func() {
		err := v.controller.ExitGmae()
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
	}()
	go v.eventBus.Publish(eventMainViewShowMessage, v.name, exitingGame, walk.MsgBoxOK)
}
