package main

// START CFG OMIT

import (
	"log" // OMIT
	// OMIT
	"github.com/sbinet-alice/fer"
	"github.com/sbinet-alice/fer/config"
)

type processor struct {
	cfg    config.Device
	idatac chan fer.Msg // HL
	odatac chan fer.Msg // HL
}

func (dev *processor) Configure(cfg config.Device) error {
	dev.cfg = cfg
	return nil
}

func (dev *processor) Init(ctl fer.Controler) error {
	idatac, err := ctl.Chan("data1", 0) // handle err // HL
	odatac, err := ctl.Chan("data2", 0) // handle err // HL
	dev.idatac = idatac
	dev.odatac = odatac
	return nil
}

// STOP CFG OMIT

// START RUN OMIT

func (dev *processor) Run(ctl fer.Controler) error {
	for {
		select {
		case data := <-dev.idatac: // HL
			out := append([]byte(nil), data.Data...)
			out = append(out, []byte(" (modified by "+dev.cfg.Name()+")")...)
			dev.odatac <- fer.Msg{Data: out} // HL
		case <-ctl.Done():
			return nil
		}
	}
}

func (dev *processor) Pause(ctl fer.Controler) error {
	return nil
}

func (dev *processor) Reset(ctl fer.Controler) error {
	return nil
}

// STOP RUN OMIT

// START MAIN OMIT

func main() {
	err := fer.Main(&processor{}) // HL
	if err != nil {
		log.Fatal(err)
	}
}

// STOP MAIN OMIT
