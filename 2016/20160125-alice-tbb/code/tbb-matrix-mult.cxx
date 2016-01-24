void operator()( const tbb::blocked_range2d<size_t>& r ) const {
    float (*a)[L] = my_a;
    float (*b)[N] = my_b;
    float (*c)[N] = my_c;
    for( size_t i=r.rows().begin(); i!=r.rows().end(); ++i ){
        for( size_t j=r.cols().begin(); j!=r.cols().end(); ++j ) {
            float sum = 0;
            for( size_t k=0; k<L; ++k )
                sum += a[i][k]*b[k][j];
            c[i][j] = sum;
        }
    }
}

