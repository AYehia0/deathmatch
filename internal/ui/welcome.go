package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type WelcomeScreen struct {
	*AnimatedScreen
}

func NewWelcomeScreen(width, height int) *WelcomeScreen {
	colors := []lipgloss.Color{"196", "202", "208", "214", "220", "226"}
	return &WelcomeScreen{
		AnimatedScreen: NewAnimatedScreen(
			width,
			height,
			"ROBOT DEATHMATCH ARENA",
			"[h] How to Play  [c] Controls  [s] Scoring",
			"Press any key to start...",
			colors,
		),
	}
}
