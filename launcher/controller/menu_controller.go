package controller

import (
	"bufio"
	"bytes"
	"launcher/api"
	"launcher/assets"
	"launcher/config"
	"launcher/util"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	gameFilename    = "MapleStory.exe"
	patchFilename   = "ijl15.dll"
	hostsDir        = "C:/Windows/System32/drivers/etc"
	hostsFilename   = "hosts"
	gpkAssetsDomain = "127.0.0.1 mxd.dorado.sdo.com"
)

var startGameArgs []string = []string{"221.231.130.70", "8484"}

type menuController struct {
	conf *config.Config
}

func NewMenuController(conf *config.Config) MenuController {
	c := &menuController{
		conf: conf,
	}
	return c
}

// StartGame implements MenuController.
func (c *menuController) StartGame() error {
	temp := path.Join(c.conf.GameDir, gameFilename)
	return util.StartProcess(temp, startGameArgs)
}

// SkipSDOAuth implements MenuController.
func (c *menuController) SkipSDOAuth() (string, error) {
	bytes, err := api.SendRequest(http.MethodGet, api.PathAuthGameSkipSDOAuth, nil, c.conf.AccessToken)
	if err != nil {
		return "", err
	}
	msg := string(bytes)
	return msg, nil
}

// KickGame implements MenuController.
func (c *menuController) KickGame() (string, error) {
	bytes, err := api.SendRequest(http.MethodGet, api.PathAuthGameKick, nil, c.conf.AccessToken)
	if err != nil {
		return "", err
	}
	msg := string(bytes)
	return msg, nil
}

// SetGameDir implements MenuController.
func (c *menuController) SetGameDir(dir string) error {
	c.conf.GameDir = dir
	return config.Save(c.conf)
}

// GetCurrentGameDir implements MenuController.
func (c *menuController) GetCurrentGameDir() string {
	return c.conf.GameDir
}

// PatchGame implements MenuController.
func (c *menuController) PatchGame() error {
	temp := path.Join(c.conf.GameDir, patchFilename)
	_, err := os.Stat(temp)
	if err != nil {
		return err
	}
	file, err := os.Create(temp)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(assets.PatchData)
	if err != nil {
		return err
	}
	return nil
}

// CheckHosts implements MenuController.
func (c *menuController) CheckHosts() bool {
	path := path.Join(hostsDir, hostsFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	hostsContent := string(data)
	return strings.Contains(hostsContent, gpkAssetsDomain)
}

// PatchHosts implements MenuController.
func (c *menuController) PatchHosts() error {
	path := path.Join(hostsDir, hostsFilename)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(gpkAssetsDomain)
	if err != nil {
		return err
	}
	return nil
}

// UnpatchHosts implements MenuController.
func (c *menuController) UnpatchHosts() error {
	path := path.Join(hostsDir, hostsFilename)
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var bs []byte
	buf := bytes.NewBuffer(bs)
	for scanner.Scan() {
		line := scanner.Text()
		length := len(line)
		if length >= 3 && strings.Contains(line, gpkAssetsDomain) {
			_, err = buf.WriteString("\n")
			if err != nil {
				return err
			}
			break
		} else {
			buf.WriteString(line + "\n")
		}
	}
	file.Truncate(0)
	file.Seek(0, 0)
	_, err = buf.WriteTo(file)
	return err
}

// ExitGmae implements MenuController.
func (c *menuController) ExitGmae() error {
	return util.KillProcess(gameFilename)
}
