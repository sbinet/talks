// START RUN OMIT
import "github.com/sbinet-alice/fer/config"

// Device is a handle to what users get to run via the Fer toolkit.
type Device interface {
	Run(ctl Controler) error
}

// STOP RUN OMIT

// START CFG OMIT
type DevConfigurer interface {
	Configure(cfg config.Device) error
}
type DevIniter interface {
	Init(ctl Controler) error
}
type DevPauser interface {
	Pause(ctl Controler) error
}
type DevReseter interface {
	Reset(ctl Controler) error
}

// STOP CFG OMIT
