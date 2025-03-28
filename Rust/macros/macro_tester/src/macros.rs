use proc_macro2::TokenStream;
use syn::DeriveInput;
use quote::quote;

pub fn json_derive(item: TokenStream) -> TokenStream {
    let ast: DeriveInput = syn::parse2(item).unwrap();
    let name = &ast.ident;
    quote! {
        impl Json for #name {
            fn to_json(&self) -> String {
                let json = to_string(self).unwrap();
            }
        }
    }
}