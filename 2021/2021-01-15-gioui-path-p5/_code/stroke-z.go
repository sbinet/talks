package main

import (
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

// START OMIT
func Z(o *op.Ops) {

	p := new(clip.Path)
	p.Begin(o)
	p.Move(f32.Pt(50, 20))
	p.Line(f32.Pt(200, 0))
	p.Line(f32.Pt(-200, 200))
	p.Line(f32.Pt(200, 0))
	p.Quad(f32.Pt(-200, 80), f32.Pt(-200, 200))
	p.Line(f32.Pt(200, 0))

	clip.Stroke{
		Path: p.End(),
		Style: clip.StrokeStyle{
			Width: 2.5,
			Cap:   clip.SquareCap,
			Join:  clip.BevelJoin,
		},
	}.Op().Add(o)

	paint.Fill(o, red)
}

// END OMIT
