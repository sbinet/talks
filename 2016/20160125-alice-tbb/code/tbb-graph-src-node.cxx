class file_cat {
private:
    std::ifstream* m_input_stream_p;
public:
    file_cat(std::ifstream* input_stream_p):
        m_input_stream_p(input_stream_p) {};

    bool operator() (std::string& msg) {
        *m_input_stream_p >> msg;
        if (m_input_stream_p->good())
            return true;
        return false;
    }
};

// ...

tbb::flow::source_node<std::string> input_node(g, file_cat(my_stream_p), false);

// ...

input_node.activate();

