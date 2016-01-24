#include <tbb/tbb.h>

double serial_sum(double x[], size_t n) {
    double the_answer = 0.0;
    for (size_t i=0; i<n; ++i) { the_answer += x[i]; }
    return the_answer;
}

class parallel_sum {
    double *my_x;
public:
    double the_answer;
    void operator()(const tbb::blocked_range<size_t>& r) {
        double *x=my_x;
        for(size_t i=r.begin(); i!=r.end(); ++i) { the_answer += x[i]; }
    }
    // This is the constructor used to split the task
    parallel_sum(parallel_sum& a, tbb::split): my_x{a.my_x}, the_answer{0.0} {};
    // This method joins (or merges) the results of two subtasks
    void join(const parallel_sum& b) {
        the_answer += b.the_answer;
    }
    parallel_sum(double x[]): my_x{x} {};
};

