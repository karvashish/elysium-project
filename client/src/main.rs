use std::net::Ipv4Addr;

mod interface;
mod wg_common;

use wg_common::{
    wg_common::{gen_private_key, gen_public_key},
    wireguard_cffi::WgKeyBase64String,
};

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
        interface::update_wireguard_ifc(IFCNAME, None, None, interface::Operation::Enable).await,
    ) {
        (Ok(()), Ok(()), Ok(())) => println!("Interface setup completed successfully"),
        (Err(e1), _, _) => eprintln!("Interface creation failed: {}", e1),
        (Ok(()), Err(e2), _) => eprintln!("Interface update failed: {}", e2),
        (Ok(()), _, Err(e3)) => eprintln!("Interface enabling failed: {}", e3),
    }
}
