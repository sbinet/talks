tbb::flow::graph g;

tbb::flow::broadcast_node<std::string> shout(g);
tbb::flow::function_node<std::string, std::string> en2fr(g, 1, translate_1);
tbb::flow::function_node<std::string, std::string> en2de(g, 1, translate_2);

tbb::flow::make_edge(shout, en2fr);
tbb::flow::make_edge(shout, en2de);
