package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	XDim = 1024
	YDim = 768

	Square = 32
	XLen   = (XDim / Square) + 1
	YLen   = (YDim / Square) + 1

	MinSupply = -10.0
	MaxSupply = 10.0
	DSupply   = MaxSupply - MinSupply
)

var supplies = make([]float64, XLen*YLen)

func loop(dt float64, cells []float64, colors []color.Color) *imdraw.IMDraw {
	dCells := make([]float64, len(cells))

	for x := 0; x < XLen; x++ {
		for y := 0; y < YLen; y++ {
			c := cells[x*YLen+y]
			if x > 0 {
				dCells[x*YLen+y] += cells[(x-1)*YLen+y] - c
			}
			if x < XLen-1 {
				dCells[x*YLen+y] += cells[(x+1)*YLen+y] - c
			}
			if y > 0 {
				dCells[x*YLen+y] += cells[x*YLen+y-1] - c
			}
			if y < YLen-1 {
				dCells[x*YLen+y] += cells[x*YLen+y+1] - c
			}
			dCells[x*YLen+y] *= 0.25
		}
	}

	for x := 0; x < XLen; x++ {
		for y := 0; y < YLen; y++ {
			n := cells[x*YLen+y]
			if supplies[x*YLen+y] != 0.0 {
				n = supplies[x*YLen+y]
			} else {
				n += dCells[x*YLen+y]
			}

			cells[x*YLen+y] = n

			hue := ((n - MinSupply) / DSupply) * 240.0
			switch {
			case hue < 0.0:
				colors[x*YLen+y] = pixel.RGB(1.0, 0, 0)
			case hue < 60.0:
				colors[x*YLen+y] = pixel.RGB(1.0, hue/60.0, 0)
			case hue < 120.0:
				colors[x*YLen+y] = pixel.RGB((120-hue)/60, 1.0, 0)
			case hue < 180.0:
				colors[x*YLen+y] = pixel.RGB(0.0, 1.0, (hue-120)/60.0)
			case hue < 240.0:
				colors[x*YLen+y] = pixel.RGB(0.0, (240-hue)/60.0, 1.0)
			default:
				colors[x*YLen+y] = pixel.RGB(0, 0, 1.0)
			}
		}
	}

	imd := imdraw.New(nil)
	for x := 0; x < XLen-1; x++ {
		for y := 0; y < YLen-1; y++ {
			imd.Color = colors[x*YLen+y]
			imd.Push(pixel.V(float64(x*Square), float64(y*Square)))
			imd.Color = colors[(x+1)*YLen+y]
			imd.Push(pixel.V(float64((x+1)*Square), float64(y*Square)))
			imd.Color = colors[(x+1)*YLen+y+1]
			imd.Push(pixel.V(float64((x+1)*Square), float64((y+1)*Square)))
			imd.Color = colors[x*YLen+y+1]
			imd.Push(pixel.V(float64(x*Square), float64((y+1)*Square)))
			imd.Polygon(0)
		}
	}

	return imd
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Heat Simulation",
		Bounds: pixel.R(0, 0, XDim, YDim),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	cells := make([]float64, XLen*YLen)
	colors := make([]color.Color, XLen*YLen)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// for i := 0; i < 10; i++ {
	// 	pos := rand.Intn(len(supplies))
	// 	for ; supplies[pos] != 0.0; pos = rand.Intn(len(supplies)) {}
	//
	// 	if i % 2 == 0 {
	// 		supplies[pos] = MinSupply
	// 	} else {
	// 		supplies[pos] = MaxSupply
	// 	}
	// }

	last := time.Now()
	for !win.Closed() {
		now := time.Now()
		dt := now.Sub(last).Seconds()
		last = now

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			p := win.MousePosition()
			pos := (int(p.X) / Square) * YLen + (int(p.Y) / Square)
			if supplies[pos] != 0.0 {
				supplies[pos] = 0.0
			} else {
				supplies[pos] = MaxSupply
			}
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			p := win.MousePosition()
			pos := (int(p.X) / Square) * YLen + (int(p.Y) / Square)
			if supplies[pos] != 0.0 {
				supplies[pos] = 0.0
			} else {
				supplies[pos] = MinSupply
			}
		}
		if win.JustPressed(pixelgl.MouseButtonMiddle) {
			supplies = make([]float64, len(supplies))
		}

		win.Clear(colornames.Aliceblue)
		imd := loop(dt, cells, colors)
		imd.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
