use std::path::Path;
use std::process::Command;
use std::{env, process};

// This build script performs several tasks during the build process:
// 1. It monitors specific environment variables (PUBKEY, IFCNAME, ADDR, CIDR) and ensures
//    that Cargo re-runs this script if any of them change.
// 2. It retrieves and validates these environment variables:
//    - PUBKEY: Optional, used for embedding a public key into the binary.
//    - IFCNAME: Optional, defaults to "wg0" if not set.
//    - ADDR and CIDR: Mandatory, representing an IPv4 address and CIDR subnet mask.
//      The script terminates the build with an error if either is missing or invalid.
// 3. After validation, the values of these variables are embedded into the binary at compile time
//    using `cargo:rustc-env`, ensuring the application doesn't depend on runtime environment variables.
// 4. The script compiles a C source file (`wireguard.c`) into a static library using GCC and AR tools:
//    - The C file is compiled into an object file (`wireguard.o`).
//    - The object file is archived into a static library (`libwireguard.a`).
// 5. Cargo is instructed to link the Rust project with the static library and re-run this script if
//    the `wireguard.c` file changes.
fn main() {
    println!("cargo:rerun-if-env-changed=PUBKEY");
    println!("cargo:rerun-if-env-changed=IFCNAME");
    println!("cargo:rerun-if-env-changed=ADDR");
    println!("cargo:rerun-if-env-changed=CIDR");
    let pub_key = env::var("PUBKEY").ok();
    if let Some(key) = &pub_key {
        println!("Using PUBKEY: {}", key);
    } else {
        println!("PUBKEY is not set, proceeding without it");
    }

    let ifc_name = env::var("IFCNAME").unwrap_or_else(|_| "wg0".to_string());
    println!("IFCNAME is not set, proceeding with default: {}", ifc_name);

    let addr = env::var("ADDR")
        .ok()
        .and_then(|addr| addr.parse::<std::net::Ipv4Addr>().ok());
    let cidr = env::var("CIDR")
        .ok()
        .and_then(|cidr| cidr.parse::<u8>().ok());

    if let (Some(addr), Some(cidr)) = (addr, cidr) {
        println!("Using ADDR: {}/{}", addr, cidr);
    } else {
        eprintln!("Error: Both ADDR and CIDR must be set and valid");
        process::exit(1);
    }

    if let Some(key) = pub_key {
        println!("cargo:rustc-env=PUBKEY={}", key);
    }
    println!("cargo:rustc-env=IFCNAME={}", ifc_name);
    println!("cargo:rustc-env=ADDR={}", addr.unwrap());
    println!("cargo:rustc-env=CIDR={}", cidr.unwrap());

    //------------------------------------------------------------------------------//
    let out_dir = env::var("OUT_DIR").unwrap();

    Command::new("gcc")
        .args(&["src/wireguard/wireguard.c", "-c", "-fPIC", "-o"])
        .arg(&format!("{}/wireguard.o", out_dir))
        .status()
        .unwrap();
    Command::new("ar")
        .args(&["crus", "libwireguard.a", "wireguard.o"])
        .current_dir(&Path::new(&out_dir))
        .status()
        .unwrap();

    println!("cargo::rustc-link-search=native={}", out_dir);
    println!("cargo::rustc-link-lib=static=wireguard");
    println!("cargo::rerun-if-changed=src/wireguard/wireguard.c");
}
