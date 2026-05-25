package client

import (
	"fmt"
	"os"
)

func upN(n int) {
	for i := 0; i < n; i++ {
		fmt.Fprint(os.Stderr, "\033[1A")
	}
}

// SelectOption represents a single choice in an interactive selector.
type SelectOption struct {
	ID    string
	Label string
	Desc  string // optional, shown dim on the same line after the label
}

// Select presents an interactive list and returns the chosen option ID.
// Use ↑/↓ or j/k to move, Enter/Space to confirm, 1-9 shortcuts, Esc/q to cancel.
func Select(prompt, detail string, options []SelectOption) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options")
	}

	restore, err := rawMode()
	if err != nil {
		return "", err
	}
	defer restore()

	hideCursor()
	defer showCursor()

	idx := 0
	renderSingle(prompt, detail, options, idx)
	sc := 1 + len(options) // prompt + options
	if detail != "" {
		sc++
	}
	defer clearSelect(sc)

	for {
		key, err := readKey()
		if err != nil {
			return "", err
		}
		switch key {
		case keyUp, keyK:
			if idx > 0 {
				idx--
			}
		case keyDown, keyJ:
			if idx < len(options)-1 {
				idx++
			}
		case keyEnter, keySpace:
			rerenderSingle(prompt, detail, options, idx)
			return options[idx].ID, nil
		case keyEsc:
			rerenderSingle(prompt, detail, options, idx)
			return "", nil
		default:
			if k := int(key - '0'); k >= 1 && k <= len(options) {
				rerenderSingle(prompt, detail, options, k-1)
				return options[k-1].ID, nil
			}
		}
		rerenderSingle(prompt, detail, options, idx)
	}
}

// SelectMulti presents a multi-select list. Tab toggles selection, Enter confirms all.
// Returns IDs of all selected options. Esc/q returns nil.
func SelectMulti(prompt, detail string, options []SelectOption) ([]string, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("no options")
	}

	restore, err := rawMode()
	if err != nil {
		return nil, err
	}
	defer restore()

	hideCursor()
	defer showCursor()

	idx := 0
	selected := make(map[int]bool)
	renderMulti(prompt, detail, options, idx, selected)
	mc := 2 + len(options) // prompt + hint + options
	if detail != "" {
		mc++
	}
	defer clearSelect(mc)

	for {
		key, err := readKey()
		if err != nil {
			return nil, err
		}
		switch key {
		case keyUp, keyK:
			if idx > 0 {
				idx--
			}
		case keyDown, keyJ:
			if idx < len(options)-1 {
				idx++
			}
		case keyTab:
			selected[idx] = !selected[idx]
		case keyEnter:
			rerenderMulti(prompt, detail, options, idx, selected)
			return selectedIDs(options, selected), nil
		case keyEsc, keyQ:
			rerenderMulti(prompt, detail, options, idx, selected)
			return nil, nil
		default:
			if k := int(key - '0'); k >= 1 && k <= len(options) {
				selected[k-1] = !selected[k-1]
			}
		}
		rerenderMulti(prompt, detail, options, idx, selected)
	}
}

func selectedIDs(options []SelectOption, selected map[int]bool) []string {
	var ids []string
	for i := range options {
		if selected[i] {
			ids = append(ids, options[i].ID)
		}
	}
	return ids
}

// --- single-select rendering ---

func renderSingle(prompt, detail string, options []SelectOption, selected int) {
	meta("\n  %s%s%s%s\n", ansiBold, prompt, ansiReset, clearLine())
	if detail != "" {
		meta("  %s%s%s%s\n", ansiDim, detail, ansiReset, clearLine())
	}
	for i, opt := range options {
		if i == selected {
			meta("  %s▶%s %s%s%s", ansiCyan, ansiReset, ansiBold, opt.Label, ansiReset)
		} else {
			meta("     %s%s%s", ansiDim, opt.Label, ansiReset)
		}
		if opt.Desc != "" {
			meta("  %s- %s%s", ansiDim, opt.Desc, ansiReset)
		}
		meta("%s\n", clearLine())
	}
}

func rerenderSingle(prompt, detail string, options []SelectOption, selected int) {
	lines := 2 + len(options) // \n (initial pos) + prompt \n + each option \n
	if detail != "" {
		lines++
	}
	upN(lines)
	renderSingle(prompt, detail, options, selected)
}

// --- multi-select rendering ---

func renderMulti(prompt, detail string, options []SelectOption, idx int, selected map[int]bool) {
	meta("\n  %s%s%s%s\n", ansiBold, prompt, ansiReset, clearLine())
	if detail != "" {
		meta("  %s%s%s%s\n", ansiDim, detail, ansiReset, clearLine())
	}
	meta("  %stab%s toggle  %senter%s confirm  %sesc%s cancel%s\n",
		ansiDim, ansiReset, ansiDim, ansiReset, ansiDim, ansiReset, clearLine())
	for i, opt := range options {
		check := " "
		if selected[i] {
			check = fmt.Sprintf("%s✓%s", ansiGreen, ansiReset)
		}
		if i == idx {
			meta("  %s▶%s [%s] %s%s%s", ansiCyan, ansiReset, check, ansiBold, opt.Label, ansiReset)
		} else {
			meta("     [%s] %s%s%s", check, ansiDim, opt.Label, ansiReset)
		}
		if opt.Desc != "" {
			meta("  %s- %s%s", ansiDim, opt.Desc, ansiReset)
		}
		meta("%s\n", clearLine())
	}
}

func rerenderMulti(prompt, detail string, options []SelectOption, idx int, selected map[int]bool) {
	lines := 3 + len(options) // initial \n + prompt \n + hint \n + each option \n
	if detail != "" {
		lines++
	}
	upN(lines)
	renderMulti(prompt, detail, options, idx, selected)
}

// --- helpers ---

func clearSelect(count int) {
	upN(count)
	for i := 0; i < count; i++ {
		meta("\r%s\n", clearLine())
	}
	upN(count)
}

func hideCursor()        { meta("\033[?25l") }
func showCursor()        { meta("\033[?25h") }
func clearLine() string { return "\033[0K" }
