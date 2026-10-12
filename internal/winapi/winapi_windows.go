// Package winapi wraps the few Win32 calls Vigilant needs.
package winapi

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	inputKeyboard     = 1
	keyeventfKeyUp    = 0x0002
	vkF15             = 0x7E
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
)

var (
	user32                = windows.NewLazySystemDLL("user32.dll")
	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procSendInput         = user32.NewProc("SendInput")
	procSetExecutionState = kernel32.NewProc("SetThreadExecutionState")
)

// input mirrors the Win32 INPUT structure for keyboard events. The trailing
// padding makes it as large as the biggest union member (MOUSEINPUT), which
// SendInput requires: 40 bytes on 64-bit Windows.
type input struct {
	typ   uint32
	_     uint32
	vk    uint16
	scan  uint16
	flags uint32
	time  uint32
	_     uint32
	extra uintptr
	_     [8]byte
}

// PressF15 simulates a press and release of the F15 key, which Windows counts
// as user activity but which no application normally reacts to.
func PressF15() error {
	inputs := [2]input{
		{typ: inputKeyboard, vk: vkF15},
		{typ: inputKeyboard, vk: vkF15, flags: keyeventfKeyUp},
	}
	n, _, err := procSendInput.Call(uintptr(len(inputs)), uintptr(unsafe.Pointer(&inputs[0])), unsafe.Sizeof(inputs[0]))
	if int(n) != len(inputs) {
		return fmt.Errorf("SendInput injected %d of %d events: %w", n, len(inputs), err)
	}
	return nil
}

// ResetIdleTimers tells Windows the system (and optionally the display) is in
// use, resetting the idle timers once. It does not leave any lasting state.
func ResetIdleTimers(display bool) error {
	flags := uintptr(esSystemRequired)
	if display {
		flags |= esDisplayRequired
	}
	if r, _, err := procSetExecutionState.Call(flags); r == 0 {
		return fmt.Errorf("SetThreadExecutionState: %w", err)
	}
	return nil
}

// Open opens a URL or folder with the default handler.
func Open(target string) error {
	p, err := windows.UTF16PtrFromString(target)
	if err != nil {
		return err
	}
	verb, _ := windows.UTF16PtrFromString("open")
	return windows.ShellExecute(0, verb, p, nil, nil, windows.SW_SHOWNORMAL)
}

// AcquireSingleInstance reports whether this is the only running Vigilant.
// The underlying mutex lives until the process exits.
func AcquireSingleInstance() (bool, error) {
	name, err := windows.UTF16PtrFromString(`Local\Vigilant-SingleInstance`)
	if err != nil {
		return false, err
	}
	_, err = windows.CreateMutex(nil, false, name)
	if err == windows.ERROR_ALREADY_EXISTS {
		return false, nil
	}
	return err == nil, err
}
