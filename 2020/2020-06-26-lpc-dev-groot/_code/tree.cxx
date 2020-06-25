void write() {
	auto f = TFile::Open("out.root", "RECREATE");
	auto t = new TTree("t", "title"); // HL

	int32_t n = 0;
	double  px = 0;
	double  arr[10];
	double  vec[20];

	t->Branch("n",   &n,  "n/I"); // HL
	t->Branch("px",  &px, "px/D");
	t->Branch("arr", arr, "arr[10]/D");
	t->Branch("vec", vec, "vec[n]/D"); // HL

	for (int i = 0; i < NEVTS; i++) {
		// fill data: n, px, arr, vec with some values
		fill_data(&n, &px, &arr, &vec);

		t->Fill(); // commit data to tree. // HL
	}

	f->Write(); // commit data to disk.
	f->Close();

	exit(0);
}
