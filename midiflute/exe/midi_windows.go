//go:build windows

package main

// winmm.dll による MIDI 入力。鍵盤が押された瞬間にコールバックで受け取る(定期的に見に行かない)。

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winmm             = windows.NewLazySystemDLL("winmm.dll")
	procInGetNumDevs  = winmm.NewProc("midiInGetNumDevs")
	procInGetDevCapsW = winmm.NewProc("midiInGetDevCapsW")
	procInOpen        = winmm.NewProc("midiInOpen")
	procInStart       = winmm.NewProc("midiInStart")
	procInStop        = winmm.NewProc("midiInStop")
	procInReset       = winmm.NewProc("midiInReset")
	procInClose       = winmm.NewProc("midiInClose")
)

const (
	callbackFunction = 0x30000
	mimData          = 0x3C3
)

type midiInCaps struct {
	Mid           uint16
	Pid           uint16
	DriverVersion uint32
	Pname         [32]uint16
	Support       uint32
}

// 受信したショートメッセージ(ステータス+データ2バイト)
var midiCh = make(chan uint32, 4096)

var midiCallback = syscall.NewCallback(func(h, msg, inst, p1, p2 uintptr) uintptr {
	if msg == mimData {
		select {
		case midiCh <- uint32(p1):
		default: // 溢れたら捨てる(通常起きない)
		}
	}
	return 0
})

func ListMidiDevices() []string {
	n, _, _ := procInGetNumDevs.Call()
	var names []string
	for i := uintptr(0); i < n; i++ {
		var caps midiInCaps
		r, _, _ := procInGetDevCapsW.Call(i, uintptr(unsafe.Pointer(&caps)), unsafe.Sizeof(caps))
		if r != 0 {
			names = append(names, "")
			continue
		}
		names = append(names, windows.UTF16ToString(caps.Pname[:]))
	}
	return names
}

type MidiIn struct {
	h    uintptr
	Name string
}

func OpenMidi(name string) (*MidiIn, error) {
	for i, n := range ListMidiDevices() {
		if n != name {
			continue
		}
		var h uintptr
		r, _, _ := procInOpen.Call(uintptr(unsafe.Pointer(&h)), uintptr(i), midiCallback, 0, callbackFunction)
		if r != 0 {
			return nil, errors.New(T("midi_busy"))
		}
		procInStart.Call(h)
		return &MidiIn{h: h, Name: name}, nil
	}
	return nil, errors.New(T("midi_notfound"))
}

func (m *MidiIn) Close() {
	if m == nil || m.h == 0 {
		return
	}
	procInStop.Call(m.h)
	procInReset.Call(m.h)
	procInClose.Call(m.h)
	m.h = 0
}
