use proc_macros::{PrintIndent, print_fn_name_attr};

pub trait PrintIndent {
    fn print_indent(&self);
}

#[derive(PrintIndent)]
pub struct Rectangle {
    pub width: u32,
    pub height: u32,
}


#[cfg(test)]
mod tests {
    use super::*;
    
    #[print_fn_name_attr]
    pub fn test() {
        println!("Hello, world!");
    }
    #[test]
    fn it_works() {
        let r = Rectangle {
            width: 30,
            height: 50,
        };
        r.print_indent();
        test();
    }
}
