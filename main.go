package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	_ "image/png"
	"log"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type Target struct {
	x int
	y int
	r int
	isDead bool
}

type Game struct {
	tick float64
	targetZones []Target
	kills int
}

var woodPlank *ebiten.Image
var water *ebiten.Image
var duck *ebiten.Image

const (
	screenWidth  = 800
	screenHeight = 600
)

func init() {

	woodPlank = loadImage(woodPlankImg)
	water = loadImage(waterImg)
	duck = loadImage(duckImg)

}

func loadImage(byteImage []byte) *ebiten.Image{
	img, _, err := image.Decode(bytes.NewReader(byteImage))
	if err != nil {
		log.Fatal(err)
	}
	ebit_image, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return ebit_image
}

// Internal state of the game (loop 1)
func (g *Game) Update(screen *ebiten.Image) error {
	g.tick ++
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft){
		x, y := ebiten.CursorPosition()
		fmt.Printf("{%d ,%d} \n", x, y)
		for i, s := range g.targetZones {
			fmt.Printf("Target {%v} \n", s)
			r := s.r
			if ((s.x - r < x) && (x < s.x + r)) && ((s.y - r < y) && (y < s.y + r)) && !g.targetZones[i].isDead{
					g.targetZones[i].isDead = true
					println("DEAD MODDAFOKAKAAH")
					time.Sleep(500)
					g.kills = g.kills + 1
			}
		}
	}
	return nil
}

// Draws the screen (loop 2)
func (g *Game) Draw(screen *ebiten.Image) {
	output := "KILLS: " + strconv.FormatInt(int64(g.kills), 10)
	println(output)
	ebitenutil.DebugPrint(screen, output)

	x, y := screen.Size()

	w_water, _ := water.Size()
	Water(x, screen, g, - (w_water / 2.0) - w_water, y - 400, w_water, float64(3))
	if !g.targetZones[0].isDead {
		g.targetZones[0] = Duck(x, y, screen, duck, g.tick, false, 400)
	}

	Water(x, screen, g, - (w_water / 2.0) - 2*w_water + 50, y - 350, w_water, float64(1))

	if !g.targetZones[1].isDead {
		g.targetZones[1] = Duck(x, y, screen, duck, g.tick, true, 400)
	}

	Water(x, screen, g, - (w_water / 2.0) - 2*w_water + 100 , y - 300, w_water , float64(3))

	w_plank, _ := woodPlank.Size()
	Floor(x, y, w_plank, screen)

}

func Duck(x int, y int, screen *ebiten.Image, duck *ebiten.Image, tick float64, reverse bool, y_start int) Target {
	w_duck, _ := duck.Size()

	var tx_duck = 0
	var ty_duck = 0
	var offset = 0
	var x_position = 0

	if !reverse{
		tx_duck = - w_duck
		ty_duck = y - y_start

		offset = int(tick / 1.5) % (x + w_duck)

		op := &ebiten.DrawImageOptions{}
		x_position = tx_duck + offset
		op.GeoM.Translate(float64(x_position), float64(ty_duck))
		if err := screen.DrawImage(duck, op); err != nil {
			log.Fatal(err)
		}
	} else {
		tx_duck = x + w_duck
		ty_duck = y - y_start

		offset = int(tick / 1.5) % (x + w_duck)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(-1,1)
		x_position = tx_duck - offset
		op.GeoM.Translate(float64(x_position), float64(ty_duck))
		if err := screen.DrawImage(duck, op); err != nil {
			log.Fatal(err)
		}
	}


	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tx_duck + offset), float64(ty_duck))
	if err := screen.DrawImage(duck, op); err != nil {
		log.Fatal(err)
	}

	return Target{x_position, ty_duck, 150, false}
}

func Water(screen_width int, screen *ebiten.Image, g *Game, tx_water int, ty_water int, w_water int, speed float64){
	num_waters := int((float64(screen_width) / float64(w_water) ) + 3)

	var m = int(g.tick / speed) % w_water
	water_movement_x := 0 + m
	water_movement_y := 0

	for i := 0; i < num_waters; i++ {
		offset := float64(tx_water +(i*w_water) + water_movement_x)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offset, float64(ty_water + water_movement_y))
		if err := screen.DrawImage(water, op); err != nil {
			log.Fatal(err)
		}
	}
}

func Floor(x int, y int, w_plank int, screen *ebiten.Image){
	tx_plank := 0
	ty_plank := y - 200
	scalePlank := 200.0/ float64(w_plank)

	for i := 0; i < x/200; i++ {
		offset := float64(tx_plank +(i*200))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scalePlank, scalePlank)
		op.GeoM.Translate(offset, float64(ty_plank))
		if err := screen.DrawImage(woodPlank, op); err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Draw image")
	var start = []Target{
		{
			x:      0,
			y:      0,
			r:      0,
			isDead: false,
		},
		{
			x:      0,
			y:      0,
			r:      0,
			isDead: false,
		},
	}
	if err := ebiten.RunGame(&Game{0, start, 0}); err != nil {
		log.Fatal(err)
	}
}
