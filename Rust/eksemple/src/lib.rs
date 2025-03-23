// For the below cases, the compiler will generate code for each type of s that is passed to the function.
// so from our main function, we will have 2 strlen functions generated for type &str and String.
// and 1 strlen2 function is generated for type &str
// This is called monomorphization.
pub fn strlen(s: impl AsRef<str>) -> usize {
    s.as_ref().len()
}

pub fn strlen2<S>(s: S) -> usize
where
    S: AsRef<str>,
{
    s.as_ref().len()
}

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

pub fn say_hei_static<H: Hei>(h: H) {
    h.hei();
}

pub fn say_hei(h: Box<dyn Hei>) {
    h.hei();
}

pub fn main() {
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
