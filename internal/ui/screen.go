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

type AnimatedScreen struct {
	width     int
	height    int
	particles []Particle
	frame     int
	title     string
	subtitle  string
	prompt    string
	titleColor []lipgloss.Color
}

func NewAnimatedScreen(width, height int, title, subtitle, prompt string, titleColor []lipgloss.Color) *AnimatedScreen {
	screen := &AnimatedScreen{
		width:      width,
		height:     height,
		particles:  make([]Particle, 30),
		title:      title,
		subtitle:   subtitle,
		prompt:     prompt,
		titleColor: titleColor,
	}

	for i := range screen.particles {
		screen.particles[i] = Particle{
			X:    rand.Intn(width),
			Y:    rand.Intn(height),
			Char: string([]rune{'·', '•', '*', '○', '●'}[rand.Intn(5)]),
		}
	}

	return screen
}

func (s *AnimatedScreen) Update() {
	s.frame++
	for i := range s.particles {
		s.particles[i].X += rand.Intn(3) - 1
		s.particles[i].Y += rand.Intn(3) - 1

		if s.particles[i].X < 0 {
			s.particles[i].X = s.width - 1
		}
		if s.particles[i].X >= s.width {
			s.particles[i].X = 0
		}
		if s.particles[i].Y < 0 {
			s.particles[i].Y = s.height - 1
		}
		if s.particles[i].Y >= s.height {
			s.particles[i].Y = 0
		}
	}
}

func (s *AnimatedScreen) Render() string {
	grid := make([][]rune, s.height)
	styles := make([][]lipgloss.Style, s.height)
	for i := range grid {
		grid[i] = make([]rune, s.width)
		styles[i] = make([]lipgloss.Style, s.width)
		for j := range grid[i] {
			grid[i][j] = ' '
			styles[i][j] = lipgloss.NewStyle()
		}
	}

	for _, p := range s.particles {
		if p.Y >= 0 && p.Y < s.height && p.X >= 0 && p.X < s.width {
			grid[p.Y][p.X] = []rune(p.Char)[0]
			styles[p.Y][p.X] = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		}
	}

	titleY := s.height/2 - 3
	titleX := (s.width - len(s.title)) / 2

	if titleY >= 0 && titleY < s.height && titleX >= 0 {
		colorIdx := (s.frame / 10) % len(s.titleColor)
		style := lipgloss.NewStyle().Foreground(s.titleColor[colorIdx]).Bold(true)

		for j, ch := range s.title {
			x := titleX + j
			if x >= 0 && x < s.width {
				grid[titleY][x] = ch
				styles[titleY][x] = style
			}
		}
	}

	if s.subtitle != "" {
		subtitleY := titleY + 2
		subtitleX := (s.width - len(s.subtitle)) / 2
		if subtitleY >= 0 && subtitleY < s.height && subtitleX >= 0 {
			subtitleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			for j, ch := range s.subtitle {
				x := subtitleX + j
				if x >= 0 && x < s.width {
					grid[subtitleY][x] = ch
					styles[subtitleY][x] = subtitleStyle
				}
			}
		}
	}

	if s.prompt != "" {
		promptY := titleY + 4
		promptX := (s.width - len(s.prompt)) / 2
		if promptY >= 0 && promptY < s.height && promptX >= 0 {
			blinkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
			if s.frame%20 < 10 {
				blinkStyle = blinkStyle.Foreground(lipgloss.Color("240"))
			}
			for j, ch := range s.prompt {
				x := promptX + j
				if x >= 0 && x < s.width {
					grid[promptY][x] = ch
					styles[promptY][x] = blinkStyle
				}
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
