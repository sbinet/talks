package mpl

func Data(n int) []float64 {
	o := make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = float64(i * i)
	}
	return o
}
