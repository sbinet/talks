#include <future>

double some_expensive_operation(std::vector<thing>& stuff) {
    // Do expensive things to stuff
    ...

    return result;
}

int main() {
    auto my_stuff = initialise_stuff();

    std::future<double> res = std::async(some_expensive_operation, std::ref(my_stuff));

    // do some other things

    double the_answer = res.get();
    ...

