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
