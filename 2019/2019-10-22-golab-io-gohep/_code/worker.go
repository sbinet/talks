type worker struct {
	slot  int
	keys  []string // nodes in data-flow (Input/Output)
	store datastore
	ctxs  []context // a Context for each component
	msg   msgstream

	evts <-chan int64    // channel of event indices
	quit <-chan struct{} // channel to notify early-abort
	done chan<- struct{} // channel to notify we are done
	errc chan<- error    // channel of errors during event processing
}

