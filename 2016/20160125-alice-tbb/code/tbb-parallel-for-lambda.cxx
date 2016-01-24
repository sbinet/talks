#include "tbb/tbb.h"

void ParallelApplyFunc(double x[], size_t n) {
    tbb::parallel_for(tbb::blocked_range<size_t>(0, n),
                        [=](const blocked_range<size_t>& r) {
                            for(size_t i=r.begin(); i!=r.end(); ++i)
                            my_func(x[i]);
                        });
}

