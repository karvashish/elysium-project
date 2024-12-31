use std::net::Ipv4Addr;

mod interface;
mod wg_common;

use wg_common::{
    wg_common::{gen_private_key, gen_public_key},
    wireguard_cffi::WgKeyBase64String,
};


// The main function serves as the entry point for setting up the Elysium Project Client. 
// It performs the following tasks:
// 1. Retrieves and validates compile-time constants embedded via the build script:
//    - IFCNAME: The name of the WireGuard interface.
//    - ADDR: The IPv4 address assigned to the interface.
//    - CIDR: The subnet mask associated with the address.
//    - PUBKEY: An optional public key, if provided during the build process.
// 2. Parses and ensures the validity of the ADDR and CIDR values. These are required to configure 
//    the WireGuard interface properly.
// 3. Logs the retrieved and validated values to provide feedback about the configuration being used.
// 4. Generates a new private and public key pair for WireGuard using cryptographic utilities.
//    These keys will be used for securing communications over the interface.
#[tokio::main]
async fn main() {
    println!("Starting Elysium Project Client setup");

    const IFCNAME: &str = env!("IFCNAME");
    const ADDR: &str = env!("ADDR");
    const CIDR: &str = env!("CIDR");

    let pub_key = option_env!("PUBKEY");
    if let Some(key) = pub_key {
        println!("PUBKEY: {}", key);
    } else {
        println!("PUBKEY is not set");
    }

    let addr = ADDR
        .parse::<Ipv4Addr>()
        .expect("Invalid IPv4 address in ADDR");
    let cidr = CIDR.parse::<u8>().expect("Invalid CIDR value in CIDR");

    println!("Interface Name: {}", IFCNAME);
    println!("Address: {}/{}", addr, cidr);

    let private_key: WgKeyBase64String = gen_private_key();
    let public_key: WgKeyBase64String = gen_public_key(&private_key);

    print!("New Priv key: {:?}\n", private_key);
    print!("New Public key: {:?}\n", public_key);

    match (
        interface::create_wireguard_ifc(IFCNAME).await,
        interface::update_wireguard_ifc(
            IFCNAME,
            Some(addr),
            Some(cidr),
            interface::Operation::Update,
        )
        .await,

        //TODO: Configure wireguard device
        interface::update_wireguard_ifc(IFCNAME, None, None, interface::Operation::Enable).await,
    ) {
        (Ok(()), Ok(()), Ok(())) => println!("Interface setup completed successfully"),
        (Err(e1), _, _) => eprintln!("Interface creation failed: {}", e1),
        (Ok(()), Err(e2), _) => eprintln!("Interface update failed: {}", e2),
        (Ok(()), _, Err(e3)) => eprintln!("Interface enabling failed: {}", e3),
    }
}
