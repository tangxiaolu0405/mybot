package client

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type cmdDef struct {
	Name    string
	Aliases []string
	Desc    string
}

var commands = []cmdDef{
	{Name: "config", Desc: "edit configuration"},
	{Name: "exit", Aliases: []string{"quit", "q"}, Desc: "exit cata"},
	{Name: "clear", Aliases: []string{"reset"}, Desc: "reset chat session"},
	{Name: "cls", Desc: "clear terminal screen"},
	{Name: "help", Desc: "show available commands"},
}

// readLine reads a line in raw mode. When the input starts with "/",
// matching commands are shown as a selectable list below the input.
// Use ↑/↓ to navigate, Enter to select the highlighted command.
func readLine() (string, error) {
	restore, err := rawMode()
	if err != nil {
		return "", err
	}
	defer restore()

	var buf []byte
	sugLines := 0
	sel := 0

	renderInput(buf, &sel, &sugLines)

	for {
		key, err := readKey()
		if err != nil {
			meta("\n")
			return "", err
		}

		inCommand := len(buf) > 0 && buf[0] == '/'
		done := false
		var result string

		switch key {
		case keyEnter:
			if inCommand && sugLines > 0 {
				matches := matchCmds(string(buf[1:]))
				if sel >= 0 && sel < len(matches) {
					buf = []byte("/" + matches[sel].Name + " ")
					sel = 0
				}
			} else {
				result = string(buf)
				done = true
			}

		case 0x03: // Ctrl+C
			done = true

		case keyEsc:
			if inCommand && sugLines > 0 {
				// Exit command mode: clear the input
				buf = nil
				sel = 0
			} else if len(buf) == 0 {
				done = true
			}

		case 0x04: // Ctrl+D
			if len(buf) == 0 {
				meta("\n")
				return "", fmt.Errorf("EOF")
			}

		case 0x7f, 0x08: // Backspace
			if len(buf) > 0 {
				buf = trimLastUTF8(buf)
				sel = 0
			}

		case keyTab:
			if inCommand {
				buf = tabComplete(buf)
				sel = 0
			}

		case keyUp:
			if inCommand && sel > 0 {
				sel--
			}

		case keyDown:
			if inCommand && sugLines > 1 {
				matches := matchCmds(string(buf[1:]))
				if sel < len(matches)-1 {
					sel++
				}
			}

		case keyK:
			buf = appendRune(buf, 'k')
			sel = 0
		case keyJ:
			buf = appendRune(buf, 'j')
			sel = 0
		case keyQ:
			buf = appendRune(buf, 'q')
			sel = 0

		default:
			if key >= 0x20 {
				buf = appendRune(buf, key)
				sel = 0
			}
		}

		if done {
			clearArea()
			meta("%s› %s%s\n", ansiBold, ansiReset, result)
			if strings.TrimSpace(result) == `"""` {
				return readMultiRaw()
			}
			return result, nil
		}

		rerenderInput(buf, &sel, &sugLines)
	}
}

func tabComplete(buf []byte) []byte {
	prefix := string(buf[1:])
	matches := matchCmds(prefix)
	if len(matches) == 0 {
		return buf
	}
	if len(matches) == 1 {
		return []byte("/" + matches[0].Name + " ")
	}
	common := matches[0].Name
	for _, m := range matches[1:] {
		common = commonPrefix(common, m.Name)
	}
	if len(common) > len(prefix) {
		return []byte("/" + common)
	}
	return buf
}

// --- rendering ---
//
// Layout (cursor starts at anchor line A, before the initial \n):
//
//	\n  › /con                 ← input line  (A+1)
//	\n    ──────────────────   ← separator   (A+2)
//	\n    ▶ /config  desc      ← suggestion  (A+3)
//	\n       /clear   desc     ← suggestion  (A+4)
//
// After renderInput/sugMove, cursor is at the end of the input line.

func renderInput(buf []byte, sel *int, sugLines *int) {
	meta("\n%s› %s%s%s", ansiBold, ansiReset, string(buf), clearLine())
	*sugLines = 0
	if len(buf) > 0 && buf[0] == '/' {
		matches := matchCmds(string(buf[1:]))
		if len(matches) > 0 {
			meta("\n  %s%s%s%s", ansiDim, sepLine(), ansiReset, clearLine())
			*sugLines = 1
			for i, c := range matches {
				m := "     "
				if i == *sel {
					m = fmt.Sprintf("  %s▶%s ", ansiCyan, ansiReset)
				}
				meta("\n%s%s/%s%s  %s%s%s%s", m, ansiBold, c.Name, ansiReset, ansiDim, c.Desc, ansiReset, clearLine())
				*sugLines++
			}
		}
	}
	// Position cursor at end of input line
	if *sugLines > 0 {
		upN(*sugLines)
	}
	meta("\r%s› %s%s", ansiBold, ansiReset, string(buf))
}

func rerenderInput(buf []byte, sel *int, sugLines *int) {
	oldSug := *sugLines
	upN(1) // from input line to anchor
	renderInput(buf, sel, sugLines)
	if *sugLines < oldSug {
		// Move past new suggestion area to first leftover line
		for i := 0; i <= *sugLines; i++ {
			meta("\r\n")
		}
		// Erase from here to end of display — clears all leftovers in one shot
		meta("\033[J")
		// Back to input line
		upN(*sugLines + 1)
		meta("\r%s› %s%s", ansiBold, ansiReset, string(buf))
	}
}

func clearArea() {
	upN(1) // from input line to anchor
	meta("\033[J") // erase from anchor to end of display
	// Cursor stays at anchor — caller's next print lands on the cleared anchor line
}

func sepLine() string { return strings.Repeat("─", 30) }

// --- matching ---

func matchCmds(prefix string) []cmdDef {
	prefix = strings.ToLower(prefix)
	var result []cmdDef
	for _, c := range commands {
		if strings.HasPrefix(strings.ToLower(c.Name), prefix) {
			result = append(result, c)
			continue
		}
		for _, a := range c.Aliases {
			if strings.HasPrefix(strings.ToLower(a), prefix) {
				result = append(result, c)
				break
			}
		}
	}
	return result
}

func commonPrefix(a, b string) string {
	n := 0
	for n < len(a) && n < len(b) && a[n] == b[n] {
		n++
	}
	return a[:n]
}

func appendRune(buf []byte, r rune) []byte {
	var enc [utf8.UTFMax]byte
	n := utf8.EncodeRune(enc[:], r)
	return append(buf, enc[:n]...)
}

func trimLastUTF8(buf []byte) []byte {
	if len(buf) == 0 {
		return buf
	}
	i := len(buf) - 1
	for i > 0 && buf[i]&0xC0 == 0x80 {
		i--
	}
	return buf[:i]
}

// --- multi-line input (""") ---

func readMultiRaw() (string, error) {
	var b strings.Builder
	for {
		line, err := readLineRaw()
		if err != nil {
			return b.String(), nil
		}
		if strings.TrimSpace(line) == `"""` {
			return b.String(), nil
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}
}

func readLineRaw() (string, error) {
	var buf []byte
	meta("%s… %s", ansiDim, ansiReset)

	for {
		key, err := readKey()
		if err != nil {
			meta("\n")
			return "", err
		}

		switch key {
		case keyEnter:
			meta("\n")
			return string(buf), nil
		case 0x03, keyEsc:
			meta("\n")
			return "", nil
		case 0x7f, 0x08:
			if len(buf) > 0 {
				buf = trimLastUTF8(buf)
				meta("\b \b")
			}
		case keyK:
			buf = appendRune(buf, 'k')
			meta("k")
		case keyJ:
			buf = appendRune(buf, 'j')
			meta("j")
		case keyQ:
			buf = appendRune(buf, 'q')
			meta("q")
		default:
			if key >= 0x20 {
				buf = appendRune(buf, key)
				meta(string(key))
			}
		}
	}
}
