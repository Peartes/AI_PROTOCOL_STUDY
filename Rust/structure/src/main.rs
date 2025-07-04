extern crate structure;
use structure::say_hello;


fn main() {
    /*
       Rust program that demonstrates a stack trace
       when a panic occurs. The program consists of multiple functions
       that call each other, leading to a panic in the deepest function.
       When the panic occurs, Rust will print a stack trace during unwinding
    */
    stack1();
}

fn stack1() {
    say_hello();
    stack2();
}

fn stack2() {
    stack3();
}

fn stack3() {
    panic!("This is a panic in stack3!");
}
