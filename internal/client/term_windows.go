package client

import (
	"fmt"
	"os"
	"syscall"
	"unicode/utf16"
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
	const (
		enableLineInput           = 0x0002
		enableEchoInput           = 0x0004
		enableProcessedInput      = 0x0001
		enableVirtualTerminalInput = 0x0200
	)
	raw := old &^ (enableLineInput | enableEchoInput | enableProcessedInput)
	raw |= enableVirtualTerminalInput
	if r, _, _ := setConsoleMode.Call(fd, uintptr(raw)); r == 0 {
		return func() error { return nil }, fmt.Errorf("SetConsoleMode failed")
	}
	return func() error {
		_, _, _ = setConsoleMode.Call(fd, uintptr(old))
		return nil
	}, nil
}

func readKeyOS() (rune, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	readConsoleW := kernel32.NewProc("ReadConsoleW")
	getStdHandle := kernel32.NewProc("GetStdHandle")
	const stdInputHandle = uint32(0xfffffff6) // -10
	handle, _, _ := getStdHandle.Call(uintptr(stdInputHandle))

	readWide := func() (uint16, error) {
		var ch uint16
		var n uint32
		r, _, _ := readConsoleW.Call(handle, uintptr(unsafe.Pointer(&ch)), 1, uintptr(unsafe.Pointer(&n)), 0)
		if r == 0 || n == 0 {
			return 0, fmt.Errorf("read failed")
		}
		return ch, nil
	}

	ch, err := readWide()
	if err != nil {
		return 0, err
	}

	switch ch {
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
	case 0x1b: // escape — read arrow key sequence
		next, err := readWide()
		if err != nil {
			return keyEsc, nil
		}
		if next == '[' {
			next, err = readWide()
			if err != nil {
				return keyEsc, nil
			}
			switch next {
			case 'A':
				return keyUp, nil
			case 'B':
				return keyDown, nil
			}
		}
		return keyEsc, nil
	}

	// Surrogate pair handling for characters outside BMP.
	if utf16.IsSurrogate(rune(ch)) {
		low, err := readWide()
		if err != nil {
			return rune(ch), nil
		}
		return utf16.DecodeRune(rune(ch), rune(low)), nil
	}

	return rune(ch), nil
}
