// Closures in Rust
// functions in rust have a type, so we can assign them to variables
// fn function_name(parameter: type) -> return_type {}
// let foo = function_name;
// the type of foo is a function item not a function pointer. A function item is a zero sized value
// carried around  at compile time that references a unique function. It's zero sized because it's
// only used at compile time and this is why this fails
// let foo = function_name::<type>;
// foo = function_name::<another_type>; // fails because they are not the same type
// foo is the type of a function that takes a type while the other foo is the type of a function that
// takes another type
// function pointers are a type that can hold a reference to a function
// function items can be coerced to function pointers automatically by the compiler
// fn add_one(x: i32) -> i32 {
//     x + 1
// }
// let f_i = add_one; // this is a function item
// let f_p: fn(i32) -> i32 = add_one; // this is a function pointer
// if we have a function 
// fn baz(_: fn(i32) -> i32) {} // this is a function that takes a function pointer
// baz(f_i); // this works because the compiler coerces the function item to a function pointer
// baz(f_p); // this works because f_p is already a function pointer
// Closures are anonymous functions that can capture their environment
// let add_one = |x: i32| x + 1;
// let f = add_one;
// let f: fn(i32) -> i32 = add_one; // this fails because add_one is a closure not a function
// Closures can capture their environment
// let x = 1;
// let add_x = |y| x + y;
// let y = 2;
// let result = add_x(y);
// Closures can capture their environment by reference
// let x = 1;

#[cfg(test)]
mod tests {
}
