void RunPipeline(int ntoken, FILE* input_file, FILE* output_file) {
    tbb::parallel_pipeline(
        ntoken, // maximum number of pieces of data that can be processed simultaneously
        tbb::make_filter<void,double>(
            tbb::filter::serial_in_order, DataReader(input_file) )
    &
        tbb::make_filter<double,double>(
            tbb::filter::parallel, Transform() )
    &
        tbb::make_filter<double,void>(
            tbb::filter::serial_in_order, DataWriter(output_file) );
}
