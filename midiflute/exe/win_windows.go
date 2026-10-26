//go:build windows

package main

import (
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	kernel32                     = windows.NewLazySystemDLL("kernel32.dll")
	procSetProcessInformation    = kernel32.NewProc("SetProcessInformation")
)

var gameExes = map[string]bool{"client_tos_x64.exe": true, "client_tos.exe": true}

// foregroundIsToS は、いま前面にあるウィンドウがToSかどうか
func foregroundIsToS() bool {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return false
	}
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h)
	buf := make([]uint16, windows.MAX_PATH)
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return false
	}
	return gameExes[strings.ToLower(filepath.Base(windows.UTF16ToString(buf[:size])))]
}

// disablePowerThrottling は、裏に回ったときにWindowsの省電力制御で処理が間引かれるのを防ぐ
func disablePowerThrottling() {
	type throttlingState struct {
		Version     uint32
		ControlMask uint32
		StateMask   uint32
	}
	st := throttlingState{Version: 1, ControlMask: 1, StateMask: 0} // EXECUTION_SPEED を無効化
	const processPowerThrottling = 4
	procSetProcessInformation.Call(uintptr(windows.CurrentProcess()), processPowerThrottling,
		uintptr(unsafe.Pointer(&st)), unsafe.Sizeof(st))
}
