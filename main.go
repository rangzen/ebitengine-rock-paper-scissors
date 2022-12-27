package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"log"
	"math/rand"
)

// GUI
const (
	screenWidth  = 300
	screenHeight = 600
	objWidth     = 48 / 2
	objHeight    = 48 / 2
	nbType       = 3
	nbObjPerType = 50
	nbObj        = nbObjPerType * nbType
)

// States
const (
	rock = iota
	paper
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

// init loads all the resources
func init() {
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

type Obj struct {
	r    int // resource: rock, paper, scissors
	x, y int
}

type Game struct {
	o [nbObj]Obj
}

func (g *Game) Update() error {
	for i := 0; i < nbObj; i++ {
		g.o[i].x += rand.Intn(3) - 1
		g.o[i].y += rand.Intn(3) - 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Rock Paper Scissors")
	for _, o := range g.o {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.x-objWidth), float64(o.y-objHeight))
		switch o.r {
		case rock:
			screen.DrawImage(rockImg, op)
		case paper:
			screen.DrawImage(paperImg, op)
		case scissors:
			screen.DrawImage(scissorsImg, op)
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth * 2, screenHeight * 2
}

func (g *Game) Init() {
	log.Println("Initiating game")
	for i := 0; i < nbObj; i++ {
		g.o[i] = Obj{
			r: i % 3,
			x: rand.Intn(screenWidth * 2),
			y: rand.Intn(screenHeight * 2),
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rock Paper Scissors")
	game := &Game{}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
