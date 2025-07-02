use serde::ser::SerializeStruct;
use std::{borrow::Cow, fmt::Debug};

/// proper interface definition in Rust should be unsurprising. Unsurprising means that
/// things that are automatically intuitive should be. One way to achieve this
/// is designing using concepts users are already familiar with and some ways
/// to achieve this are naming conventions, common traits and ergonomic traits
///
/// common traits that makes it easy for people to use your interfaces
/// some common traits to impl are defined below
///
pub struct MyInterface {
    pub i_name: String,
    pub i_value: u32,
}
/// Where possible, we should again avoid surprising the
/// user and eagerly implement most of the standard traits even if we do not
/// need them immediately.
///
/// One of the trait to impl is the Debug trait. Almost all users expect your types
/// to be printable using the `{:?}`. You could also derive the trait
impl Debug for MyInterface {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "MyInterface {{ i_name: {}, i_value: {} }}",
            self.i_name, self.i_value
        )
    }
}
/// another trait is Send and Sync. Send means the type can be safely used in multiple threads
/// A type that is not Sync can’t be shared through an Arc or placed
/// in a static variable.
unsafe impl Send for MyInterface {}

unsafe impl Sync for MyInterface {}

/// One step further down in the hierarchy of expected traits is the comparison traits: PartialEq, PartialOrd, Hash, Eq, and Ord. The PartialEq trait is
// particularly desirable, because users will at some point inevitably have two
// instances of your type that they wish to compare with == or assert_eq!.
impl PartialEq for MyInterface {
    fn eq(&self, other: &Self) -> bool {
        self.i_name == other.i_name && self.i_value == other.i_value
    }
}
/// Finally, for most types, it makes sense to implement the serde crate’s
/// Serialize and Deserialize traits
impl serde::Serialize for MyInterface {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let mut state = serializer.serialize_struct("MyInterface", 2)?;
        state.serialize_field("i_name", &self.i_name)?;
        state.serialize_field("i_value", &self.i_value)?;
        state.end()
    }
}

/// Rust does not auto impl traits for references to types that impl a trait. e.g. even though we impl
/// Serialize for MyInterface, we need to impl Serialize for &MyInterface. This is because some trait method
/// might take ownership of the type or require an exclusive reference to the type.

/// wrapper types in Rust allow for some form of inheritance inn Rust. Wrapper types like AsRef and Deref allow you to call methods of a type U
/// on a type T that impl AsRef<U> or Deref<Target=U>
pub struct MyInterfaceWrapper {}

/// allow MyInterfaceWrapper to be used as a MyInterface or in a more technical sense, allow MyInterfaceWrapper reference to be used as a MyInterface reference
impl AsRef<MyInterface> for MyInterfaceWrapper {
    fn as_ref(&self) -> &MyInterface {
        todo!()
    }
}
// we can call methods of MyInterface on MyInterfaceWrapper
impl MyInterface {
    pub fn my_interface_name(&self) {
        println!("MyInterface name is {}", self.i_name);
    }
}

mod flexible {
    use super::*;
    /*
        FLEXIBLE INTERFACES

        - Avoid restricting callers more than you must
        - Accept generic arguments whenever reasonable
          * This lets the caller give you the cheapest, simplest type they have
          * Example: impl AsRef<str> instead of &str or String

        - Return types can also be flexible
          * Cow<'_, str> is a great way to let the implementation decide
            if a value is borrowed or owned

        - In function signatures:
            fn do_stuff(s: String) -> String
                - expects caller to allocate and you to allocate
                - hard to evolve without breaking

            fn do_stuff(s: &str) -> Cow<'_, str>
                - caller can pass by ref
                - you can return borrowed or owned
                - easier to change later

            fn do_stuff<T: AsRef<str>>(s: T) -> impl AsRef<str>
                - even more flexible
                - maximizes what can be passed and what can be returned

        - This is critical because relaxing restrictions is backwards compatible,
          but tightening them is a breaking change.

        - Rust empowers you to design this kind of flexibility
          through generics, trait bounds, lifetimes, and ownership choices.

        💡 Takeaway: be minimal in restrictions, maximal in promises.
    */
    // Write a function process_lines that takes anything that can be turned into a str,
    // and returns something that can be used as a str, while potentially allocating if needed.
    #[allow(dead_code)]
    pub fn process_lines<'a>(input: &'a impl AsRef<str>) -> Cow<'a, str> {
        // In the body, if the string is less than 10 characters, return it as-is (borrowed); otherwise, return an uppercase owned string.
        let input_str = input.as_ref();
        if input_str.len() < 10 {
            Cow::Borrowed(input_str) // return as owned string
        } else {
            Cow::Owned(input_str.to_uppercase()) // return as owned string
        }
    }
}

#[allow(dead_code)]
mod generics {
    use std::fmt::Debug;

    /*
        GENERIC ARGUMENTS

        - Generics let you write one piece of code that works for many concrete types
        - They increase flexibility *and* safety:
            * The compiler will verify that only compatible types can be used
            * You don't have to write code for each type manually

        - Use trait bounds to express what your generic types are capable of
            Example:
                fn compare<T: PartialOrd>(a: T, b: T) -> Ordering

        - Generics can apply to:
            * Types (structs, enums, etc.)
            * Functions
            * Traits

        - You can even make trait bounds generic themselves, via higher-ranked trait bounds.

        - Generics also work with lifetimes and const generics in modern Rust:
            * lifetimes: tie reference lifetimes together
            * const generics: parameterize types over compile-time values (array sizes, etc)

        - 💡 Takeaway: design with generics so your interfaces are future-proof,
          reusable, and composable, while still type-safe.
    */
    // A generic compare_and_print function that takes any two comparable items and prints which one is larger, using a PartialOrd bound.
    pub fn compare_and_print<T: PartialOrd + Debug>(input1: &T, input2: &T) {
        if input1 < input2 {
            println!("{:?} is less than {:?}", input1, input2);
        } else if input1 > input2 {
            println!("{:?} is greater than {:?}", input1, input2);
        } else {
            println!("{:?} is equal to {:?}", input1, input2);
        }
    }

    /*
        HASVALUE TRAIT

        - Demonstrates using an associated type to access inner fields generically
        - Any type implementing HasValue promises a `.value()` method returning
          something PartialOrd + Debug
    */
    pub trait HasValue {
        type Value: PartialOrd + Debug;

        fn value(&self) -> &Self::Value;
    }

    // generic function comparing inner values
    pub fn compare_inner_values<T: HasValue>(a: &T, b: &T) {
        if a.value() < b.value() {
            println!("{:?} is less than {:?}", a.value(), b.value());
        } else if a.value() > b.value() {
            println!("{:?} is greater than {:?}", a.value(), b.value());
        } else {
            println!("{:?} is equal to {:?}", a.value(), b.value());
        }
    }

    // sample type
    pub struct NumberWrapper {
        pub value: i32,
    }

    impl HasValue for NumberWrapper {
        type Value = i32;

        fn value(&self) -> &Self::Value {
            &self.value
        }
    }
}

#[cfg(test)]
mod tests {

    use super::*;

    #[test]
    fn it_works() {
        let ty = MyInterface {
            i_name: "test".to_string(),
            i_value: 42,
        };
        // test our type is debuggable
        assert_eq!(
            format!("{:?}", ty),
            "MyInterface { i_name: test, i_value: 42 }"
        );
        // test our type is comparable
        let ty2 = MyInterface {
            i_name: "test".to_string(),
            i_value: 42,
        };
        assert_eq!(ty, ty2);
    }

    #[test]
    fn test_flexible_interface() {
        // test the flexible interface
        let input = "short";
        let result = flexible::process_lines(&input);
        assert_eq!(result, Cow::Borrowed("short"));

        let input = "this is a long string";
        let result = flexible::process_lines(&input);
        assert_eq!(
            result,
            Cow::Owned::<str>("THIS IS A LONG STRING".to_string())
        );
    }

    #[test]
    fn test_generic_compare_and_print() {
        // test the generic compare_and_print function
        let a = 5;
        let b = 10;
        generics::compare_and_print(&a, &b);
        let c = "hello";
        let d = "world";
        generics::compare_and_print(&c, &d);
        // We can't assert the output directly, but we can visually check it in the console.
    }

    #[test]
    fn test_compare_inner_values() {
        let a = generics::NumberWrapper { value: 5 };
        let b = generics::NumberWrapper { value: 10 };
        generics::compare_inner_values(&a, &b);
        // again, visually check console
    }
}
