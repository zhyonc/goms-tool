package util

import (
	"os/exec"
	"syscall"
)

func GetScreenResolution() (int, int) {
	lazyProc := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`)
	screenWidth, _, _ := lazyProc.Call(uintptr(0))
	screenHeight, _, _ := lazyProc.Call(uintptr(1))
	return int(screenWidth), int(screenHeight)
}

func IsLessThanMinLength(minLength int, texts ...string) bool {
	for _, text := range texts {
		if len(text) < minLength {
			return true
		}
	}
	return false
}

func StartProcess(path string, args []string) error {
	cmd := exec.Command(path, args...)
	return cmd.Run()
}

func KillProcess(processName string) error {
	cmd := exec.Command("taskkill", "/IM", processName, "/F")
	return cmd.Run()
}
