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

	rss := map[string][]float64{
		"linear": []float64{
			16,
			16,
			32,
			64,
			128,
			256,
			512,
			1024,
		},
		"fads": []float64{
			16,
			16,
			16,
			35,
			52,
			84,
			149,
			232,
		},
	}

	//p.CheckedCmd("set multiplot")
	p.CheckedCmd("set title \"darwin - RSS (MB) - 10000 events\"")
	p.CheckedCmd("set grid x")
	p.CheckedCmd("set grid y")
	p.CheckedCmd("set xrange [-0.5:64.5]")
	//set yrange [-40:20]
	//p.CheckedCmd("set logscale y")

	p.SetXLabel("nbr of procs")
	p.SetYLabel("RSS (MB)")

	p.PlotXY(
		x,
		rss["linear"],
		"linear",
	)

	p.PlotXY(
		x,
		rss["fads"],
		"fads",
	)

	p.CheckedCmd("set style linespoints")
	p.CheckedCmd("set terminal png")
	p.CheckedCmd("set output 'darwin-rss.png'")
	p.CheckedCmd("replot")

	p.CheckedCmd("q")
	return
}
