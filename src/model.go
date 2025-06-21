package src

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	brightness int
	min        int
	max        int
	step       int
	width      int
	err        error
}

func InitialModel() Model {
	return Model{
		brightness: 50,
		min:        0,
		max:        100,
		step:       5,
		width:      40,
	}
}

func (m Model) Init() tea.Cmd {
	return getCurrentBrightness
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "down", "j":
			if m.brightness > m.min {
				m.brightness -= 1
				if m.brightness < m.min {
					m.brightness = m.min
				}
				return m, setBrightnessCmd(m.brightness)
			}
		case "up", "k":
			if m.brightness < m.max {
				m.brightness += 1
				if m.brightness > m.max {
					m.brightness = m.max
				}
				return m, setBrightnessCmd(m.brightness)
			}
		case "left", "h":
			if m.brightness > m.min {
				m.brightness -= m.step
				if m.brightness < m.min {
					m.brightness = m.min
				}
				return m, setBrightnessCmd(m.brightness)
			}
		case "right", "l":
			if m.brightness < m.max {
				m.brightness += m.step
				if m.brightness > m.max {
					m.brightness = m.max
				}
				return m, setBrightnessCmd(m.brightness)
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Quick brightness presets
			val, _ := strconv.Atoi(msg.String())
			m.brightness = min(val*10, m.max)
			return m, setBrightnessCmd(m.brightness)
		}

	case brightnessMsg:
		m.brightness = msg.brightness
		m.err = nil

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Current brightness display
	brightnessDisplay := fmt.Sprintf("Current Brightness: %d%%", m.brightness)
	b.WriteString(brightnessDisplay + "\n\n")

	// Slider
	slider := m.renderSlider()
	b.WriteString(slider + "\n\n")

	// Controls
	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Render("j k : Adjust brightness by 1%\nh l : Adjust brightness by step (5%)\n0-9 : Quick presets (0=0%, 9=90%)\nq / esc / ctrl+c : Quit")

	b.WriteString(controls)

	// Error display
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
		b.WriteString("\n\n" + errorStyle.Render("Error: "+m.err.Error()))
	}

	return b.String()
}

func (m Model) renderSlider() string {
	// Calculate slider position
	pos := int(float64(m.brightness) / float64(m.max) * float64(m.width))

	var slider strings.Builder

	for i := range m.width {
		if i < pos {
			slider.WriteString("█")
		} else {
			slider.WriteString("░")
		}
	}

	return slider.String()
}
