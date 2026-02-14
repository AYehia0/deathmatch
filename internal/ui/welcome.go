package ui

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Particle struct {
	X, Y int
	Char string
}

type WelcomeScreen struct {
	width     int
	height    int
	particles []Particle
	frame     int
}

func NewWelcomeScreen(width, height int) *WelcomeScreen {
	ws := &WelcomeScreen{
		width:     width,
		height:    height,
		particles: make([]Particle, 30),
	}

	for i := range ws.particles {
		ws.particles[i] = Particle{
			X:    rand.Intn(width),
			Y:    rand.Intn(height),
			Char: string([]rune{'·', '•', '*', '○', '●'}[rand.Intn(5)]),
		}
	}

	return ws
}

func (ws *WelcomeScreen) Update() {
	ws.frame++
	for i := range ws.particles {
		ws.particles[i].X += rand.Intn(3) - 1
		ws.particles[i].Y += rand.Intn(3) - 1

		if ws.particles[i].X < 0 {
			ws.particles[i].X = ws.width - 1
		}
		if ws.particles[i].X >= ws.width {
			ws.particles[i].X = 0
		}
		if ws.particles[i].Y < 0 {
			ws.particles[i].Y = ws.height - 1
		}
		if ws.particles[i].Y >= ws.height {
			ws.particles[i].Y = 0
		}
	}
}

func (ws *WelcomeScreen) Render() string {
	colors := []lipgloss.Color{"196", "202", "208", "214", "220", "226"}

	grid := make([][]rune, ws.height)
	styles := make([][]lipgloss.Style, ws.height)
	for i := range grid {
		grid[i] = make([]rune, ws.width)
		styles[i] = make([]lipgloss.Style, ws.width)
		for j := range grid[i] {
			grid[i][j] = ' '
			styles[i][j] = lipgloss.NewStyle()
		}
	}

	for _, p := range ws.particles {
		if p.Y >= 0 && p.Y < ws.height && p.X >= 0 && p.X < ws.width {
			grid[p.Y][p.X] = []rune(p.Char)[0]
			styles[p.Y][p.X] = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		}
	}

	title := "ROBOT DEATHMATCH ARENA"
	titleY := ws.height/2 - 3
	titleX := (ws.width - len(title)) / 2

	if titleY >= 0 && titleY < ws.height && titleX >= 0 {
		colorIdx := (ws.frame / 10) % len(colors)
		style := lipgloss.NewStyle().Foreground(colors[colorIdx]).Bold(true)

		for j, ch := range title {
			x := titleX + j
			if x >= 0 && x < ws.width {
				grid[titleY][x] = ch
				styles[titleY][x] = style
			}
		}
	}

	help := "[h] How to Play  [c] Controls  [s] Scoring"
	helpY := titleY + 2
	helpX := (ws.width - len(help)) / 2
	if helpY >= 0 && helpY < ws.height && helpX >= 0 {
		helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		for j, ch := range help {
			x := helpX + j
			if x >= 0 && x < ws.width {
				grid[helpY][x] = ch
				styles[helpY][x] = helpStyle
			}
		}
	}

	prompt := "Press any key to start..."
	promptY := helpY + 2
	promptX := (ws.width - len(prompt)) / 2
	if promptY >= 0 && promptY < ws.height && promptX >= 0 {
		blinkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
		if ws.frame%20 < 10 {
			blinkStyle = blinkStyle.Foreground(lipgloss.Color("240"))
		}
		for j, ch := range prompt {
			x := promptX + j
			if x >= 0 && x < ws.width {
				grid[promptY][x] = ch
				styles[promptY][x] = blinkStyle
			}
		}
	}

	var output strings.Builder
	for i := range grid {
		for j := range grid[i] {
			output.WriteString(styles[i][j].Render(string(grid[i][j])))
		}
		if i < len(grid)-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}

func init() {
}
