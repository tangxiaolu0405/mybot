//go:build linux

package client

import (
	"fmt"
	"os"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

func rawModeOS() (func() error, error) {
	fd := int(os.Stdin.Fd())
	var old syscall.Termios
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&old))); err != 0 {
		return func() error { return nil }, fmt.Errorf("tcgetattr: %v", err)
	}
	raw := old
	raw.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&raw))); err != 0 {
		return func() error { return nil }, fmt.Errorf("tcsetattr: %v", err)
	}
	return func() error {
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&old)))
		return nil
	}, nil
}

func readKeyOS() (rune, error) {
	var buf [4]byte
	n, err := os.Stdin.Read(buf[:1])
	if err != nil || n == 0 {
		return 0, err
	}
	if buf[0] != 0x1b { // not escape
		switch buf[0] {
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
		}
		// Multi-byte UTF-8 sequence: read remaining bytes and decode.
		if buf[0] >= 0x80 {
			extra := utf8BytesRemaining(buf[0])
			for i := 0; i < extra; i++ {
				if _, err := os.Stdin.Read(buf[1+i : 2+i]); err != nil {
					return rune(buf[0]), nil
				}
			}
			r, _ := utf8.DecodeRune(buf[:1+extra])
			if r != utf8.RuneError || extra == 0 {
				return r, nil
			}
		}
		return rune(buf[0]), nil
	}
	// escape sequence
	if _, err := os.Stdin.Read(buf[1:2]); err != nil {
		return keyEsc, nil
	}
	if buf[1] != '[' {
		return keyEsc, nil
	}
	if _, err := os.Stdin.Read(buf[2:3]); err != nil {
		return keyEsc, nil
	}
	switch buf[2] {
	case 'A':
		return keyUp, nil
	case 'B':
		return keyDown, nil
	}
	return keyEsc, nil
}
