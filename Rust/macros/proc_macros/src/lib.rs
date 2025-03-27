use proc_macro::TokenStream;
use quote::quote;
use syn::{DeriveInput, ItemFn};

#[proc_macro_derive(PrintIndent)]
pub fn print_indent_derive(item: TokenStream) -> TokenStream {
    let ast: DeriveInput = syn::parse(item).unwrap();

    impl_print_indent(&ast)
}

fn impl_print_indent(ast: &DeriveInput) -> TokenStream {
    let name = &ast.ident;

    quote! {
        impl PrintIndent for #name {
            fn print_indent(&self) {
                println!("Hello, I am a {}", stringify!(#name));
            }
        }
    }
    .into()
}

#[proc_macro_attribute]
pub fn print_fn_name_attr(_attr: TokenStream, item: TokenStream) -> TokenStream {
    let ast = syn::parse_macro_input!(item as ItemFn);
    let vis = ast.vis;
    let sig = ast.sig;
    let sig_ident = &sig.ident;
    let block = ast.block;
    quote! {
        #vis #sig {
            println!("calling fn: {}", stringify!(#sig_ident));
            #block
        }
    }
    .into()
}
