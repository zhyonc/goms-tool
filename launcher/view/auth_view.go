package view

import (
	"launcher/controller"

	"github.com/asaskevich/EventBus"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type Action uint8

const (
	LoginSubmit Action = iota
	SignupSubmit
	RetrieveSubmit
	PasswordSubmit
)

type authView struct {
	name       string
	eventBus   EventBus.Bus
	controller controller.AuthController
	tabWidget  *walk.TabWidget
}

func NewAuthView(eventBus EventBus.Bus, controller controller.AuthController) AuthView {
	v := &authView{
		name:       "AuthView",
		eventBus:   eventBus,
		controller: controller,
	}
	v.CreateSubscribe()
	return v
}

// CreateSubscribe implements AuthView.
func (v *authView) CreateSubscribe() {
	v.eventBus.Subscribe(eventAuthViewSetVisible, v.SetVisible)
}

// SetVisible implements AuthView.
func (v *authView) SetVisible(visible bool) {
	v.tabWidget.SetVisible(visible)
}

// CreateWidget implements AuthView.
func (v *authView) CreateWidget() declarative.Widget {
	return declarative.TabWidget{
		AssignTo: &v.tabWidget,
		Pages: []declarative.TabPage{
			v.createLoginTabPage(),
			v.createSignupTabPage(),
			v.createRetrieveTabPage(),
		},
		Visible: true,
	}
}

// AutoLogin implements AuthView.
func (v *authView) AutoLogin() {
	ok := v.controller.Renew()
	if !ok {
		return
	}
	v.authViewAccessHomeView()
}

func (v *authView) authViewAccessHomeView() {
	ok := make(chan bool, 1)
	v.eventBus.Publish(eventHomeViewAccess, ok)
	if <-ok {
		v.SetVisible(false)
		v.eventBus.Publish(eventHomeViewSetVisible, true)
		v.eventBus.Publish(eventMenuViewSetVisible, true)
	}
}

func (v *authView) createBaseTabPage(title string, labelTexts []string, maxLengths []int, action Action) declarative.TabPage {
	editGroupBox := declarative.GroupBox{
		Title:    authGroupBoxTitlte,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{},
	}
	edits := make([]*walk.LineEdit, len(labelTexts))
	for i := 0; i < len(labelTexts); i++ {
		passwordMode := true
		if i == 0 {
			passwordMode = false
		}
		editGroupBox.Children = append(editGroupBox.Children,
			declarative.Label{Text: labelTexts[i]},
			declarative.LineEdit{
				AssignTo:     &edits[i],
				MaxLength:    maxLengths[i],
				MinSize:      declarative.Size{Width: textEditMaxWidth, Height: textEditMaxHeight},
				MaxSize:      declarative.Size{Width: textEditMaxWidth, Height: textEditMaxHeight},
				PasswordMode: passwordMode,
			},
		)
	}
	tagPage := declarative.TabPage{
		Title:  title,
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			editGroupBox,
			declarative.PushButton{
				Text: submitButtonText,
				OnClicked: func() {
					v.onSubmitButtonClicked(action, edits)
				},
			},
			declarative.PushButton{
				Text: clearButtonText,
				OnClicked: func() {
					for _, edit := range edits {
						edit.SetText("")
					}
				},
			},
		},
	}
	return tagPage
}

func (v *authView) createLoginTabPage() declarative.TabPage {
	labelTexts := []string{usernameLabelText, passwordLabelText}
	maxLengths := []int{usernameMaxLength, passwordMaxLength}
	return v.createBaseTabPage(loginTabPageTitle, labelTexts, maxLengths, LoginSubmit)
}

func (v *authView) createSignupTabPage() declarative.TabPage {
	labelTexts := []string{usernameLabelText, passwordLabelText, passwordAgainLableText, secondPasswordLableText}
	maxLengths := []int{usernameMaxLength, passwordMaxLength, passwordMaxLength, secondPasswordMaxLength}
	return v.createBaseTabPage(signupTabPageTitle, labelTexts, maxLengths, SignupSubmit)
}

func (v *authView) createRetrieveTabPage() declarative.TabPage {
	labelTexts := []string{usernameLabelText, secondPasswordLableText, passwordLabelText}
	maxLengths := []int{usernameMaxLength, secondPasswordMaxLength, passwordMaxLength}
	return v.createBaseTabPage(retrieveTabPageTitle, labelTexts, maxLengths, RetrieveSubmit)
}

func (v *authView) onSubmitButtonClicked(action Action, edits []*walk.LineEdit) {
	switch action {
	case LoginSubmit:
		if len(edits) < 1 {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconError)
			return
		}
		username := edits[0].Text()
		password := edits[1].Text()
		if len(username) < usernameMinLength || len(password) < passwordMinLength {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconError)
			return
		}
		err := v.controller.Login(username, password)
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
		v.authViewAccessHomeView()
	case SignupSubmit:
		if len(edits) < 3 {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconError)
			return
		}
		username := edits[0].Text()
		password := edits[1].Text()
		passwordAgain := edits[2].Text()
		secondPassword := edits[3].Text()
		if len(username) < usernameMinLength || len(password) < passwordMinLength || len(passwordAgain) < passwordMinLength || len(secondPassword) < secondPasswordMinLength {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconError)
			return
		}
		if password != passwordAgain {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errDifferentPassword, walk.MsgBoxIconError)
			return
		}
		err := v.controller.Signup(username, password, secondPassword)
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
		v.eventBus.Publish(eventMainViewShowMessage, v.name, signupOK, walk.MsgBoxOK)
	case RetrieveSubmit:
		if len(edits) < 2 {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconError)
			return
		}
		username := edits[0].Text()
		secondPassword := edits[1].Text()
		newPassword := edits[2].Text()
		if len(username) < usernameMinLength || len(secondPassword) < secondPasswordMinLength || len(newPassword) < passwordMinLength {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, errEditTextLength, walk.MsgBoxIconError)
			return
		}
		err := v.controller.Retrieve(username, secondPassword, newPassword)
		if err != nil {
			v.eventBus.Publish(eventMainViewShowMessage, v.name, err.Error(), walk.MsgBoxIconError)
			return
		}
		v.eventBus.Publish(eventMainViewShowMessage, v.name, retrieveOK, walk.MsgBoxOK)
	}

}
