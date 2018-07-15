package main

import (
	"image"
	"os"

	_ "image/png"

	"github.com/direvus/sudoku"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

type scene int

const (
	start scene = iota
	end
	game
)

const (
	width = 900
)

var input bool
var state scene
var x, y int

var puzzle = sudoku.GenerateSolution()
var mask = puzzle.MinimalMask()
var board = puzzle.ApplyMask(&mask)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func updateBoard(puz *sudoku.Puzzle, value byte, imd *imdraw.IMDraw) {
	puz[9*x+y] = value
	imd.Clear()
	input = false
}

func run() {
	// initialize window
	cfg := pixelgl.WindowConfig{
		Title:  "Sudoku",
		Bounds: pixel.R(0, 0, width, width),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// initialize imdraw
	imd := imdraw.New(nil)

	// initialize font
	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 64 * width / 900,
	})
	atlas := text.NewAtlas(face, text.ASCII)

	// initialize batch
	batch := pixel.NewBatch(&pixel.TrianglesData{}, atlas.Picture())

	for !win.Closed() {
		switch state {
		// start menu
		case start:
			win.Clear(colornames.Snow)
			msg := text.New(pixel.ZV, atlas)
			msg.WriteString("Sudoku")
			msg.DrawColorMask(win,
				pixel.IM.Moved(win.Bounds().Center().Sub(msg.Bounds().Center())),
				colornames.Black)

			if win.JustPressed(pixelgl.MouseButtonLeft) ||
				win.JustPressed(pixelgl.KeyEnter) ||
				win.JustPressed(pixelgl.KeySpace) {
				state = game
			}
		// end menu
		case end:
			win.Clear(colornames.Snow)
			msg := text.New(pixel.ZV, atlas)
			msg.WriteString("Congratulations, You Win!")
			msg.DrawColorMask(win,
				pixel.IM.Moved(win.Bounds().Center().Sub(msg.Bounds().Center())),
				colornames.Black)

			if win.JustPressed(pixelgl.MouseButtonLeft) ||
				win.JustPressed(pixelgl.KeyEnter) ||
				win.JustPressed(pixelgl.KeySpace) {
				win.SetClosed(true)
			}
		// the actual game
		case game:
			// win condition
			if board.Equal(puzzle) {
				state = end
			}

			// select a box
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				if input {
					imd.Clear()
				}
				pos := win.MousePosition()
				x = int(pos.X) / (width / 9)
				y = int(pos.Y) / (width / 9)

				if !mask[9*x+y] {
					imd.Color = colornames.Paleturquoise
					imd.Push(pixel.V(float64(x*width/9)+1, float64(y*width/9)+1),
						pixel.V(float64((x+1)*width/9)-1, float64((y+1)*width/9)-1))
					imd.Rectangle(0)
					input = true
				}
			}
			// act on user input
			if input && !mask[9*x+y] {
				if win.JustPressed(pixelgl.Key1) || win.JustPressed(pixelgl.KeyKP1) {
					updateBoard(&board, '1', imd)
				}
				if win.JustPressed(pixelgl.Key2) || win.JustPressed(pixelgl.KeyKP2) {
					updateBoard(&board, '2', imd)
				}
				if win.JustPressed(pixelgl.Key3) || win.JustPressed(pixelgl.KeyKP3) {
					updateBoard(&board, '3', imd)
				}
				if win.JustPressed(pixelgl.Key4) || win.JustPressed(pixelgl.KeyKP4) {
					updateBoard(&board, '4', imd)
				}
				if win.JustPressed(pixelgl.Key5) || win.JustPressed(pixelgl.KeyKP5) {
					updateBoard(&board, '5', imd)
				}
				if win.JustPressed(pixelgl.Key6) || win.JustPressed(pixelgl.KeyKP6) {
					updateBoard(&board, '6', imd)
				}
				if win.JustPressed(pixelgl.Key7) || win.JustPressed(pixelgl.KeyKP7) {
					updateBoard(&board, '7', imd)
				}
				if win.JustPressed(pixelgl.Key8) || win.JustPressed(pixelgl.KeyKP8) {
					updateBoard(&board, '8', imd)
				}
				if win.JustPressed(pixelgl.Key9) || win.JustPressed(pixelgl.KeyKP9) {
					updateBoard(&board, '9', imd)
				}
				if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyKP0) ||
					win.JustPressed(pixelgl.KeyBackspace) || win.JustPressed(pixelgl.KeySpace) {
					updateBoard(&board, '0', imd)
				}
			}

			// set up lines for drawing
			for i := 1; i < 9; i++ {
				imd.Color = colornames.Black
				if i%3 == 0 {
					imd.Push(pixel.V(float64(i*width/9), 0), pixel.V(float64(i*width/9), width))
					imd.Line(6)

					imd.Push(pixel.V(0, float64(i*width/9)), pixel.V(width, float64(i*width/9)))
					imd.Line(6)
				} else {
					imd.Push(pixel.V(float64(i*width/9), 0), pixel.V(float64(i*width/9), width))
					imd.Line(3)

					imd.Push(pixel.V(0, float64(i*width/9)), pixel.V(width, float64(i*width/9)))
					imd.Line(3)
				}
			}

			// set up numbers for drawing
			batch.Clear()
			for i, sa := range board {
				num := text.New(pixel.ZV, atlas)
				num.WriteByte(sa)
				num.DrawColorMask(batch,
					pixel.IM.
						Scaled(
							pixel.ZV, float64(width)/900).
						Moved(
							pixel.V((float64(i/9)+0.3)*width/9,
								(float64(i%9)+0.25)*width/9)),
					colornames.Black)
			}

			// draw the scene to the window
			win.Clear(colornames.Snow)
			batch.Draw(win)
			imd.Draw(win)
		}

		// update the window
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
