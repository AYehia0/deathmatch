package game

type Position struct {
	X, Y int
}

type Game struct {
	Width  int
	Height int
	Player Position
}

func New(width, height int) *Game {
	return &Game{
		Width:  width,
		Height: height,
		Player: Position{X: width / 2, Y: height / 2},
	}
}

func (g *Game) MovePlayer(dx, dy int) {
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	if newX >= 0 && newX < g.Width && newY >= 0 && newY < g.Height {
		g.Player.X = newX
		g.Player.Y = newY
	}
}
