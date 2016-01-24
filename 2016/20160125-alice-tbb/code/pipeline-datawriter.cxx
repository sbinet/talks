class DataWriter {
private:
    FILE* my_output;
public:
    DataWriter(FILE* out): my_output{out} {};
    void operator()(double const answer) const {
       fprintf(my_output, "%lf\n", answer);
    }
};
