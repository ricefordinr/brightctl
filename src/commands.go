package src

import tea "github.com/charmbracelet/bubbletea"

// Commands
func getCurrentBrightness() tea.Msg {
	brightness, err := getBrightness()
	if err != nil {
		return errMsg{err}
	}
	return brightnessMsg{brightness}
}

func setBrightnessCmd(brightness int) tea.Cmd {
	return func() tea.Msg {
		err := setBrightness(brightness)
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}
