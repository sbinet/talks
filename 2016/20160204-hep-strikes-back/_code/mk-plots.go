//+build ignore

package main

import (
	"fmt"
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
)

const (
	nevts = 10000
)

type Data struct {
	nprocs float64
	cpu    float64
	rss    float64
}

type Datas []Data

func (d Datas) CPU() plotter.XYs {
	o := make(plotter.XYs, len(d))
	for i := range d {
		o[i].X = d[i].nprocs
		o[i].Y = d[i].cpu
	}
	return o
}

func (d Datas) RSS() plotter.XYs {
	o := make(plotter.XYs, len(d))
	for i := range d {
		o[i].X = d[i].nprocs
		o[i].Y = d[i].rss / 1000.0
	}
	return o
}

func (d Datas) Hz() plotter.XYs {
	o := make(plotter.XYs, len(d))
	for i := range d {
		o[i].X = d[i].nprocs
		o[i].Y = 1.0 / (d[i].cpu / nevts)
	}
	return o
}

type Results struct {
	fads    Datas
	delphes Datas
}

/*
	darwin := []Data{
		{0, 828575 / 1000.0, 16000},
		{1, 828241 / 1000.0, 16000},
		{2, 419421 / 1000.0, 16000},
		{4, 215883 / 1000.0, 35000},
		{8, 115512 / 1000.0, 52000},
		{16, 103179 / 1000.0, 84000},
		{32, 100532 / 1000.0, 149000},
		{64, 101248 / 1000.0, 232000},
	}
*/

var (
	data Results
)

func main() {
	fads := Datas([]Data{
		{0, 654566 / 1000.0, 18896},
		{1, 653227 / 1000.0, 17900},
		{2, 338351 / 1000.0, 25280},
		{4, 171694 / 1000.0, 38280},
		{8, 87464 / 1000.0, 60072},
		{16, 50631 / 1000.0, 86896},
		{32, 52236 / 1000.0, 119028},
		{40, 50342 / 1000.0, 126540},
		{64, 48167 / 1000.0, 154476},
	})

	delphes := Datas([]Data{
		{0, 645780 / 1000.0, 63676},
		{1, 645780 / 1000.0, 63676},
		{2, 645780 / 1000.0, 2 * 63676},
		{4, 645780 / 1000.0, 4 * 63676},
		{8, 645780 / 1000.0, 8 * 63676},
		{16, 645780 / 1000.0, 8 * 63676},
		{32, 645780 / 1000.0, 8 * 63676},
		{40, 645780 / 1000.0, 8 * 63676},
		{64, 645780 / 1000.0, 8 * 63676},
	})

	data = Results{
		fads:    fads,
		delphes: delphes,
	}

	do_linux_cpu("lhcb3")
	do_linux_hz("lhcb3")
	do_linux_rss("lhcb3")

}

func do_linux_cpu(name string) {
	// Make a plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = fmt.Sprintf("linux - %d events", nevts)
	p.X.Label.Text = "nbr of procs"
	p.X.Max = 70
	p.Y.Label.Text = "CPU time (s)"
	p.Y.Min = 0
	p.Y.Max = 800
	// p.Y.Scale = plot.LogScale{}
	p.Legend.Top = true

	p.Add(plotter.NewGrid())
	err = plotutil.AddLinePoints(p,
		"fads", data.fads.CPU(),
		"delphes", data.delphes.CPU(),
	)
	if err != nil {
		panic(err)
	}

	// naive linear scaling
	naive := plotter.NewFunction(func(x float64) float64 {
		switch x {
		case 0:
			return data.fads.CPU()[1].Y
		default:
			return data.fads.CPU()[1].Y / x
		}
	})
	naive.Color = color.RGBA{B: 255, A: 255}
	p.Add(naive)
	p.Legend.Add("naive-scaling", naive)

	// save the plot
	err = p.Save(6, 4, fmt.Sprintf("%s-cpu.pdf", name))
	if err != nil {
		panic(err)
	}
	return
}

func do_linux_hz(name string) {
	// Make a plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = fmt.Sprintf("linux - %d events", nevts)
	p.X.Label.Text = "nbr of procs"
	p.X.Max = 70
	p.Y.Label.Text = "Event rate (Hz)"
	p.Y.Min = 0
	p.Y.Max = 250
	// p.Y.Scale = plot.LogScale{}
	p.Legend.Top = true

	p.Add(plotter.NewGrid())
	err = plotutil.AddLinePoints(p,
		"fads", data.fads.Hz(),
		"delphes", data.delphes.Hz(),
	)

	// naive linear scaling
	naive := plotter.NewFunction(func(x float64) float64 { return x * data.fads.Hz()[1].Y })
	naive.Color = color.RGBA{B: 255, A: 255}

	p.Add(naive)
	p.Legend.Add("naive-scaling", naive)

	// save the plot
	err = p.Save(6, 4, fmt.Sprintf("%s-hz.pdf", name))
	if err != nil {
		panic(err)
	}
	return
}

func do_linux_rss(name string) {
	// Make a plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = fmt.Sprintf("linux - %d events", nevts)
	p.X.Label.Text = "nbr of procs"
	p.X.Max = 70
	p.Y.Label.Text = "RSS (MB)"
	p.Y.Min = 0
	//p.Y.Scale = plot.LogScale{}
	//p.Y.Tick.Marker = myDefaultTicks{3,}
	//p.Legend.Top = true
	//p.Legend.Left = true

	p.Add(plotter.NewGrid())
	err = plotutil.AddLinePoints(p,
		"fads", data.fads.RSS(),
		"delphes", data.delphes.RSS(),
	)
	if err != nil {
		panic(err)
	}

	// save the plot
	err = p.Save(6, 4, fmt.Sprintf("%s-rss.pdf", name))
	if err != nil {
		panic(err)
	}
	return
}
