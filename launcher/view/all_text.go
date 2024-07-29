package view

// LauncherView
const (
	// Main Window
	launcherViewTitle  string = "Launcher"
	launcherViewWidth  int    = 400
	launcherViewHeight int    = 300
	// Auth Widget
	authGroupBoxTitlte      string = "Auth"
	textEditMaxWidth        int    = 200
	textEditMaxHeight       int    = 20
	loginTabPageTitle       string = "Login"
	signupTabPageTitle      string = "Signup"
	retrieveTabPageTitle    string = "Retrieve"
	usernameLabelText       string = "Username"
	usernameMinLength       int    = 5
	usernameMaxLength       int    = 20
	passwordLabelText       string = "Password"
	passwordAgainLableText  string = "PasswordAgain"
	passwordMinLength       int    = 5
	passwordMaxLength       int    = 20
	secondPasswordLableText string = "SecondPassword"
	secondPasswordMinLength int    = 6
	secondPasswordMaxLength int    = 10
	submitButtonText        string = "Submit"
	clearButtonText         string = "Clear"
	// Home Widget
	accountTabPageTitle      string = "Account"
	passwordTabPageTitle     string = "Password"
	logoutTabPageTitle       string = "Logout"
	isForeverBannedLabelText string = "IsForeverBanned"
	cashPointLabelText       string = "CashPoint"
	maplePointLabelText      string = "MaplePoint"
	isSecondPassword         string = "IsSecondPassword"
	newPasswordLabelText     string = "NewPassword"

	// Menu Widget
	menuGroupBoxTitlte    string = "Menu"
	startGameButtonText   string = "Start Game"
	skipSDOAuthButtonText string = "Skip SDO Auth"
	kickGameButtonText    string = "Kick Game"
	setGameDirButtonText  string = "Set Game Dir"
	openGameDirButtonText string = "Open Game Dir"
	patchGameButtonText   string = "Patch Game"
	patchHostsButtonText  string = "Patch Hosts"
	exitGameButtonText    string = "Exit Game"
)

// MsgBox Info
const (
	checkExit        string = "Do you want to exit? "
	signupOK         string = "Signup OK"
	retrieveOK       string = "Retrieve OK"
	updatePasswordOK string = "Update password OK"
	startingGame     string = "Starting game..."
	setGameDirOK     string = "Set game dir ok"
	patchGameOK      string = "Patch game ok"
	patchHostsOK     string = "Patch hosts ok"
	unpatchHostsOK   string = "Unpatch hosts ok"
	exitingGame      string = "Exiting game..."
	handleOK         string = "Handle OK"
)

// MsgBox Error
const (
	errEditTextLength    string = "Less than minimum length"
	errDifferentPassword string = "Password is different with PasswordAgain"
	errDuplicateHosts    string = "Hosts has been patched, Do you want to unpatch the hosts?"
)
