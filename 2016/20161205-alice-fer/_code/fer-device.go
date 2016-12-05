// Device is a handle to what users get to run via the Fer toolkit.
type Device interface {
	Configure(cfg config.Device) error
	Init(ctl Controler) error
	Run(ctl Controler) error
	Pause(ctl Controler) error
	Reset(ctl Controler) error
}

