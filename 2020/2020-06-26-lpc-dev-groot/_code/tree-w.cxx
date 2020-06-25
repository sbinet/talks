auto t = new TTree("t", "my tree");
auto n = int32_t(0);
auto d = struct{
    int32_t i32;
	int64_t i64;
	double  f64;
};
t->Branch("n", &n, "n/I"); // HL
t->Branch("d", &d);        // HL

// -> leaf_n = TLeaf<int32_t>(t, "n");
// -> leaf_d = TLeaf<struct> (t, "d");
