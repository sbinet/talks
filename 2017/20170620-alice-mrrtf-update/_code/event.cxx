const int ARRAYSZ = 10;

struct P3 {
	int32_t Px; double  Py; int32_t Pz;
};

struct Event {
	TString  Beg;
	int16_t  I16; int32_t  I32; int64_t  I64; float    F32; double   F64; 
	P3       P3;
	std::string StlStr;

	std::vector<int16_t> StlVecI16; std::vector<int32_t> StlVecI32; // ...
	std::vector<std::string> StlVecStr;

	int16_t ArrayI16[ARRAYSZ]; int32_t ArrayI32[ARRAYSZ]; // ...

	int32_t N;
	int16_t *DynArrayI16; //[N]
	// ...
};
