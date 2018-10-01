// START CTL OMIT

// Controler controls devices execution and gives a device access to input and
// output data channels.
type Controler interface {
	Logger
	Chan(name string, i int) (chan Msg, error)
	Done() chan Cmd
}

// STOP CTL OMIT

// Logger gives access to printf-like facilities
type Logger interface {
	Fatalf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}

// START MSG OMIT

// Msg is a quantum of data being exchanged between devices.
type Msg struct {
	Data []byte // Data is the message payload.
	Err  error  // Err indicates whether an error occured.
}

// Cmd describes commands to be sent to a device, via a channel.
type Cmd byte

// STOP MSG OMIT
