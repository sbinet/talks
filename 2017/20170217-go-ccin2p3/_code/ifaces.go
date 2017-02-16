package ifaces

type Stringer interface {
	String() string
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Empty interface{}
