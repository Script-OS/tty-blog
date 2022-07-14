package view

import (
	"fmt"
	"github.com/muesli/termenv"
	"golang.org/x/term"
	"os"
	"strings"
	"tty-blog/cmd/view/renderer/input"
	"tty-blog/cmd/view/renderer/webmedia"
)

type model struct {
	ready           bool
	lines           []string
	y               int
	width           int
	height          int
	MouseWheelDelta int
}

func (m *model) Update(msg input.Msg) bool {

	switch msg := msg.(type) {
	case input.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return false
		}
		if m.ready {
			switch msg.String() {
			case "pgdown":
				m.ViewDown()

			case "pgup":
				m.ViewUp()

			case "down":
				m.LineDown(1)

			case "up":
				m.LineUp(1)
			}
		}

	case input.MouseMsg:
		if m.ready {
			switch msg.Type {
			case input.MouseWheelUp:
				m.LineUp(m.MouseWheelDelta)

			case input.MouseWheelDown:
				m.LineDown(m.MouseWheelDelta)
			}
		}
	}

	return true
}

func (m *model) AtTop() bool {
	return m.y <= 0
}

func (m *model) AtBottom() bool {
	return m.y >= m.maxY()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *model) maxY() int {
	return max(0, len(m.lines)-m.height)
}

func (m *model) SetY(n int) {
	m.y = min(m.maxY(), max(0, n))
}

func (m *model) ViewDown() {
	if m.AtBottom() {
		return
	}
	m.SetY(m.y + m.height)
}

func (m *model) ViewUp() {
	if m.AtTop() {
		return
	}
	m.SetY(m.y - m.height)
}

func (m *model) LineDown(n int) {
	if m.AtBottom() || n == 0 {
		return
	}
	m.SetY(m.y + n)
}

func (m *model) LineUp(n int) {
	if m.AtTop() || n == 0 {
		return
	}
	m.SetY(m.y - n)
}

func (m model) visibleLines() (lines []string) {
	if len(m.lines) > 0 {
		top := max(0, m.y)
		bottom := min(len(m.lines), max(top, m.y+m.height))
		lines = m.lines[top:bottom]
	}
	return lines
}

func (m model) View() {
	if !m.ready {
		return
	}

	fmt.Print(webmedia.ResetWebmedia())
	lines := m.visibleLines()
	for i, line := range lines {
		termenv.MoveCursor(i, 0)
		fmt.Print(line)
	}
	for i := len(lines); i < m.height; i += 1 {
		termenv.MoveCursor(i, 0)
		fmt.Print(strings.Repeat(" ", m.width))
	}

	return
}

func RenderInPage(content string) {
	m := model{
		lines:           strings.Split(content, "\n"),
		y:               0,
		MouseWheelDelta: 3,
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	termenv.AltScreen()
	termenv.HideCursor()

	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	m.width = w
	m.height = h
	m.ready = true

	m.View()

	quit := false
	for !quit {
		msgs, err := input.ReadInputs(os.Stdin)
		if err != nil {
			break
		}
		if len(msgs) == 0 {
			continue
		}
		for _, msg := range msgs {
			if !m.Update(msg) {
				quit = true
			}
		}

		m.View()
	}

	termenv.ShowCursor()
	termenv.ExitAltScreen()
}
