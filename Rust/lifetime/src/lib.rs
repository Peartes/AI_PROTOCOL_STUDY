#[allow(dead_code)]
pub struct Rectangle<'a, 'b> {
    width: &'a mut u32,
    height: &'b mut u32,
}

impl<'a, 'b> Rectangle<'a, 'b> {
    pub fn area(&self) -> Box<u32> {
        let res = Box::new(*self.width * *self.height);
        res
    }

    pub fn extend(&mut self, other: &Self) {
        *(self.width) += *other.width;
        *(self.height) += *other.height;
    }
}

impl PartialEq for Rectangle<'_, '_> {
    fn eq(&self, other: &Self) -> bool {
        self.width == other.width && self.height == other.height
    }
}

#[cfg(test)]
mod tests {
    use std::alloc::Layout;

    use super::*;

    #[test]
    fn it_works() {
        let mut width = 10;
        let mut height = 20;
        let mut rect = Rectangle {
            width: &mut width,
            height: &mut height,
        };
        let area = rect.area();
        assert_eq!(*area, 200);

        let mut new_width = Box::new(5);
        let mut new_height = 10;

        let new_rect = unsafe { std::alloc::alloc(Layout::new::<Rectangle>()) } as *mut Rectangle;

        unsafe {
            std::ptr::write(
                new_rect,
                Rectangle {
                    width: &mut new_width,
                    height: &mut new_height,
                },
            );
        }
        let read_rect = unsafe { std::ptr::read(new_rect) };
        assert_eq!(*read_rect.width, 5);
        assert_eq!(*read_rect.height, 10);
        rect.extend(&read_rect);
        assert_eq!(*rect.width, 15);
        assert_eq!(*rect.height, 30);
    }
}
