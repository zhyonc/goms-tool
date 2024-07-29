package view

import (
	"launcher/controller"
	"launcher/util"

	"github.com/asaskevich/EventBus"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type homeView struct {
	name                     string
	eventBus                 EventBus.Bus
	controller               controller.HomeController
	tagWidget                *walk.TabWidget
	usernameLabel            *walk.Label
	isForeverBannedLabel     *walk.Label
	cashPointLabel           *walk.Label
	maplePointLabel          *walk.Label
	isSecondPasswordCheckBox *walk.CheckBox
	passwordEdit             *walk.LineEdit
	newPasswordEdit          *walk.LineEdit
}

func NewHomeView(eventBus EventBus.Bus, controller controller.HomeController) HomeView {
	v := &homeView{
		name:       "HomeView",
		eventBus:   eventBus,
		controller: controller,
	}
	v.CreateSubscribe()
	return v
}

// CreateSubscribe implements HomeView.
func (v *homeView) CreateSubscribe() {
	v.eventBus.Subscribe(eventHomeViewSetVisible, v.SetVisible)
	v.eventBus.Subscribe(eventHomeViewAccess, v.Access)
}

// SetVisible implements HomeView.
func (v *homeView) SetVisible(visible bool) {
	v.tagWidget.SetVisible(visible)
}

// Access implements HomeView.
func (v *homeView) Access(ok chan bool) {
	ok <- v.updateAccount() == nil
}

func (v *homeView) updateAccount() error {
	texts, err := v.controller.GetAccount()
	if err != nil {
		return err
	}
	v.usernameLabel.SetText(texts[0])
	v.isForeverBannedLabel.SetText(texts[1])
	v.cashPointLabel.SetText(texts[2])
	v.maplePointLabel.SetText(texts[3])
	return nil
}

// CreateWidget implements HomeView.
func (v *homeView) CreateWidget() declarative.Widget {
	return declarative.TabWidget{
		AssignTo: &v.tagWidget,
		Visible:  false,
		Pages: []declarative.TabPage{
			v.createAccountTabPage(),
			v.createPasswordTabPage(),
			v.createLogoutTabPage(),
		},
		OnCurrentIndexChanged: v.onCurrentIndexChanged,
	}
}

func (v *homeView) onCurrentIndexChanged() {
	if v.tagWidget.CurrentIndex() == 0 {
		err := v.updateAccount()
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err, walk.MsgBoxIconError)
		}
	}
	if v.tagWidget.CurrentIndex() == 2 {
		// logout index is 2 now, there is no other good function to trigger logout
		_ = v.tagWidget.SetCurrentIndex(0)
		err := v.controller.Logout()
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
		v.SetVisible(false)
		v.eventBus.Publish(eventAuthViewSetVisible, true)
		v.eventBus.Publish(eventMenuViewSetVisible, false)
	}
}

func (v *homeView) createAccountTabPage() declarative.TabPage {
	return declarative.TabPage{
		Title:  accountTabPageTitle,
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.GroupBox{
				Layout: declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.Label{Text: usernameLabelText},
					declarative.Label{AssignTo: &v.usernameLabel},
					declarative.Label{Text: isForeverBannedLabelText},
					declarative.Label{AssignTo: &v.isForeverBannedLabel},
					declarative.Label{Text: cashPointLabelText},
					declarative.Label{AssignTo: &v.cashPointLabel},
					declarative.Label{Text: maplePointLabelText},
					declarative.Label{AssignTo: &v.maplePointLabel},
				},
			},
		},
	}
}

func (v *homeView) createPasswordTabPage() declarative.TabPage {
	return declarative.TabPage{
		Title:  passwordTabPageTitle,
		Layout: declarative.Flow{},
		Children: []declarative.Widget{
			declarative.GroupBox{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.CheckBox{
						AssignTo:         &v.isSecondPasswordCheckBox,
						Text:             isSecondPassword,
						OnCheckedChanged: v.onIsSecondPasswordChecked,
					},
					declarative.Label{Text: passwordLabelText},
					declarative.LineEdit{
						AssignTo:     &v.passwordEdit,
						MaxSize:      declarative.Size{Width: textEditMaxWidth, Height: textEditMaxHeight},
						PasswordMode: true,
					},
					declarative.Label{Text: newPasswordLabelText},
					declarative.LineEdit{
						AssignTo:     &v.newPasswordEdit,
						MaxSize:      declarative.Size{Width: textEditMaxWidth, Height: textEditMaxHeight},
						PasswordMode: true,
					},
					declarative.PushButton{
						Text:      submitButtonText,
						OnClicked: v.onPasswordButtonClicked,
					},
					declarative.PushButton{
						Text:      clearButtonText,
						OnClicked: v.onClearButtonClicked,
					},
				},
			},
		},
	}
}

func (v *homeView) onIsSecondPasswordChecked() {
	v.passwordEdit.SetText("")
	v.newPasswordEdit.SetText("")
	if v.isSecondPasswordCheckBox.Checked() {
		v.passwordEdit.SetMaxLength(secondPasswordMaxLength)
		v.newPasswordEdit.SetMaxLength(secondPasswordMaxLength)
	} else {
		v.passwordEdit.SetMaxLength(passwordMaxLength)
		v.newPasswordEdit.SetMaxLength(passwordMaxLength)
	}
}

func (v *homeView) onPasswordButtonClicked() {
	if v.isSecondPasswordCheckBox.Checked() {
		if util.IsLessThanMinLength(secondPasswordMinLength, v.passwordEdit.Text(), v.newPasswordEdit.Text()) {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconWarning)
			return
		}
	} else if util.IsLessThanMinLength(passwordMinLength, v.passwordEdit.Text(), v.newPasswordEdit.Text()) {
		v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconWarning)
		return
	}
	err := v.controller.UpdatePassword(v.isSecondPasswordCheckBox.Checked(), v.passwordEdit.Text(), v.newPasswordEdit.Text())
	if err != nil {
		v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
		return
	}
	v.eventBus.Publish(eventMainViewShowMessage, v.name, updatePasswordOK, walk.MsgBoxOK)
}

func (v *homeView) onClearButtonClicked() {
	v.passwordEdit.SetText("")
	v.newPasswordEdit.SetText("")
}

func (v *homeView) createLogoutTabPage() declarative.TabPage {
	return declarative.TabPage{
		Title: logoutTabPageTitle,
	}
}
