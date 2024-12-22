fn main() { 
    println!("Elysium Project Client Running");
    if let Some(pub_key) = option_env!("PUBKEY") {
        println!("PUBKEY: {}", pub_key);
    } else {
        println!("PUBKEY is not set");
    }
} 
