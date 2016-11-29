func (tsk *task2) Process(ctx fwk.Context) error {
	store := ctx.Store()
	// blocks until data for this event/slot is available
	v, err := store.Get(tsk.input) // HLxxx
	if err != nil {
		return err
	}
	i := v.(int64)
	o := tsk.fct(i)
	err = store.Put(tsk.output, o) // HLxxx
	return err
}
