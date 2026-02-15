package game

import (
	"math/rand/v2"
)

type Position struct {
	X, Y int
}

type EntityType int

const (
	EntityRobot EntityType = iota
	EntityObstacle
	EntityJunk
	EntityShrub
)

type Entity struct {
	Pos  Position
	Type EntityType
}

type Difficulty struct {
	RobotCount    int
	ObstacleCount int
	MinSpawnDist  int
}

type Game struct {
	Width            int
	Height           int
	Player           Position
	Entities         []Entity
	GameOver         bool
	Teleports        int
	EMPs             int
	EMPTurnsLeft     int
	Score            int
	Level            int
	ConsecutiveKills int
	Blasters         int
	BlasterActive    bool
	BlasterTarget    Position
}

func New(width, height int, difficulty Difficulty) *Game {
	playerPos := Position{X: width / 2, Y: height / 2}
	occupied := []Position{playerPos}

	entities := []Entity{}

	robots := generatePositions(width, height, difficulty.RobotCount, difficulty.MinSpawnDist, occupied)
	for _, pos := range robots {
		entities = append(entities, Entity{Pos: pos, Type: EntityRobot})
	}
	occupied = append(occupied, robots...)

	obstacles := generatePositions(width, height, difficulty.ObstacleCount, 0, occupied)
	for _, pos := range obstacles {
		entities = append(entities, Entity{Pos: pos, Type: EntityObstacle})
	}

	return &Game{
		Width:     width,
		Height:    height,
		Player:    playerPos,
		Entities:  entities,
		Teleports: 5,
		EMPs:      3,
		Blasters:  2,
		Level:     1,
	}
}

func (g *Game) NextLevel() {
	g.Level++
	playerPos := Position{X: g.Width / 2, Y: g.Height / 2}
	occupied := []Position{playerPos}

	entities := []Entity{}

	robotCount := 10 + (g.Level-1)*2
	obstacleCount := 15 + (g.Level-1)*3
	minSpawnDist := max(3, 5-(g.Level-1)/2)

	robots := generatePositions(g.Width, g.Height, robotCount, minSpawnDist, occupied)
	for _, pos := range robots {
		entities = append(entities, Entity{Pos: pos, Type: EntityRobot})
	}
	occupied = append(occupied, robots...)

	obstacles := generatePositions(g.Width, g.Height, obstacleCount, 0, occupied)
	for _, pos := range obstacles {
		entities = append(entities, Entity{Pos: pos, Type: EntityObstacle})
	}

	g.Player = playerPos
	g.Entities = entities
	g.EMPTurnsLeft = 0
	g.BlasterActive = false
	g.Teleports += 5
	g.EMPs += 3
	g.Blasters += 1
	g.Score += 50
	g.ConsecutiveKills = 0
}

func generatePositions(width, height, count, minDist int, occupied []Position) []Position {
	positions := make([]Position, 0, count)

	for len(positions) < count {
		pos := Position{X: rand.IntN(width), Y: rand.IntN(height)}

		if isOccupied(pos, occupied) || isOccupied(pos, positions) {
			continue
		}

		if minDist > 0 && !isFarEnough(pos, occupied, minDist) {
			continue
		}

		positions = append(positions, pos)
	}

	return positions
}

func isOccupied(pos Position, positions []Position) bool {
	for _, p := range positions {
		if p.X == pos.X && p.Y == pos.Y {
			return true
		}
	}
	return false
}

func isFarEnough(pos Position, positions []Position, minDist int) bool {
	minDistSq := minDist * minDist
	for _, p := range positions {
		dx := pos.X - p.X
		dy := pos.Y - p.Y
		if dx*dx+dy*dy < minDistSq {
			return false
		}
	}
	return true
}

func (g *Game) MovePlayer(dx, dy int) {
	if g.GameOver {
		return
	}

	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	if newX < 0 || newX >= g.Width || newY < 0 || newY >= g.Height {
		return
	}

	newPos := Position{X: newX, Y: newY}

	for _, entity := range g.Entities {
		if entity.Pos.X == newPos.X && entity.Pos.Y == newPos.Y {
			g.GameOver = true
			return
		}
	}

	g.Player = newPos
	g.MoveRobots()
	g.CheckCollisions()
}

func (g *Game) MoveRobots() {
	if g.EMPTurnsLeft > 0 {
		g.EMPTurnsLeft--
		return
	}

	for i := range g.Entities {
		if g.Entities[i].Type != EntityRobot {
			continue
		}

		dx := 0
		dy := 0

		if g.Entities[i].Pos.X < g.Player.X {
			dx = 1
		} else if g.Entities[i].Pos.X > g.Player.X {
			dx = -1
		}

		if g.Entities[i].Pos.Y < g.Player.Y {
			dy = 1
		} else if g.Entities[i].Pos.Y > g.Player.Y {
			dy = -1
		}

		newPos := Position{
			X: g.Entities[i].Pos.X + dx,
			Y: g.Entities[i].Pos.Y + dy,
		}

		if newPos.X < 0 || newPos.X >= g.Width || newPos.Y < 0 || newPos.Y >= g.Height {
			continue
		}

		hitObstacle := false
		hitJunk := false
		for _, entity := range g.Entities {
			if entity.Pos == newPos {
				if entity.Type == EntityObstacle {
					hitObstacle = true
					break
				}
				if entity.Type == EntityJunk {
					hitJunk = true
					break
				}
			}
		}

		if hitObstacle {
			continue
		}

		if hitJunk {
			g.Entities[i].Type = EntityJunk
		} else {
			g.Entities[i].Pos = newPos
		}
	}
}

func (g *Game) CheckCollisions() {
	posMap := make(map[Position][]int)

	for i, entity := range g.Entities {
		if entity.Type == EntityRobot {
			posMap[entity.Pos] = append(posMap[entity.Pos], i)

			if entity.Pos == g.Player {
				g.GameOver = true
				return
			}
		}
	}

	toRemove := make(map[int]bool)
	var junkPositions []Position

	for pos, indices := range posMap {
		if len(indices) > 1 {
			for _, idx := range indices {
				toRemove[idx] = true
			}
			junkPositions = append(junkPositions, pos)

			killCount := len(indices)
			g.ConsecutiveKills += killCount
			multiplier := 1 + g.ConsecutiveKills/5
			g.Score += 10 * killCount * multiplier
		}
	}

	newEntities := []Entity{}
	for i, entity := range g.Entities {
		if !toRemove[i] {
			newEntities = append(newEntities, entity)
		}
	}

	for _, pos := range junkPositions {
		newEntities = append(newEntities, Entity{Pos: pos, Type: EntityJunk})
	}

	g.Entities = newEntities

	robotCount := 0
	for _, entity := range g.Entities {
		if entity.Type == EntityRobot {
			robotCount++
		}
	}

	if robotCount == 0 {
		g.NextLevel()
	}
}

func (g *Game) Teleport() bool {
	if g.Teleports <= 0 || g.GameOver {
		return false
	}

	occupiedMap := make(map[Position]bool)
	for _, entity := range g.Entities {
		occupiedMap[entity.Pos] = true
	}

	maxAttempts := 100
	for range maxAttempts {
		newPos := Position{
			X: rand.IntN(g.Width),
			Y: rand.IntN(g.Height),
		}

		if !occupiedMap[newPos] {
			g.Player = newPos
			g.Teleports--
			g.Score -= 2
			return true
		}
	}

	return false
}

func (g *Game) UseEMP() bool {
	if g.EMPs <= 0 || g.GameOver {
		return false
	}

	g.EMPs--
	g.EMPTurnsLeft = 5
	return true
}

func (g *Game) ToggleBlaster() bool {
	if g.GameOver {
		return false
	}

	if !g.BlasterActive {
		if g.Blasters <= 0 {
			return false
		}
		g.BlasterActive = true
		g.BlasterTarget = g.Player
		return true
	}

	g.BlasterActive = false
	g.Blasters--

	killCount := 0
	newEntities := []Entity{}

	for _, entity := range g.Entities {
		inBlastZone := entity.Pos.X >= g.BlasterTarget.X-1 && entity.Pos.X <= g.BlasterTarget.X+1 &&
			entity.Pos.Y >= g.BlasterTarget.Y-1 && entity.Pos.Y <= g.BlasterTarget.Y+1

		if entity.Type == EntityRobot && inBlastZone {
			killCount++
		} else {
			newEntities = append(newEntities, entity)
		}
	}

	g.Entities = newEntities

	if killCount > 0 {
		g.ConsecutiveKills += killCount
		multiplier := 1 + g.ConsecutiveKills/5
		g.Score += 10 * killCount * multiplier
	}

	playerInBlastZone := g.Player.X >= g.BlasterTarget.X-1 && g.Player.X <= g.BlasterTarget.X+1 &&
		g.Player.Y >= g.BlasterTarget.Y-1 && g.Player.Y <= g.BlasterTarget.Y+1

	if playerInBlastZone {
		g.GameOver = true
	}

	robotCount := 0
	for _, entity := range g.Entities {
		if entity.Type == EntityRobot {
			robotCount++
		}
	}

	if robotCount == 0 && !g.GameOver {
		g.NextLevel()
	}

	return true
}

func (g *Game) MoveBlasterTarget(dx, dy int) {
	if !g.BlasterActive {
		return
	}

	newX := g.BlasterTarget.X + dx
	newY := g.BlasterTarget.Y + dy

	if newX >= 1 && newX < g.Width-1 && newY >= 1 && newY < g.Height-1 {
		g.BlasterTarget.X = newX
		g.BlasterTarget.Y = newY
	}
}
