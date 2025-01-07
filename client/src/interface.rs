use futures_util::{TryFutureExt, TryStreamExt};
use rtnetlink::new_connection;
use std::net::Ipv4Addr;

pub enum Operation {
    Enable,
    Update,
}

pub async fn create_wireguard_ifc(name: &str) -> Result<(), String> {
    let (connection, handle, _) =
        new_connection().map_err(|e| format!("Connection setup failed: {e}"))?;
    tokio::spawn(connection);

    handle
        .link()
        .add()
        .wireguard(name.to_string())
        .execute()
        .await
        .map_err(|e| {
            if e.to_string().contains("File exists") {
                println!("Interface {} already exists", name);
                String::from("Interface already exists")
            } else {
                eprintln!("Error creating interface {}: {}", name, e);
                e.to_string()
            }
        })
        .map(|_| println!("Successfully created interface {}...", name))
}

pub async fn update_wireguard_ifc(
    name: &str,
    addr: Option<Ipv4Addr>,
    cidr: Option<u8>,
    op: Operation,
) -> Result<(), String> {
    let (connection, handle, _) =
        new_connection().map_err(|e| format!("Connection setup failed: {e}"))?;
    tokio::spawn(connection);

    let mut links = handle.link().get().match_name(name.to_string()).execute();
    if let Some(link) = links
        .try_next()
        .map_err(|e| format!("Error retrieving link {name}: {e}"))
        .await?
    {
        match op {
            Operation::Update => {
                if let (Some(addr), Some(cidr)) = (addr, cidr) {
                    handle
                        .address()
                        .add(link.header.index, std::net::IpAddr::V4(addr), cidr)
                        .execute()
                        .map_err(|e| format!("Error adding address to {name}: {e}"))
                        .await?
                } else {
                    return Err("Address and CIDR must be provided for Operation::Update".into());
                }
            }

            Operation::Enable => {
                handle
                    .link()
                    .set(link.header.index)
                    .up()
                    .execute()
                    .map_err(|e| format!("Error bringing interface {name} up: {e}"))
                    .await?
            }
        }

        println!("Interface {} updated and enabled successfully", name);
    } else {
        return Err(format!("No link named {name} found"));
    }
    Ok(())
}
