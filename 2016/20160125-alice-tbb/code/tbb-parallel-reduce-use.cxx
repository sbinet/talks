double Apply_parallel_sum(const double x[], size_t n) {
    parallel_sum ps(x);
    parallel_reduce(tbb::blocked_range<size_t>(0, n), ps);
    return ps.the_answer;
}

