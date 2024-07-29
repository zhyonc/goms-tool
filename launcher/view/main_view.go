package view

import (
	"bytes"
	"image"
	_ "image/png"
	"launcher/assets"
	"launcher/util"

	"github.com/asaskevich/EventBus"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type mainView struct {
	name       string
	eventBus   EventBus.Bus
	mainWindow *walk.MainWindow
}

func NewMainView(eventBus EventBus.Bus, authView AuthView, homeView HomeView, menuView MenuView) MainView {
	// Create a new main window
	v := &mainView{
		name:     "MainView",
		eventBus: eventBus,
	}
	// Show in the screen center
	screenWidth, screenHeight := util.GetScreenResolution()
	posX := (screenWidth - launcherViewWidth) / 2
	posY := (screenHeight - launcherViewHeight) / 2
	declarative.MainWindow{
		AssignTo: &v.mainWindow,
		Title:    launcherViewTitle,
		Size:     declarative.Size{Width: launcherViewWidth, Height: launcherViewHeight},
		Layout:   declarative.Grid{Columns: 3},
		Bounds:   declarative.Rectangle{X: posX, Y: posY},
		Children: []declarative.Widget{
			authView.CreateWidget(),
			homeView.CreateWidget(),
			menuView.CreateWidget(),
		},
	}.Create()
	v.mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		result := walk.MsgBox(v.mainWindow, "Confirm", checkExit, walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
		if result == walk.DlgCmdYes {
			return
		}
		*canceled = true
	})
	// Load icon from a file
	img, _, err := image.Decode(bytes.NewReader(assets.IconData))
	if err != nil {
		v.ShowMessage("", err.Error(), walk.MsgBoxIconError)
		return nil
	}
	icon, err := walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		v.ShowMessage("", err.Error(), walk.MsgBoxIconError)
		return nil
	}
	err = v.mainWindow.SetIcon(icon)
	if err != nil {
		v.ShowMessage("", err.Error(), walk.MsgBoxIconError)
		return nil
	}
	// Disable minimize and/or maximize ^win.WS_MINIMIZEBOX
	win.SetWindowLong(v.mainWindow.Handle(), win.GWL_STYLE,
		win.GetWindowLong(v.mainWindow.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
	// Init subscribe
	v.CreateSubscribe()
	return v
}

// CreateSubscribe implements MainView.
func (v *mainView) CreateSubscribe() {
	v.eventBus.Subscribe(eventMainViewShowMessage, v.ShowMessage)
	v.eventBus.Subscribe(eventMainViewShowMessageResult, v.ShowMessageResult)
}

// Run implements MainView.
func (v *mainView) Run() {
	v.mainWindow.Run()
}

// ShowMessage implements MainView.
func (v *mainView) ShowMessage(title string, message string, style walk.MsgBoxStyle) {
	walk.MsgBox(v.mainWindow, title, message, style)
}

// ShowMessageResult implements MainView.
func (v *mainView) ShowMessageResult(title string, message string, eventTopic string) {
	result := walk.MsgBox(v.mainWindow, title, message, walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
	if result == walk.DlgCmdYes {
		go func() {
			errChan := make(chan error, 1)
			v.eventBus.Publish(eventTopic, errChan)
			err := <-errChan
			if err != nil {
				v.ShowMessage(v.name, err.Error(), walk.MsgBoxIconError)
			} else {
				v.ShowMessage(v.name, handleOK, walk.MsgBoxOK)
			}
		}()
	}
}
