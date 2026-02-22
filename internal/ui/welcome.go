package ui

import (
	"github.com/ayehia0/deathmatch/internal/game"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeScreen struct {
	*AnimatedScreen
	topScores []game.ScoreEntry
}

func NewWelcomeScreen(width, height int) *WelcomeScreen {
	topScores := game.GetTopScores(3)

	subtitle := ""
	if len(topScores) > 0 {
		subtitle = "TOP SCORES: "
		for i, s := range topScores {
			if i > 0 {
				subtitle += " | "
			}
			subtitle += formatInt(i+1) + ". " + s.Name + " Lvl" + formatInt(s.Level) + " " + formatInt(s.Score) + "pts"
		}
	}

	colors := []lipgloss.Color{"196", "202", "208", "214", "220", "226"}
	return &WelcomeScreen{
		AnimatedScreen: NewAnimatedScreen(
			width,
			height,
			"ROBOT DEATHMATCH ARENA",
			"",
			subtitle,
			"[h] How to Play  [c] Controls  [s] Scoring",
			colors,
		),
		topScores: topScores,
	}
}
