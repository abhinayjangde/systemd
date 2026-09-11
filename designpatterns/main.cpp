#include<iostream>
using namespace std;

// defining the abstract class (blueprint) - cannot create an object of it
class Employee {
    public:
    virtual void work() = 0;
    virtual ~Employee() = default;
    void greet(){
        cout << "welcome to company"<<endl;
    }
};

class Developer : public Employee {
    public:
    void work() override {
        cout << "writing and debugging code..."<<endl;
    }
};

int main() {
    Developer emp;
    
    emp.greet();
    emp.work();

    Employee *empPtr = new Developer();
    empPtr->work();

    delete empPtr;
    return 0;
}