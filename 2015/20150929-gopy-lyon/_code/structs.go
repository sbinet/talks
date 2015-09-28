package structs

type Struct struct {
	private int
	Public int
}

func (s Struct) Private() int {
	return s.private
}

func NewStruct(v int) Struct {
	return Struct{
		private: v+1,
		Public: v,
	}
}
