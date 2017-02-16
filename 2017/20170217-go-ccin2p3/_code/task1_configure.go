import "reflect"

// configure START OMIT
func (tsk *task1) Configure(ctx fwk.Context) error {
	err := tsk.DeclOutPort(tsk.i1prop, reflect.TypeOf(int64(1.0))) // HLxxx
	if err != nil {
		return err
	}
	// ...
	return err
}

// configure END OMIT
