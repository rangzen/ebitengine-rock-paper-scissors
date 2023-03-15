package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GUI
const (
	screenWidth   = 300
	screenHeight  = 600
	objScale      = .33
	objHalfWidth  = 48 / 2
	objHalfHeight = 48 / 2
	nbType        = 3
	nbObjPerType  = 50
	nbObj         = nbObjPerType * nbType
)

// Types
// The order reflect force in the game not the name
// paper > rock > scissors (> paper, etc.)
const (
	paper = iota
	rock
	scissors
)

//go:embed resources/icons8-rock-48.png
var rockPng []byte
var rockImg *ebiten.Image

//go:embed resources/icons8-page-facing-up-48.png
var paperPng []byte
var paperImg *ebiten.Image

//go:embed resources/icons8-scissors-48.png
var scissorsPng []byte
var scissorsImg *ebiten.Image

type NNS interface {
	Neighbor([]Obj, int) (int, error)
}

// init loads all the resources
func init() {
	rand.Seed(time.Now().Unix())
	rockDecoded, _, err := image.Decode(bytes.NewReader(rockPng))
	if err != nil {
		log.Fatal("decoding rock:", err)
	}
	rockImg = ebiten.NewImageFromImage(rockDecoded)

	paperDecoded, _, err := image.Decode(bytes.NewReader(paperPng))
	if err != nil {
		log.Fatal("decoding paper:", err)
	}
	paperImg = ebiten.NewImageFromImage(paperDecoded)

	scissorsDecoded, _, err := image.Decode(bytes.NewReader(scissorsPng))
	if err != nil {
		log.Fatal("decoding scissors:", err)
	}
	scissorsImg = ebiten.NewImageFromImage(scissorsDecoded)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rock Paper Scissors")
	game := &Game{
		nns: &Linear{},
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Obj struct {
	t    int // type: rock, paper, scissors
	x, y int
}

func (o Obj) String() string {
	var t string
	switch o.t {
	case paper:
		t = "P"
	case rock:
		t = "R"
	case scissors:
		t = "S"
	}
	return fmt.Sprintf("%s(%d,%d)", t, o.x, o.y)
}

type Game struct {
	paused bool
	nns    NNS
	o      [nbObj]Obj
}

func (g *Game) Update() error {
	// Pause
	spacePressed := inpututil.IsKeyJustReleased(ebiten.KeySpace)
	if spacePressed {
		g.paused = !g.paused
		log.Println("paused:", g.paused)
	}
	if g.paused {
		return nil
	}

	// Create next state by copy
	var next [nbObj]Obj
	for i := 0; i < nbObj; i++ {
		next[i] = Obj{
			t: g.o[i].t,
			x: g.o[i].x,
			y: g.o[i].y,
		}
	}

	// Update next state
	for i := 0; i < nbObj; i++ {
		// Pass if the Obj has already a new state
		if g.o[i].t != next[i].t {
			continue
		}

		n, err := g.nns.Neighbor(g.o[:], i)
		if err != nil {
			log.Println("searching neighbor:", err)
			g.paused = true
			log.Println("game paused")
			continue
		}

		// Move
		switch g.o[i].t {
		case paper:
			switch g.o[n].t {
			case rock:
				next[i].x, next[i].y = pursue(g.o[i], g.o[n])
			case scissors:
				next[i].x, next[i].y = evade(g.o[i], g.o[n])
			}
		case rock:
			switch g.o[n].t {
			case scissors:
				next[i].x, next[i].y = pursue(g.o[i], g.o[n])
			case paper:
				next[i].x, next[i].y = evade(g.o[i], g.o[n])
			}
		case scissors:
			switch g.o[n].t {
			case paper:
				next[i].x, next[i].y = pursue(g.o[i], g.o[n])
			case rock:
				next[i].x, next[i].y = evade(g.o[i], g.o[n])
			}
		}

		// Catch
		if sqrDist(next[i], g.o[n]) <= 2 {
			next[n].t = g.o[i].t
		}
	}

	// Disjoint groups
	for i := 0; i < nbObj; i++ {
		for j := i + 1; j < nbObj; j++ {
			if next[i].t == next[j].t && next[i].x == next[j].x && next[i].y == next[j].y {
				next[i].x += rand.Intn(3) - 1
				next[i].x = clamp(next[i].x, 0, screenWidth)
				next[i].y += rand.Intn(3) - 1
				next[i].y = clamp(next[i].y, 0, screenHeight)
			}
		}
	}

	g.o = next
	return nil
}

func pursue(me Obj, prey Obj) (int, int) {
	var newX = me.x
	var newY = me.y
	if me.x < prey.x {
		newX++
	}
	if me.x > prey.x {
		newX--
	}
	if me.y < prey.y {
		newY++
	}
	if me.y > prey.y {
		newY--
	}
	newX = clamp(newX, 0, screenWidth)
	newY = clamp(newY, 0, screenHeight)
	return newX, newY
}

func evade(me Obj, predator Obj) (int, int) {
	var newX = me.x
	var newY = me.y
	if me.x < predator.x {
		newX--
	}
	if me.x > predator.x {
		newX++
	}
	if me.y < predator.y {
		newY--
	}
	if me.y > predator.y {
		newY++
	}
	newX = clamp(newX, 0, screenWidth)
	newY = clamp(newY, 0, screenHeight)

	// If cornered, move randomly
	if newX == me.x && newY == me.y {
		if newX == predator.x {
			newX += rand.Intn(3) - 1
			newX = clamp(newX, 0, screenWidth)
		} else {
			newY += rand.Intn(3) - 1
			newY = clamp(newY, 0, screenHeight)
		}
	}
	return newX, newY
}

func clamp(val int, min int, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Rock Paper Scissors")
	for _, o := range g.o {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(objScale, objScale)
		op.GeoM.Translate(float64(o.x)-objHalfWidth*objScale, float64(o.y)-objHalfHeight*objScale)
		switch o.t {
		case rock:
			screen.DrawImage(rockImg, op)
		case paper:
			screen.DrawImage(paperImg, op)
		case scissors:
			screen.DrawImage(scissorsImg, op)
		}
	}

	// Status
	s := [3]int{0, 0, 0}
	for _, o := range g.o {
		s[o.t]++
	}

	// Draw the message.
	tutorial := "Space: Pause"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s\n\nRock    : %3d\nPaper   : %3d\nScissors: %3d", ebiten.ActualTPS(), ebiten.ActualFPS(),
		tutorial,
		s[rock], s[paper], s[scissors])
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Init() {
	log.Println("Initiating game")
	for i := 0; i < nbObj; i++ {
		g.o[i] = Obj{
			t: i % 3,
			x: rand.Intn(screenWidth),
			y: rand.Intn(screenHeight),
		}
	}
}

// Beware,
// the following code search the nearest neighbor but from a different type
// than object given in parameter (index).

// Circular implements an index neighbor search.
// It is for testing purpose only.
type Circular struct{}

func (l *Circular) Neighbor(objs []Obj, i int) (int, error) {
	for j := 1; j < len(objs)-1; j++ {
		n := (i + j) % len(objs)
		if objs[i].t != objs[n].t {
			return n, nil
		}
	}
	return 0, errors.New("no neighbor found")
}

// Linear implements a linear search of the nearest neighbor.
// https://en.wikipedia.org/wiki/Nearest_neighbor_search#Linear_search
type Linear struct{}

func (l *Linear) Neighbor(objs []Obj, i int) (int, error) {
	bestIndex := -1
	bestDist := screenWidth * screenWidth * screenHeight * screenHeight
	for j := 1; j < len(objs)-1; j++ {
		n := (i + j) % len(objs)
		if objs[i].t != objs[n].t {
			dist := sqrDist(objs[i], objs[n])
			if dist < bestDist {
				bestIndex = n
				bestDist = dist
			}
		}
	}
	if bestIndex == -1 {
		return 0, errors.New("no neighbor found")
	}
	return bestIndex, nil
}

// sqrDist returns the square of the distance between two objects.
// To avoid a square root, we compare the square of the distance.
func sqrDist(o1 Obj, o2 Obj) int {
	return (o1.x-o2.x)*(o1.x-o2.x) + (o1.y-o2.y)*(o1.y-o2.y)
}
