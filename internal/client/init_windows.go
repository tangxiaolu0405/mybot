//go:build windows

package client

import (
	"os"
	"syscall"
	"unsafe"
)

func init() {
	enableVirtualTerminal(os.Stdout)
	enableVirtualTerminal(os.Stderr)
}

func enableVirtualTerminal(f *os.File) {
	fd := f.Fd()
	var mode uint32
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	if r, _, _ := getConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode))); r != 0 {
		const enableVirtualTerminalProcessing = 0x0004
		setConsoleMode.Call(fd, uintptr(mode|enableVirtualTerminalProcessing))
	}
}
