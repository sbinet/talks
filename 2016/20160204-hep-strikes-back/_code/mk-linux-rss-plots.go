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

	rss := map[string][]float64{
		"delphes": []float64{55632, 55632, 55632},
		"fads":    []float64{17724, 16688, 25548},
	}

	//p.CheckedCmd("set multiplot")
	p.CheckedCmd("set title \"linux - RSS (kB) - 500 events\"")
	p.CheckedCmd("set grid x")
	p.CheckedCmd("set grid y")
	p.CheckedCmd("set xrange [-0.5:2.5]")
	//set yrange [-40:20]
	p.CheckedCmd("set style linespoints")

	p.SetXLabel("nbr of procs")
	p.SetYLabel("RSS (kB)")

	p.PlotX(
		rss["delphes"],
		"delphes",
	)

	p.PlotX(
		rss["fads"],
		"fads",
	)

	p.CheckedCmd("set terminal png")
	p.CheckedCmd("set output 'linux-rss.png'")
	p.CheckedCmd("replot")

	p.CheckedCmd("q")
	return
}
