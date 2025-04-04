#![allow(dead_code)]

use std::{fmt::Debug, io, marker::Tuple, vec};
// I have thought of trait bounds as a way to specify the types that a generic can be.
//     * For example, if we have a generic type T, we can specify that T must implement the trait Foo.
//     * This means that T can be any type that implements the Foo trait.
//     * This is useful because it allows us to write generic code that can work with any type that implements the Foo trait.
//     * For example, we can write a function that takes a generic type T and a trait bound Foo.
//     * This means that T can be any type that implements the Foo trait.
pub trait Foo {}
pub type Bar = u32;
impl Foo for Bar {}
// here we are using the trait bound Foo to specify that T must implement the trait Foo
pub fn bar<T>(_ty: T)
where
    T: Foo,
{
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_bar() {
        bar(1); // this works because Bar implements the Foo trait
                // bar("hello"); // this doesn't work because String doesn't implement the Foo trait
        assert!(true);
    }
}

// but trait bounds are more than that. Think of them as a way to specify what your item requires
// they don't have to be of the type T: trait they can be any type that is even non-local e.g. String
// and don't even need to include generic types. the generic params don't even have to appear on the left hand
// of the function. a trait bound String: Clone is valid and means that the type String must implement the trait Clone
// some valid trait bounds
// this trait bound means that T must implement the trait Clone
fn foo<T: Clone>(_t: T) {}

// this trait bound means that the error type must implement the From trait
// this means that the error type must be convertible from the type T
fn foo2<T>(_t: T)
where
    io::Error: From<T>,
{
}

// trait bounds can also be used on associated types using the T::(ass_ty): trait
// if there is any ambiguity, use the qualified syntax <Type as T>::Item: Trait
trait Input
where
    Self::Item: Debug,
{
    type Item;
}

struct Garbage {}
struct NotDebuggable {}

impl Input for Garbage {
    type Item = String; // works because String impl Debug
                        // type Item = NotDebuggable;
}
// let's express a trait that can only be implemented for a type that impl the IntoIterator
// one way to do this is using an associated item as shown below
trait TakeNext {
    type Item: IntoIterator;
    fn take_next(_ty: Self::Item) -> <Self::Item as IntoIterator>::Item;
}

// now when impl this trait, we must use an item that is iterable
impl TakeNext for Vec<Garbage> {
    type Item = Vec<String>;

    fn take_next(_ty: Self::Item) -> String {
        "string".to_string()
    }
}
// but a better way is using supertrait although this only expresses that the trait is impl for types that impl IntoIterator
trait TakeNext2: IntoIterator {}
//impl TakeNext2 for Garbage {} // this will not work because Garbage does not impl IntoIterator

// another approach would be to use a blanket impl of the trait with a trait bound on the implementation
trait TakeNext3 {}
// we are saying that any item of type T that impl IntoIterator and FnOnce<Tuple> can be used as a TakeNext3
impl<T: FnOnce(String, String) + IntoIterator> TakeNext3 for T {}

impl IntoIterator for Garbage {
    type Item = String;

    type IntoIter = vec::IntoIter<Self::Item>;

    fn into_iter(self) -> Self::IntoIter {
        todo!()
    }
}

impl FnOnce<(String, String)> for Garbage {
    type Output = (); // for object safety, the return type has to be a () i.e. the call does not return any value
                      // the definition of the call method makes a reference to Self (Self::Output) which breaks object safety

    extern "rust-call" fn call_once(self, _args: (String, String)) -> Self::Output {
        ()
    }
}

#[cfg(test)]
mod test_trait_bound {
    use super::*;

    #[test]
    fn test_take_next() {
        let _ty: &dyn TakeNext3 = &Garbage {}; // works because Vec<T> impl IntoIterator
    }
}
