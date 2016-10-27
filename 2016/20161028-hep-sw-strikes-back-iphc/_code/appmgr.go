func (app *appmgr) run(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.state = fsm.Running

	switch app.nprocs {
	case 0:
		err = app.runSequential(ctx) // HLxxx
	default:
		err = app.runConcurrent(ctx) // HLxxx
	}

	return err
}

