// FFI is a way to call other code written in different programming languages.
// It also allows us to expose Rust functions to other languages.
// This is made possible by using the `extern` keyword.
// This keyword tells the Rust compiler that the symbols defined in this block
// should be available for linking from other languages (during the linking phase).

/*
    FFI is, ultimately, all about accessing bytes that originate somewhere out-
    side your application’s Rust code. For that, Rust provides two primary build-
    ing blocks: symbols, which are names assigned to particular addresses in a
    given segment of your binary that allow you to share memory (be it for data
    or code) between the external origin and your Rust code, and calling conven-
    tions that provide a common understanding of how to call functions stored
    in such shared memory. We’ll look at each of these in turn
*/
#[unsafe(no_mangle)]
/*
    The #[no_mangle] attribute ensures that RS_DEBUG retains that name dur-
    ing compilation rather than having the compiler assign it another symbol
    name to, for example, distinguish it from another (non-FFI) RS_DEBUG static
    variable elsewhere in the program.
*/
pub static RS_DEBUG: bool = true;
 unsafe extern "C" {
static FOREIGN_DEBUG: bool;
}
