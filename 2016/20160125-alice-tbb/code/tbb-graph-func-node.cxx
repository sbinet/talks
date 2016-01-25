tbb::flow::function_node< float, float > squarer( 
		g, 
		tbb::flow::unlimited, 
		[](const float &v) -> float {
			return v*v; 
		} 
);
