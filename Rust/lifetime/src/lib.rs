pub struct Rectangle<'a, 'b> {
    width: &'a u32,
    height: &'b u32,
}

impl<'a, 'b> Rectangle<'a, 'b> {
    pub fn area(&self) -> Box<u32> {
        let res = Box::new(self.width * self.height);
        res
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let width = 10;
        let height = 20;
        let rect = Rectangle {
            width: &width,
            height: &height,
        };
        let area = rect.area();
        assert_eq!(*area, 200);
    }
}