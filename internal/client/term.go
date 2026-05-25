package client

// Key codes returned by readKey.
const (
	keyUp    = 0x100 + iota
	keyDown
	keyEnter
	keyEsc
	keySpace
	keyK
	keyJ
	keyQ
	keyTab
)

// rawMode enables raw terminal input. Returns a restore function.
func rawMode() (func(), error) {
	r, err := rawModeOS()
	if err != nil {
		return func() {}, err
	}
	return func() { _ = r() }, nil
}

// readKey reads a single keypress from stdin in raw mode.
func readKey() (rune, error) {
	return readKeyOS()
}

// utf8BytesRemaining returns the number of continuation bytes expected
// after the given UTF-8 lead byte.
func utf8BytesRemaining(first byte) int {
	if first&0xE0 == 0xC0 {
		return 1
	}
	if first&0xF0 == 0xE0 {
		return 2
	}
	if first&0xF8 == 0xF0 {
		return 3
	}
	return 0
}
