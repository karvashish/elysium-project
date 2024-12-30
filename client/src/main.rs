use std::net::Ipv4Addr;

mod interface;

#[tokio::main]
async fn main() {
    println!("Starting Elysium Project Client setup");

    let pub_key = option_env!("PUBKEY");
    if let Some(key) = pub_key {
        println!("PUBKEY: {}", key);
    } else {
        eprintln!("PUBKEY is not set");
    }

    let ifc_name = option_env!("IFCNAME").unwrap_or("wg0");

    let addr = option_env!("ADDR").and_then(|addr| addr.parse::<Ipv4Addr>().ok());
    let cidr = option_env!("CIDR").and_then(|cidr| cidr.parse::<u8>().ok());

    let (addr, cidr) = match (addr, cidr) {
        (Some(a), Some(c)) => (a, c),
        _ => {
            eprintln!("Error: Both ADDR and CIDR must be set");
            std::process::exit(1);
        }
    };

    match (
        interface::create_wireguard_ifc(ifc_name).await,
        interface::update_and_enable_ifc(ifc_name, addr, cidr).await,
    ) {
        (Ok(()), Ok(())) => println!("Interface setup completed successfully"),
        (Err(e1), Ok(())) => eprintln!("Interface creation failed: {}", e1),
        (Ok(()), Err(e2)) => eprintln!("Interface enabling failed: {}", e2),
        (Err(e1), Err(e2)) => eprintln!("Both creation and enabling failed: {}, {}", e1, e2),
    }
}
