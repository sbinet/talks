package main

import (
	"image/color"
)

// START OMIT

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(400, 400)               // sets the canvas dimensions in pixels
	p5.Background(color.Gray{Y: 220}) // sets canvas background color
}

func draw() {
	p5.Stroke(color.Black)              // sets the stroked path color
	p5.StrokeWidth(2)                   // sets the stroked path pen size
	p5.Fill(color.RGBA{R: 255, A: 208}) // sets the fill color of the "current" shape
	p5.Ellipse(50, 50, 80, 80)          // draws an ellipse at (50,50), with 80 major/minor-axis

	p5.Fill(color.RGBA{B: 255, A: 208})
	p5.Quad(50, 50, 80, 50, 80, 120, 60, 120) // draws a Bézier curve from (50,50) to (60,120)

	p5.Fill(color.RGBA{G: 255, A: 208})
	p5.Rect(200, 200, 50, 100)

	p5.TextSize(24)                   // sets the text size
	p5.Text("Hello, World!", 10, 300) // draws text at (10,300)
}

// END OMIT
