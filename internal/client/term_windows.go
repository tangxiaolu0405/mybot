package client

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func rawModeOS() (func() error, error) {
	fd := os.Stdin.Fd()
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	var old uint32
	if r, _, _ := getConsoleMode.Call(fd, uintptr(unsafe.Pointer(&old))); r == 0 {
		return func() error { return nil }, fmt.Errorf("GetConsoleMode failed")
	}
	// disable line input and echo
	const (
		enableLineInput  = 0x0002
		enableEchoInput  = 0x0004
		enableProcessedInput = 0x0001
	)
	raw := old &^ (enableLineInput | enableEchoInput | enableProcessedInput)
	if r, _, _ := setConsoleMode.Call(fd, uintptr(raw)); r == 0 {
		return func() error { return nil }, fmt.Errorf("SetConsoleMode failed")
	}
	return func() error {
		_, _, _ = setConsoleMode.Call(fd, uintptr(old))
		return nil
	}, nil
}

func readKeyOS() (rune, error) {
	var buf [1]byte
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	readConsole := kernel32.NewProc("ReadConsoleA")
	getStdHandle := kernel32.NewProc("GetStdHandle")
	const stdInputHandle = uint32(0xfffffff6) // -10
	handle, _, _ := getStdHandle.Call(uintptr(stdInputHandle))

	var n uint32
	r, _, _ := readConsole.Call(handle, uintptr(unsafe.Pointer(&buf[0])), 1, uintptr(unsafe.Pointer(&n)), 0)
	if r == 0 || n == 0 {
		return 0, fmt.Errorf("read failed")
	}
	b := buf[0]

	switch b {
	case '\r', '\n':
		return keyEnter, nil
	case ' ':
		return keySpace, nil
	case 'j':
		return keyJ, nil
	case 'k':
		return keyK, nil
	case 'q':
		return keyQ, nil
	case '\t':
		return keyTab, nil
	case 0x1b: // escape
		// read arrow key sequence on Windows
		var seq [2]byte
		r, _, _ := readConsole.Call(handle, uintptr(unsafe.Pointer(&seq[0])), 2, uintptr(unsafe.Pointer(&n)), 0)
		if r == 0 || n < 2 {
			return keyEsc, nil
		}
		if seq[0] == '[' {
			switch seq[1] {
			case 'A':
				return keyUp, nil
			case 'B':
				return keyDown, nil
			}
		}
		return keyEsc, nil
	}
	return rune(b), nil
}
