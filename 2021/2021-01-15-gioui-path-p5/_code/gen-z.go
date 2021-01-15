// +build ignore

package main

import (
	"image/color"
	"image/png"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/app/headless"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

const (
	xsz = 500
	ysz = 500
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(
			unit.Dp(xsz),
			unit.Dp(ysz),
		))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case key.Event:
			switch e.Name {
			case key.NameEscape, "Q":
				screenshot()
				w.Close()
			}

		case system.FrameEvent:
			ops := new(op.Ops)
			do(ops)
			e.Frame(ops)
		}
	}
}

func screenshot() {
	ops := new(op.Ops)
	win, err := headless.NewWindow(xsz, ysz)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	do(ops)
	err = win.Frame(ops)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	img, err := win.Screenshot()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	f, err := os.Create("out.png")
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	err = f.Close()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
}

var (
	red   = color.NRGBA{R: 255, A: 255}
	black = color.NRGBA{A: 255}
)

func do(o *op.Ops) {
	newPath := newZigZagPath
	//	newPath = newLine

	{
		stk := op.Push(o)
		p := newPath(o)
		clip.Stroke{
			Path: p,
			Style: clip.StrokeStyle{
				Width: 20,
				Cap:   clip.FlatCap,
				Join:  clip.BevelJoin,
				Miter: 5,
			},
		}.Op().Add(o)
		paint.Fill(o, red)
		stk.Pop()
	}
	{
		stk := op.Push(o)
		p := newPath(o)
		clip.Stroke{
			Path: p,
			Style: clip.StrokeStyle{
				Width: 4,
			},
		}.Op().Add(o)
		paint.Fill(o, black)
		stk.Pop()
	}
}

func newLine(o *op.Ops) clip.PathSpec {
	p := new(clip.Path)
	p.Begin(o)
	p.Move(f32.Pt(50, 50))
	p.Line(f32.Pt(100, 0))
	p.Line(f32.Pt(0, 100))
	return p.End()
}

func newZigZagPath(o *op.Ops) clip.PathSpec {
	const x = 4
	p := new(clip.Path)
	p.Begin(o)
	p.Move(f32.Pt(10*x, 5*x))
	p.Line(f32.Pt(50*x, 0*x))
	p.Line(f32.Pt(-50*x, 50*x))
	p.Line(f32.Pt(50*x, 0))
	p.Quad(f32.Pt(-50*x, 20*x), f32.Pt(-50*x, 50*x))
	p.Line(f32.Pt(50*x, 0))
	return p.End()
}
