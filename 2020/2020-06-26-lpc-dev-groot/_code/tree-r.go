package main

import (
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	// START OMIT
	f, err := groot.Open("f.root")
	defer f.Close()
	o, err := f.Get("t")

	v := struct {
		N int32 `groot:"n"`
		D struct {
			I32 int32   `groot:"i32"`
			I64 int64   `groot:"i64"`
			F64 float64 `groot:"f64"`
		} `groot:"d"`
	}{}

	r, err := rtree.NewReader(o.(rtree.Tree), rtree.ReadVarsFromStruct(&v))
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error { // HL
		fmt.Printf(
			"evt=%d, n=%d, d.i32=%d, d.i64=%d, d.f64=%v\n",
			ctx.Entry, v.N, v.D.I32, v.D.I64, v.D.F64,
		)
		return nil
	})
	// END OMIT
}
