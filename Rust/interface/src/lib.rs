use serde::ser::SerializeStruct;
use std::fmt::Debug;

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
        // test our type is serializable
        let serialized = serde_json::to_string(&ty).unwrap();
        assert_eq!(serialized, r#"{"i_name":"test","i_value":42}"#);
    }
}
