package ui

import (
	"strings"
	"time"

	"github.com/ahmedyahia/deathmatch/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type state int

const (
	welcomeState state = iota
	helpState
	gameState
)

type helpTab int

const (
	howToPlayTab helpTab = iota
	controlsTab
	scoringTab
)

type Model struct {
	game          *game.Game
	state         state
	width         int
	height        int
	welcomeScreen *WelcomeScreen
	viewport      viewport.Model
	activeTab     helpTab
}

func NewModel() Model {
	return Model{
		state: welcomeState,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.game == nil {
			m.game = game.New(msg.Width-4, msg.Height-2)
		}
		if m.welcomeScreen == nil {
			m.welcomeScreen = NewWelcomeScreen(msg.Width, msg.Height)
		}
		m.viewport = viewport.New(msg.Width, msg.Height-4)
		m.viewport.SetContent(m.getHelpContent())
		return m, nil
	case tickMsg:
		if m.state == welcomeState && m.welcomeScreen != nil {
			m.welcomeScreen.Update()
		}
		return m, tick()
	case tea.KeyMsg:
		if m.state == welcomeState {
			switch msg.String() {
			case "h":
				m.state = helpState
				m.activeTab = howToPlayTab
				m.viewport.SetContent(m.getHelpContent())
				return m, nil
			case "c":
				m.state = helpState
				m.activeTab = controlsTab
				m.viewport.SetContent(m.getHelpContent())
				return m, nil
			case "s":
				m.state = helpState
				m.activeTab = scoringTab
				m.viewport.SetContent(m.getHelpContent())
				return m, nil
			default:
				m.state = gameState
				return m, nil
			}
		}

		if m.state == helpState {
			switch msg.String() {
			case "q":
				m.state = welcomeState
				return m, nil
			case "tab":
				m.activeTab = (m.activeTab + 1) % 3
				m.viewport.SetContent(m.getHelpContent())
				m.viewport.GotoTop()
				return m, nil
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}

		if m.state == gameState {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				m.game.MovePlayer(0, -1)
			case "down", "j":
				m.game.MovePlayer(0, 1)
			case "left", "h":
				m.game.MovePlayer(-1, 0)
			case "right", "l":
				m.game.MovePlayer(1, 0)
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.state == welcomeState {
		if m.welcomeScreen != nil {
			return m.welcomeScreen.Render()
		}
		return ""
	}
	if m.state == helpState {
		return m.renderHelp()
	}
	return gameView(m.game)
}

func (m Model) getHelpContent() string {
	var content string
	switch m.activeTab {
	case howToPlayTab:
		content = `# HOW TO PLAY

You are trapped in a sealed arena, hunted by robots.
Your only weapon is **movement**—lure robots into each other, shape the arena 
with debris and shrubs, and survive as long as possible.

## Arena Survival
- Trapped in a walled arena with hostile robots
- You are unarmed; survival depends on movement and positioning

## Robot Collisions & Junk
- Robots chase you and collide with each other
- Collisions create **radioactive junk**
- Robots hitting junk self-destruct
- Junk is **deadly to humans**

## Shrubs
- Shrubs block movement but not vision
- Robots may explode or crush shrubs
- Useful as temporary shields
- Running into shrubs causes damage and score loss

## Defensive Tools
- **Teleporter**: limited uses, random safe relocation
- **EMP**: single use, disables robots for 3 turns
- Tools reset each level

## Endless Progression
- Clear robots to advance
- Each level adds more enemies
- High score is the only goal

> There is no escape, no final level, and no winning—only a score to beat.`

	case controlsTab:
		content = `# CONTROLS

## Movement
Use **arrow keys** or **hjkl** (vim keys) to move:
- **h** - Move left
- **j** - Move down
- **k** - Move up
- **l** - Move right

## Actions
- **t** - Use teleporter (limited uses)
- **e** - Use EMP (single use per level)
- **q** - Quit game

## Navigation
- **Tab** - Switch between help tabs
- **Arrow keys / j/k** - Scroll help text
- **q** - Return to welcome screen`

	case scoringTab:
		content = `# SCORING

## Points
- **Destroy robot**: +10 points
- **Survive level**: +50 points
- **Consecutive kills**: bonus multiplier

## Penalties
- **Hit shrub**: -5 points
- **Use teleporter**: -2 points

## Difficulty Multiplier
- **Easy**: 1x points
- **Normal**: 1.5x points
- **Hard**: 2x points

## High Score
- Your best score is saved
- Try to beat your personal record
- Survive as long as possible`
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.width-4),
	)
	rendered, _ := r.Render(content)
	return rendered
}

func (m Model) renderHelp() string {
	tabs := []string{"How to Play [h]", "Controls [c]", "Scoring [s]"}
	var tabsRendered []string
	
	for i, tab := range tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if helpTab(i) == m.activeTab {
			style = style.Foreground(lipgloss.Color("226")).Bold(true)
		} else {
			style = style.Foreground(lipgloss.Color("240"))
		}
		tabsRendered = append(tabsRendered, style.Render(tab))
	}
	
	header := lipgloss.JoinHorizontal(lipgloss.Top, tabsRendered...)
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("\nTab: switch | ↑↓/jk: scroll | q: back")
	
	return header + "\n\n" + m.viewport.View() + footer
}

func welcomeView(width, height int) string {
	content := "DEATHMATCH\n\nPress any key to continue..."

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(content)
}

func gameView(g *game.Game) string {
	var arena strings.Builder

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if x == g.Player.X && y == g.Player.Y {
				arena.WriteString("@")
			} else {
				arena.WriteString(" ")
			}
		}
		if y < g.Height-1 {
			arena.WriteString("\n")
		}
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	return boxStyle.Render(arena.String())
}
