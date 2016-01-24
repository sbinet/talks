add_to_concurrent_vector(concurrent_vector<double> v, double new_stuff[], size_t n) {
    auto my_iterator = v.grow_by(n);
    for (size_t i=0; i<n; ++i) {
        *my_iterator++ = new_stuff[i];
	}
}

