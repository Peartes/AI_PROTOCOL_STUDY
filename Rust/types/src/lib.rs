#![feature(unboxed_closures, fn_traits, impl_trait_in_assoc_type)]
mod trait_bound; 
// let's explore types in Rust
// one important concept for how types are represented in memory using Rust is alignment
// alignment is the requirement that data be stored in memory at addresses that are multiples of a certain value
// this value is called the alignment of the type
// for example, a u32 type has an alignment of 4 bytes, meaning it must be stored at an address that is a multiple of 4
// this is important for performance reasons, as accessing data that is not aligned can be slower

/*
    * The `Positions` struct represents the positions of different entities in a game.
    * It contains four fields: `me`, `you`, `them`, and `us`, all of which are 32-bit unsigned integers.
    * The alignment of each field is 4 bytes, which is the same as the size of a u32.
    *
    * This struct can be used to store the positions of different players or entities in a game,
    * allowing for efficient access and manipulation of their positions.

    * The `#[repr(C)]` attribute is used to ensure that the struct has a C-compatible layout,
    * which means that the fields will be laid out in memory in the same order as they are declared.
    * This is important for interoperability with C code, as it ensures that the struct can be passed
    * between Rust and C without any issues.

    * The layout of the struct is as follows:
    * +------+-------+------+------+
    * |  me  |  you  | them |  us  |
    * +------+-------+------+------+
    * |  4   |   4   |  4   |  4   |
    * +------+-------+------+------+
    * The total size of the struct is 16 bytes, which is the sum of the sizes of all four fields.
    * The alignment of the struct is 4 bytes, which is the same as the alignment of its fields.
    * The size of the struct must also be aligned to the size of its largest field, which is 4 bytes in this case.
    * If for example we had a u64 field and a u8, the size of the struct would be 17 bytes and the alignment would be 8 bytes.
    pub struct Positions {
        pub me: u8, // this is a u32 so the alignment is 1 bytes
        pub you: u32, // this is a u32 so the alignment is 4 bytes
        pub them: u32, // this is a u32 so the alignment is 4 bytes
        pub us: u64, // this is a u64 so the alignment is 8 bytes
    }
    * Because the largest field is a u64, the alignment of the struct is 8 bytes.
    * This means that the struct must be stored at an address that is a multiple of 8.
    * The size of the struct is 1 + 3 + 4 + 4 + 4 + 8 = 24 bytes.
    * The layout of the struct is as follows:
    * +------+-------+------+------+
    * |  me  |  you  | them |  us  |
    * +------+-------+------+------+
    * |  1   |   4   |  4   |  8   |
    * +------+-------+------+------+
    * The total size of the struct is 24 bytes, which is the sum of the sizes of all four fields plus the padding.
*/
#[repr(C)]
pub struct Positions {
    pub me: u32, // this is a u32 so the alignment is 4 bytes
    pub you: u32, // this is a u32 so the alignment is 4 bytes
    pub them: u32, // this is a u32 so the alignment is 4 bytes
    pub us: u32, // this is a u32 so the alignment is 4 bytes
    // size of the struct is 4 + 4 + 4 + 4 = 16 bytes
    // the alignment of the struct is 4 bytes
    // so 0 bytes of padding will be added to the struct to make it 16 bytes
}

#[repr(C)]
pub struct LargePositions {
    pub me: u8, // this is a u32 so the alignment is 1 bytes
    pub you: u16, // this is a u16 so the alignment is 2 bytes
    pub them: u32, // this is a u32 so the alignment is 4 bytes
    pub us: u64, // this is a u64 so the alignment is 8 bytes
    // size of the struct is 1 + 1 + 2 + 4 + 8 = 16 bytes
    // the alignment of the struct is 8 bytes
}

/*
    * Rust memory layout is a bit different from C/C++ memory layout.
    * In Rust, the layout of a struct is determined by the order of its fields and their types.
    * The layout of the struct is as follows:
    * +------+-------+------+------+
    * |  us  | them  | you  |  me  |
    * +------+-------+------+------+
    * |  8   |   4   |  2   |  1   |
    * +------+-------+------+------+
    * The total size of the struct is 8 + 4 + 2 + 1 = 15 bytes, since alignment is 8bytes, 
    * size of the struct is 16bytes
    * This is because rust optimizes the layout and can re-arrange the fields in memory to minimize padding.
    * This means the fields ordering is not guaranteed but less memory is wasted.
    * For this example, rust arranges the fields from last to first because it avoids all forms of padding
    * Other types might be different but this is a simple deterministic example.

*/
#[repr(Rust)]
pub struct RustLargePositions {
    pub me: u8, // this is a u32 so the alignment is 1 bytes
    pub you: u16, // this is a u16 so the alignment is 2 bytes
    pub them: u32, // this is a u32 so the alignment is 4 bytes
    pub us: u64, // this is a u64 so the alignment is 8 bytes
    // size of the struct is 1 + 1 + 2 + 4 + 8 = 16 bytes
    // the alignment of the struct is 8 bytes
}

#[cfg(test)]
mod tests {
    use std::{alloc::{alloc, Layout}, mem, ptr};

    use super::*;

    #[test]
    fn positions() {
        assert_eq!(mem::size_of::<Positions>(), 16);
        assert_eq!(mem::align_of::<Positions>(), 4);

        assert_eq!(mem::size_of::<LargePositions>(), 16);
        assert_eq!(mem::align_of::<LargePositions>(), 8);

        // because of the #[repr(C)] attribute, the layout of the struct is guaranteed to be the same as the C layout
        // this means that the fields will be laid out in memory in the same order as they are declared
        // let's check the layout of the struct
        let positions = Positions {
            me: 1,
            you: 2,
            them: 3,
            us: 4,
        };
        // let's get a raw pointer to the struct
        unsafe {
            let raw_position = &positions as *const Positions as *const u32;
            // let's get the position of the fields in memory
            let me_pointer = alloc(Layout::new::<u32>()) as *mut u32;
            ptr::copy(raw_position, me_pointer, 1);
            assert_eq!(*me_pointer, 1);
            let you_pointer = alloc(Layout::new::<u32>()) as *mut u32;
            ptr::copy(raw_position.add(1), you_pointer, 1);
            assert_eq!(*you_pointer, 2);
            let them_pointer = alloc(Layout::new::<u32>()) as *mut u32;
            ptr::copy(raw_position.add(2), them_pointer, 1);
            assert_eq!(*them_pointer, 3);
            let us_pointer = alloc(Layout::new::<u32>()) as *mut u32;
            ptr::copy(raw_position.add(3), us_pointer, 1);
            assert_eq!(*us_pointer, 4);
        }
    }

    #[test]
    fn large_positions() {
        let large_positions = LargePositions {
            me: 1,
            you: 2,
            them: 3,
            us: 4,
        };
        // let's get a raw pointer to the struct
        unsafe {
            let raw_large_position = &large_positions as *const LargePositions as *const u8;
            // let's get the position of the fields in memory
            assert_eq!(*raw_large_position, 1);
            assert_eq!(*(raw_large_position.add(1) as *const u16), 2);
            assert_eq!(*(raw_large_position.add(4) as *const u32), 3);
            assert_eq!(*(raw_large_position.add(8) as *const u64), 4);
        }
    }

    #[test]
    fn rust_large_positions() {
        let rust_large_positions = RustLargePositions {
            me: 1,
            you: 2,
            them: 3,
            us: 4,
        };
        assert_eq!(mem::size_of::<RustLargePositions>(), 16);
        assert_eq!(mem::align_of::<RustLargePositions>(), 8);
        // let's get a raw pointer to the struct
        unsafe {
            // first get a pointer to the struct. we cast this as a *const u64 because rust makes the us field the first field
            let raw_rust_large_position = &rust_large_positions as *const RustLargePositions as *const u64;
            // let's get the position of the fields in memory
            assert_eq!(*raw_rust_large_position, 4);
            // the next field is the them field, which is a u32
            assert_eq!(*(raw_rust_large_position.add(1) as *const u32), 3);
            // the next field is the you field, which is a u16
            assert_eq!(*(((raw_rust_large_position as *const u32).add(3)) as *const u16), 2);
            // the last field is the me field, which is a u8
            assert_eq!(*(((raw_rust_large_position as *const u16).add(7)) as *const u8), 1);
        }
    }
}
