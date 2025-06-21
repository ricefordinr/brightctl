package src

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getBrightness() (int, error) {
	// Try common backlight paths
	paths := []string{
		"/sys/class/backlight/intel_backlight/brightness",
		"/sys/class/backlight/acpi_video0/brightness",
		"/sys/class/backlight/nvidia_backlight/brightness",
	}

	maxPaths := []string{
		"/sys/class/backlight/intel_backlight/max_brightness",
		"/sys/class/backlight/acpi_video0/max_brightness",
		"/sys/class/backlight/nvidia_backlight/max_brightness",
	}

	for i, path := range paths {
		if _, err := os.Stat(path); err == nil {
			current, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			max, err := os.ReadFile(maxPaths[i])
			if err != nil {
				continue
			}

			currentVal, err := strconv.Atoi(strings.TrimSpace(string(current)))
			if err != nil {
				continue
			}

			maxVal, err := strconv.Atoi(strings.TrimSpace(string(max)))
			if err != nil {
				continue
			}

			return int(float64(currentVal) / float64(maxVal) * 100), nil
		}
	}

	// Fallback to xrandr if available
	cmd := exec.Command("xrandr", "--verbose")
	output, err := cmd.Output()
	if err != nil {
		return 50, fmt.Errorf("no backlight control found")
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Brightness:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				brightness, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err == nil {
					return int(brightness * 100), nil
				}
			}
		}
	}

	return 50, fmt.Errorf("could not determine brightness")
}

func setBrightness(brightness int) error {
	// Try sysfs first
	paths := []string{
		"/sys/class/backlight/intel_backlight/brightness",
		"/sys/class/backlight/acpi_video0/brightness",
		"/sys/class/backlight/nvidia_backlight/brightness",
	}

	maxPaths := []string{
		"/sys/class/backlight/intel_backlight/max_brightness",
		"/sys/class/backlight/acpi_video0/max_brightness",
		"/sys/class/backlight/nvidia_backlight/max_brightness",
	}

	for i, path := range paths {
		if _, err := os.Stat(path); err == nil {
			max, err := os.ReadFile(maxPaths[i])
			if err != nil {
				continue
			}

			maxVal, err := strconv.Atoi(strings.TrimSpace(string(max)))
			if err != nil {
				continue
			}

			newVal := int(float64(brightness) / 100.0 * float64(maxVal))
			err = os.WriteFile(path, []byte(strconv.Itoa(newVal)), 0644)
			if err == nil {
				return nil
			}
		}
	}

	// Fallback to xrandr
	brightnessFloat := float64(brightness) / 100.0
	cmd := exec.Command("sh", "-c", fmt.Sprintf("xrandr --output $(xrandr | grep ' connected' | cut -d' ' -f1 | head -1) --brightness %.2f", brightnessFloat))
	return cmd.Run()
}
