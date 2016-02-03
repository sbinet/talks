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

	cpu := map[string][]float64{
		"delphes": []float64{50.01, 50.01, 50.01},
		"fads":    []float64{90.81, 90.63, 48.71},
	}

	//p.CheckedCmd("set multiplot")
	p.CheckedCmd("set title \"linux - CPU time (s) - 500 events\"")
	p.CheckedCmd("set grid x")
	p.CheckedCmd("set grid y")
	p.CheckedCmd("set xrange [-0.5:2.5]")
	//set yrange [-40:20]
	p.CheckedCmd("set style linespoints")

	p.SetXLabel("nbr of procs")
	p.SetYLabel("CPU time (s)")

	p.PlotX(
		//[]float64{0, 1, 2},
		cpu["delphes"],
		"delphes",
	)

	p.PlotX(
		//[]float64{0, 1, 2},
		cpu["fads"],
		"fads",
	)

	p.CheckedCmd("set terminal png")
	p.CheckedCmd("set output 'linux-cpu.png'")
	p.CheckedCmd("replot")

	p.CheckedCmd("q")
	return
}
