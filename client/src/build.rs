use std::env;

fn main() {
    let my_var = env::var("PUBKEY").unwrap_or_else(|_| "default_value".to_string());
    
    println!("cargo:rustc-env=SECRET={}", my_var);
}
