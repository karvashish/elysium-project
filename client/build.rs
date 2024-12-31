use std::path::Path;
use std::process::Command;
use std::{env, process};

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
