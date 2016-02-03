//+build ignore

package main

import (
	"fmt"

	gnuplot "github.com/sbinet/go-gnuplot"
)

func main() {
	fname := ""
	persist := true
	debug := true

	p, err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		panic(fmt.Errorf("err: %v\n", err))
	}
	defer p.Close()

	x := []float64{0, 1, 2, 4, 8, 16, 32, 64}

	cpu := map[string][]float64{
		"linear": []float64{
			826.78,
			826.78,
			413.39,
			206.695,
			103.3475,
			51.67375,
			25.836875,
			12.9184375,
		},
		"fads": []float64{
			826.78,
			825.78,
			419.27,
			216.02,
			117.47,
			103.90,
			100.37,
			101.28,
		},
	}

	//p.CheckedCmd("set multiplot")
	p.CheckedCmd("set title \"darwin - CPU time (s) - 10000 events\"")
	p.CheckedCmd("set grid x")
	p.CheckedCmd("set grid y")
	p.CheckedCmd("set xrange [-0.5:64.5]")
	//set yrange [-40:20]
	//p.CheckedCmd("set logscale y")

	p.SetXLabel("nbr of procs")
	p.SetYLabel("CPU time (s)")

	p.PlotXY(
		x,
		cpu["linear"],
		"linear",
	)

	p.PlotXY(
		x,
		cpu["fads"],
		"fads",
	)

	p.CheckedCmd("set style linespoints")
	p.CheckedCmd("set terminal png")
	p.CheckedCmd("set output 'darwin-cpu.png'")
	p.CheckedCmd("replot")

	p.CheckedCmd("q")
	return
}
