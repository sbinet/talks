class DataReader {
private:
    FILE *my_input;

public:
    DataReader(FILE* in):
        my_input{in} {};

    DataReader(const DataReader& a):
        my_input{a.my_input} {};

    ~DataReader() {};

    double operator()(tbb::flow_control& fc) const {
        double number;
        int rc = fscanf(my_input, "%lf\n", &number);
        if (rc != 1) {
            fc.stop();
            return 0.0;
        }
        return number;
    }
};
