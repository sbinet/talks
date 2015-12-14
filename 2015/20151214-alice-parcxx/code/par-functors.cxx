#include <thread>
#include <iostream>

class FuncObj
{
public:
	void operator()(void) {
		std::cout << std::this_thread::get_id() << std::endl;
	}
};

int main(int argc, char **argv) {
	FuncObj f;
	std::thread t(f);
	t.join();
	return 0;
}

