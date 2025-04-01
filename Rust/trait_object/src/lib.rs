use std::mem;

// For the below cases, the compiler will generate code for each type of s that is passed to the function.
// so from our main function, we will have 2 strlen functions generated for type &str and String.
// and 1 strlen2 function is generated for type &str
// This is called monomorphization.
pub trait Hei {
    fn hei(&self);
}

impl Hei for String {
    fn hei(&self) {
        println!("Hei from String");
    }
}

impl Hei for &str {
    fn hei(&self) {
        println!("Hei from str");
    }
}

// compiler will generate code for each concrete type that is passed to the function
pub fn say_hei_static<H: Hei>(h: H) {
    h.hei();
}

pub fn strlen(s: impl AsRef<str>) -> usize {
    s.as_ref().len()
}

pub fn strlen2<S>(s: S) -> usize
where
    S: AsRef<str>,
{
    s.as_ref().len()
}

// In this case, the compiler does does not know the concrete type of h, this is called type erasure
// so the compiler will generate a fat pointer for the trait object, which contains the pointer to the concrete and the pointer to the vtable
// the vtable contains the pointers to the functions that are implemented for the trait
// vtable struct looks similar to this
/*
struct Vtable {
    destructor: fn(*mut u8), // to free memory
    size: usize, // size of the concrete type
    align: usize,
    drop_in_place: fn(*mut u8), // to drop the concrete type
    hei: fn(*const u8),
}
*/
// the vtable is created at compile time and is stored in the binary
pub fn say_hei(h: Box<dyn Hei>) {
    h.hei();
}

pub fn main_() {
    let s = "hello"; // s: &'static str
    let s2 = String::from("world"); // s2: String
    println!("{}", strlen(s));
    println!("{}", strlen(s2.clone()));
    println!("{}", strlen2(s));

    say_hei_static(s);
    say_hei_static(s2.clone());
    say_hei(Box::new(s));
    say_hei(Box::new(s2));
}

// For multiple trait objects, the compiler cannot as of now generate fat pointers for each trait hence
// it requires that you create a new trait that combines the other traits
// and it creates a vtable for the new trait with pointers to all the methods of the combined traits
pub trait Combined: Hei + AsRef<str> {}

pub fn say_combined(c: &dyn Combined) {
    c.hei();
    let s = c.as_ref();
    println!("{}", s.len());
}

// Now there is a marker called Self:Sized that is used as a trait bound to ensure that the implementing
// type must be sized. This can be useful for methods that require the type to be sized. If a trait is used as a
// trait object and some of the methods require the type to be sized, then that trait is not object safe and
// cannot ber used as a trait object
/*
pub trait Test {
    fn test<U>(param: U);
}

pub fn fail(_: &dyn Test) {
    // this will not compile because the type is not sized
}
*/
// to use such trait, we need to mark the methods that require the type to be sized Self: Sized
// that way the compiler will not generate code for the method when the type is not sized
/*
pub trait Test {
    fn test<U>(param: U)
    where
        Self: Sized;
}

pub fn pass(_: &dyn Test) {
    // this will compile because the type is sized
}
*/

// if you want a trait to be implemented on only sized types, you can use mark the trait as Sized
pub trait SizedTrait
where
    Self: Sized,
{
    fn size(self); // this takes a reference so it does not require that the type be sized
}

// this works because the type is sized
impl SizedTrait for String {
    fn size(self) {
        println!("{}", mem::size_of_val(&self)); // this will print the size of the string struct
    }
}

// this will not compile if we had it as an str instead of &str
// impl SizedTrait for str {
impl SizedTrait for &str {
    fn size(self) {
        println!("{}", mem::size_of_val(self)); // this will print the size of the pointer to the str
    }
}

#[test]
pub fn test_sized_trait() {
    let s = "hello world";
    let s2 = String::from("hello world");
    s.size();
    s2.size();
}

// Trait objects must be object safe and that means 4 things
// 1. The trait cannot have generic type parameters because those are monomorphized at compile time
// and the compiler cannot generate code for all the types that could be passed to the function
// 2. The trait cannot have unconstrained associated types i.e. the assocaited types must be fixed and not
// depend on an external type e.g impl<T> Trait for Type<T> { type associated T; } this breaks because dynamic dispatch
// erases the concrete type of the trait object hence the associated type cannot be determined
// 3. The trait cannot return Self because dynamic dispatch erases the concrete type of the trait object
// hence the return type cannot be determined
// 4. The trait cannot have static methods because they are not associated with any instance of the trait
// hence they cannot be called on the trait object. The compiler needs to know the concrete type of the trait
// except if the methods are marked as requiring sized types with Self: Sized, then the compiler ignores the method
// when the type is not sized. The compiler needs to know the concrete type of the trait
/*
pub trait NotObjectSafe {
    fn object_safe();
}

pub fn test_not_object_safe(_: &dyn ObjectSafe) {
    // this will not compile because the trait has a static method
}

pub trait ObjectSafe {
    fn object_safe()
    where
        Self: Sized;
}

pub fn test_object_safe(_: &dyn ObjectSafe) {
    // this will compile because the trait has a method that requires the type to be sized
}
*/

// so it knows the vtalbe to use to call the method
// 5. The trait cannot have associated constants because they are not associated with any instance of the trait
// hence they cannot be called on the trait object. The compiler needs to know the concrete type of the trait

/*
pub trait ObjectSafe {
    const VALUE: i32;
    fn object_safe(&self);
}

pub fn test_object_safe(_: &dyn ObjectSafe) {
    // this will not compile because the trait has an associated constant
}
*/

// note that even though the method of an object safe trait can take self as a value
/*
pub trait ObjectSafe {
    fn object_safe(self);
}
// the trait can only be implemented if the concrete type implementing the trait is sized
impl ObjectSafe for String {
    // this works because String is sized
    fn object_safe(self) {}
}

impl ObjectSafe for str {
    // this will not compile because str is not sized
    fn object_safe(self) {}
}
*/
