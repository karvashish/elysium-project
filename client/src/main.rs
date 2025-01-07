use std::net::Ipv4Addr;

mod interface;
mod wg_common;

use wg_common::{
    wg_common::{gen_private_key, gen_public_key, list_device_names, update_device},
    wireguard_cffi::WgKeyBase64String,
};

use interface::{create_wireguard_ifc, update_wireguard_ifc, Operation};

/*
This is the main entry point for the Elysium Project Client setup. It performs the following tasks:

1. **Environment Variable Retrieval and Parsing**:
   - Retrieves compile-time constants embedded by the build script.
   - Ensures mandatory variables are correctly parsed (e.g., IP addresses and CIDR values).
   - Optionally retrieves `CLIENTPUB` if provided.

2. **Logging Setup Information**:
   - Logs the retrieved environment variables to verify correctness.

3. **WireGuard Key Pair Generation**:
   - Generates a new private key and corresponding public key using the WireGuard CFFI interface.

4. **WireGuard Interface Management**:
   - Creates the WireGuard interface with the name specified by `IFCNAME`.
   - Updates the interface with the provided address, CIDR, and configuration.
   - Updates the device's configuration with the generated private key and port number.
   - Enables the interface.

5. **Error Handling**:
   - Handles errors at each step, logging specific failures if any operation does not complete successfully.

6. **List Available WireGuard Interfaces**:
   - Lists all available WireGuard interfaces at the end of the setup process.
*/
#[tokio::main(flavor = "current_thread")]
async fn main() {
    println!("Starting Elysium Project Client setup");

    const IFCNAME: &str = env!("IFCNAME");
    const ADDR: &str = env!("ADDR");
    const CIDR: &str = env!("CIDR");
    const SERVERPUB: &str = env!("SERVERPUB");
    const SERVERENDPOINT: &str = env!("SERVERENDPOINT");
    const SERVERIP: &str = env!("SERVERIP");

    let client_pub = option_env!("CLIENTPUB");
    if let Some(key) = client_pub {
        println!("CLIENTPUB: {}", key);
    } else {
        println!("CLIENTPUB is not set");
    }

    let addr = ADDR
        .parse::<Ipv4Addr>()
        .expect("Invalid IPv4 address in ADDR");
    let cidr = CIDR.parse::<u8>().expect("Invalid CIDR value in CIDR");

    println!("Interface Name: {}", IFCNAME);
    println!("Address: {}/{}", addr, cidr);
    println!("Server Public Key: {}", SERVERPUB);
    println!("Server endpoint: {}", SERVERENDPOINT);

    let private_key: WgKeyBase64String = gen_private_key();
    let public_key: WgKeyBase64String = gen_public_key(&private_key);

    println!("New Private Key: {:?}", private_key);
    println!("New Public Key: {:?}", public_key);

    let server_ip_cidr = format!("{}/{}", SERVERIP, 32);
    let server_ip: &str = &server_ip_cidr;

    match (
        create_wireguard_ifc(IFCNAME).await,
        update_wireguard_ifc(IFCNAME, Some(addr), Some(cidr), Operation::Update).await,
        update_device(
            &private_key,
            54161,
            IFCNAME,
            SERVERPUB,
            SERVERENDPOINT,
            server_ip,
        ),
        update_wireguard_ifc(IFCNAME, None, None, Operation::Enable).await,
    ) {
        (Ok(()), Ok(()), Ok(()), Ok(())) => {
            println!("Interface setup completed successfully.");
        }
        (Err(e1), _, _, _) => {
            eprintln!("Interface creation failed: {}", e1);
        }
        (Ok(()), Err(e2), _, _) => {
            eprintln!("Interface update failed: {}", e2);
        }
        (Ok(()), Ok(()), Err(e3), _) => {
            eprintln!("Device update failed: {}", e3);
        }
        (Ok(()), Ok(()), Ok(()), Err(e4)) => {
            eprintln!("Interface enabling failed: {}", e4);
        }
    }

    println!(
        "Available WireGuard interfaces: {}",
        list_device_names().join(", ")
    );
}
