auto f = TFile::Open("out.root", "READ");
auto t = f->Get<TTree>("t");

auto n = int32_t(0);
auto d = struct{
    int32_t i32;
	int64_t i64;
	double  f64;
};

t->SetBranchAddress("n", &n); // HL
t->SetBranchAddress("d", &d); // HL

for (int64_t i = 0; i < t->GetEntries(); i++) {
    t->GetEntry(i);
	printf("evt=%d, n=%d, d.i32=%d, d.i64=%d, d.f64=%f\n",
			i, n, d.i32, d.i64, d.f64);
}
