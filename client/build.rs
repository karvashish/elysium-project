use std::path::Path;
use std::process::Command;
use std::{env, process};

/*
This build script performs the following tasks:

1. Ensures the Cargo build process reruns this script if specific environment variables change.
2. Validates that mandatory environment variables are set; if any are missing, it terminates the build process with an error message.
3. Parses and validates mandatory environment variables to ensure they are in the correct format (e.g., IP addresses and CIDR values).
4. Handles optional environment variables like `CLIENTPUB`, logging whether they are set or not.
5. Sets a default value for `IFCNAME` if it is not provided.
6. Embeds the environment variable values into the binary at compile time by using the `cargo:rustc-env` directive.
7. Compiles a C source file (`wireguard.c`) into a static library (`libwireguard.a`) using GCC, and links it with the Rust code.
8. Informs Cargo to link the generated static library and sets up the required linker search paths.
9. Ensures that changes to the `src/wireguard/wireguard.c` file trigger a rebuild.
10. Adds logic to set default environment variables when building in debug mode but requires these variables in release mode.**
*/

fn main() {
    let is_release = env::var("PROFILE").unwrap_or_default() == "release";

    for var in ["CLIENTPUB", "IFCNAME", "ADDR", "CIDR", "SERVERPUB", "SERVERENDPOINT", "SERVERIP"] {
        println!("cargo:rerun-if-env-changed={}", var);
    }

    let required_vars = ["ADDR", "CIDR", "SERVERPUB", "SERVERENDPOINT", "SERVERIP"];

    let mut missing_vars = vec![];

    for var in &required_vars {
        if env::var(var).is_err() {
            missing_vars.push(*var);
        }
    }

    if is_release {
        if !missing_vars.is_empty() {
            eprintln!("Error: Missing mandatory environment variables in release mode: {:?}", missing_vars);
            process::exit(1);
        }
    } else {
        if env::var("ADDR").is_err() {
            env::set_var("ADDR", "192.168.1.2");
        }
        if env::var("CIDR").is_err() {
            env::set_var("CIDR", "24");
        }
        if env::var("SERVERPUB").is_err() {
            env::set_var("SERVERPUB", "default-server-pub-key");
        }
        if env::var("SERVERENDPOINT").is_err() {
            env::set_var("SERVERENDPOINT", "default.endpoint.com:51820");
        }
        if env::var("SERVERIP").is_err() {
            env::set_var("SERVERIP", "192.168.1.1");
        }
    }

    let addr = env::var("ADDR").unwrap().parse::<std::net::Ipv4Addr>().expect("Invalid ADDR");
    let cidr = env::var("CIDR").unwrap().parse::<u8>().expect("Invalid CIDR");
    let server_pub = env::var("SERVERPUB").unwrap();
    let endpoint = env::var("SERVERENDPOINT").unwrap();
    let server_ip = env::var("SERVERIP").unwrap().parse::<std::net::Ipv4Addr>().expect("Invalid SERVERIP");

    let client_pub = env::var("CLIENTPUB").ok();
    if let Some(key) = &client_pub {
        println!("Using CLIENTPUB: {}", key);
    } else {
        println!("CLIENTPUB is not set, proceeding without it");
    }

    let ifc_name = env::var("IFCNAME").unwrap_or_else(|_| "wg0".to_string());
    println!("Using IFCNAME: {}", ifc_name);

    println!("Using ADDR: {}/{}", addr, cidr);
    println!("Using SERVERPUB: {}", server_pub);
    println!("Using SERVERENDPOINT: {}", endpoint);
    println!("Using SERVERIP: {}", server_ip);

    if let Some(key) = client_pub {
        println!("cargo:rustc-env=CLIENTPUB={}", key);
    }
    println!("cargo:rustc-env=IFCNAME={}", ifc_name);
    println!("cargo:rustc-env=ADDR={}", addr);
    println!("cargo:rustc-env=CIDR={}", cidr);
    println!("cargo:rustc-env=SERVERPUB={}", server_pub);
    println!("cargo:rustc-env=SERVERENDPOINT={}", endpoint);
    println!("cargo:rustc-env=SERVERIP={}", server_ip);

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
