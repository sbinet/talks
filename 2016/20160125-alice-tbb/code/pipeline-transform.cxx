class Transform {
public:
    double operator()(double const number) const {
        double answer=0.0;
        if (number > 0.0)
            answer = some_expensive_calculation(number)
        return answer;
    }
};
