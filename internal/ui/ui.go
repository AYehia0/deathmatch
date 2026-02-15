package ui

import (
	"strings"
	"time"

	"github.com/ahmedyahia/deathmatch/internal/game"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type state int

const (
	welcomeState state = iota
	helpState
	gameState
	gameOverState
)

type helpTab int

const (
	howToPlayTab helpTab = iota
	controlsTab
	scoringTab
)

type Model struct {
	game           *game.Game
	state          state
	width          int
	height         int
	welcomeScreen  *WelcomeScreen
	gameOverScreen *AnimatedScreen
	viewport       viewport.Model
	activeTab      helpTab
	finalScore     int
	finalLevel     int
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
			m.game = game.New((msg.Width-4)/2, msg.Height-5, game.Difficulty{
				RobotCount:    10,
				ObstacleCount: 15,
				MinSpawnDist:  5,
			})
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
		if m.state == gameOverState && m.gameOverScreen != nil {
			m.gameOverScreen.Update()
		}
		if m.state == gameState && m.game != nil && m.game.GameOver {
			m.state = gameOverState
			m.finalScore = m.game.Score
			m.finalLevel = m.game.Level
			colors := []lipgloss.Color{"9", "196", "160", "124"}
			m.gameOverScreen = NewAnimatedScreen(
				m.width,
				m.height,
				"GAME OVER",
				"Level: "+formatInt(m.finalLevel)+"  Score: "+formatInt(m.finalScore),
				"[r] Restart  [q] Quit",
				colors,
			)
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
			case "esc":
				if m.game.BlasterActive {
					m.game.BlasterActive = false
				}
			case "t":
				if !m.game.BlasterActive {
					m.game.Teleport()
				}
			case "e":
				if !m.game.BlasterActive {
					m.game.UseEMP()
				}
			case "b":
				m.game.ToggleBlaster()
			case "up", "k":
				if m.game.BlasterActive {
					m.game.MoveBlasterTarget(0, -1)
				} else {
					m.game.MovePlayer(0, -1)
				}
			case "down", "j":
				if m.game.BlasterActive {
					m.game.MoveBlasterTarget(0, 1)
				} else {
					m.game.MovePlayer(0, 1)
				}
			case "left", "h":
				if m.game.BlasterActive {
					m.game.MoveBlasterTarget(-1, 0)
				} else {
					m.game.MovePlayer(-1, 0)
				}
			case "right", "l":
				if m.game.BlasterActive {
					m.game.MoveBlasterTarget(1, 0)
				} else {
					m.game.MovePlayer(1, 0)
				}
			}
		}

		if m.state == gameOverState {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "r":
				m.game = game.New((m.width-4)/2, m.height-5, game.Difficulty{
					RobotCount:    10,
					ObstacleCount: 15,
					MinSpawnDist:  5,
				})
				m.state = gameState
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
	if m.state == gameOverState {
		if m.gameOverScreen != nil {
			return m.gameOverScreen.Render()
		}
		return ""
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

func gameView(g *game.Game) string {
	grid := make([][]string, g.Height)
	blasterGrid := make([][]bool, g.Height)
	for i := range grid {
		grid[i] = make([]string, g.Width)
		blasterGrid[i] = make([]bool, g.Width)
		for j := range grid[i] {
			grid[i][j] = "  "
		}
	}

	if g.BlasterActive {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				y := g.BlasterTarget.Y + dy
				x := g.BlasterTarget.X + dx
				if y >= 0 && y < g.Height && x >= 0 && x < g.Width {
					blasterGrid[y][x] = true
				}
			}
		}
	}

	for _, entity := range g.Entities {
		if entity.Pos.Y >= 0 && entity.Pos.Y < g.Height && entity.Pos.X >= 0 && entity.Pos.X < g.Width {
			grid[entity.Pos.Y][entity.Pos.X] = renderEntity(entity)
		}
	}

	if g.Player.Y >= 0 && g.Player.Y < g.Height && g.Player.X >= 0 && g.Player.X < g.Width {
		grid[g.Player.Y][g.Player.X] = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("@@")
	}

	var arena strings.Builder
	for y, row := range grid {
		for x, cell := range row {
			if blasterGrid[y][x] {
				if cell == "  " {
					arena.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("░░"))
				} else {
					arena.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("240")).Render(cell))
				}
			} else {
				arena.WriteString(cell)
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

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	empStatus := ""
	if g.EMPTurnsLeft > 0 {
		empStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Render(" (ACTIVE: " + formatInt(g.EMPTurnsLeft) + " turns)")
	}

	blasterStatus := ""
	if g.BlasterActive {
		blasterStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Render(" [TARGETING MODE - Press 'b' to fire, 'esc' to cancel]")
	}

	status := statusStyle.Render(
		"Level: " + formatInt(g.Level) +
			"  Score: " + formatInt(g.Score) +
			"  [t] Teleports: " + formatInt(g.Teleports) +
			"  [e] EMPs: " + formatInt(g.EMPs) + empStatus +
			"  [b] Blasters: " + formatInt(g.Blasters) + blasterStatus +
			"  [q] Quit",
	)

	return boxStyle.Render(arena.String()) + "\n" + status
}

func formatInt(n int) string {
	if n < 0 {
		return "0"
	}
	if n < 10 {
		return string(rune('0' + n))
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+(n%10))) + result
		n /= 10
	}
	return result
}

func renderEntity(e game.Entity) string {
	switch e.Type {
	case game.EntityRobot:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("RR")
	case game.EntityObstacle:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("##")
	case game.EntityJunk:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("**")
	case game.EntityShrub:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("&&")
	default:
		return "  "
	}
}
