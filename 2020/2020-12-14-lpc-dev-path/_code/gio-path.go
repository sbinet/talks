package gio

import (
	"gioui.org/f32"
	"gioui.org/op"
)

type Path struct{}

func (p *Path) Begin(ops *op.Ops)             {} // start a path
func (p *Path) MoveTo(to f32.Point)           {} // move the pen
func (p *Path) LineTo(to f32.Point)           {} // draw a line
func (p *Path) QuadTo(ctl, to f32.Point)      {} // draw a quadratic Bézier curve
func (p *Path) Cube(ctl0, ctl1, to f32.Point) {} // draw a cubic Bézier curve
func (p *Path) Close()                        {} // close a path
