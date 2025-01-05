use bitflags::bitflags;
use std::{fmt, net::{Ipv4Addr, Ipv6Addr}, os::raw::c_char};

//---------------------------------------------- Constants ----------------------------------------------//

const AF_INET: u16 = 2;
const AF_INET6: u16 = 10;

//---------------------------------------------- Enums and Flags ----------------------------------------------//

bitflags! {
    #[repr(transparent)]
    #[derive( Clone, Copy, PartialEq, Eq)]
    pub struct WgDeviceFlags: u32 {
        const REPLACE_PEERS = 1 << 0;
        const HAS_PRIVATE_KEY = 1 << 1;
        const HAS_PUBLIC_KEY = 1 << 2;
        const HAS_LISTEN_PORT = 1 << 3;
        const HAS_FWMARK = 1 << 4;
    }
}

impl fmt::Debug for WgDeviceFlags {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let mut active_flags = vec![];

        if self.contains(WgDeviceFlags::REPLACE_PEERS) {
            active_flags.push("REPLACE_PEERS");
        }
        if self.contains(WgDeviceFlags::HAS_PRIVATE_KEY) {
            active_flags.push("HAS_PRIVATE_KEY");
        }
        if self.contains(WgDeviceFlags::HAS_PUBLIC_KEY) {
            active_flags.push("HAS_PUBLIC_KEY");
        }
        if self.contains(WgDeviceFlags::HAS_LISTEN_PORT) {
            active_flags.push("HAS_LISTEN_PORT");
        }
        if self.contains(WgDeviceFlags::HAS_FWMARK) {
            active_flags.push("HAS_FWMARK");
        }

        write!(f, "WgDeviceFlags({})\n", active_flags.join(" | "))
    }
}

bitflags! {
    #[repr(transparent)]
    #[derive( Clone, Copy, PartialEq, Eq)]
    pub struct WgPeerFlags: u32 {
        const REMOVE_ME = 1 << 0;
        const REPLACE_ALLOWEDIPS = 1 << 1;
        const HAS_PUBLIC_KEY = 1 << 2;
        const HAS_PRESHARED_KEY = 1 << 3;
        const HAS_PERSISTENT_KEEPALIVE_INTERVAL = 1 << 4;
    }
}

impl fmt::Debug for WgPeerFlags {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let mut active_flags = vec![];

        if self.contains(WgPeerFlags::REMOVE_ME) {
            active_flags.push("REMOVE_ME");
        }
        if self.contains(WgPeerFlags::REPLACE_ALLOWEDIPS) {
            active_flags.push("REPLACE_ALLOWEDIPS");
        }
        if self.contains(WgPeerFlags::HAS_PUBLIC_KEY) {
            active_flags.push("HAS_PUBLIC_KEY");
        }
        if self.contains(WgPeerFlags::HAS_PRESHARED_KEY) {
            active_flags.push("HAS_PRESHARED_KEY");
        }
        if self.contains(WgPeerFlags::HAS_PERSISTENT_KEEPALIVE_INTERVAL) {
            active_flags.push("HAS_PERSISTENT_KEEPALIVE_INTERVAL");
        }

        write!(f, "WgPeerFlags({})\n", active_flags.join(" | "))
    }
}

//---------------------------------------------- Helper Structs ----------------------------------------------//

/// A wrapper for an IPv4 address compatible with FFI.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct FfiIpv4Addr {
    pub octets: [u8; 4],
}

impl From<std::net::Ipv4Addr> for FfiIpv4Addr {
    fn from(ip: std::net::Ipv4Addr) -> Self {
        FfiIpv4Addr {
            octets: ip.octets(),
        }
    }
}

impl From<FfiIpv4Addr> for std::net::Ipv4Addr {
    fn from(ip: FfiIpv4Addr) -> Self {
        std::net::Ipv4Addr::from(ip.octets)
    }
}

/// A wrapper for an IPv6 address compatible with FFI.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct FfiIpv6Addr {
    pub segments: [u16; 8],
}

impl From<std::net::Ipv6Addr> for FfiIpv6Addr {
    fn from(ip: std::net::Ipv6Addr) -> Self {
        FfiIpv6Addr {
            segments: ip.segments(),
        }
    }
}

impl From<FfiIpv6Addr> for std::net::Ipv6Addr {
    fn from(ip: FfiIpv6Addr) -> Self {
        std::net::Ipv6Addr::from(ip.segments)
    }
}

/// A wrapper for a WireGuard key.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct WgKey(pub [u8; 32]);

/// A base64-encoded WireGuard key.
#[repr(C)]
#[derive(Clone, Copy, PartialEq, Eq)]
pub struct WgKeyBase64String(pub [u8; 45]);

impl From<&str> for WgKeyBase64String {
    fn from(input: &str) -> Self {
        let mut array = [0u8; 45];
        let bytes = input.as_bytes();

        let len = bytes.len().min(45);
        array[..len].copy_from_slice(&bytes[..len]);

        WgKeyBase64String(array)
    }
}

impl fmt::Debug for WgKeyBase64String {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let as_str: String = self.0.iter().map(|&c| c as char).collect();
        write!(f, "{}", as_str)
    }
}

//---------------------------------------------- Core Structs ----------------------------------------------//

/// Represents a WireGuard peer.
#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct WgPeer {
    pub flags: WgPeerFlags,
    pub public_key: WgKey,
    pub preshared_key: WgKey,
    pub endpoint: WgEndpoint,
    pub last_handshake_time: Timespec64,
    pub rx_bytes: u64,
    pub tx_bytes: u64,
    pub persistent_keepalive_interval: u16,
    pub first_allowed_ip: *mut WgAllowedIp,
    pub last_allowed_ip: *mut WgAllowedIp,
    pub next_peer: *mut WgPeer,
}


impl WgPeer {
    pub fn init(public_key: WgKey, endpoint: WgEndpoint, first_allowed_ip: *mut WgAllowedIp) -> Self {
        WgPeer {
            flags: WgPeerFlags::HAS_PUBLIC_KEY,
            public_key,
            preshared_key: WgKey([0; 32]),
            endpoint,
            last_handshake_time: Timespec64 {
                tv_sec: 0,
                tv_nsec: 0,
            },
            rx_bytes: 0,
            tx_bytes: 0,
            persistent_keepalive_interval: 0,
            first_allowed_ip,
            last_allowed_ip: first_allowed_ip,
            next_peer: std::ptr::null_mut(),
        }
    }
}


/// Represents a WireGuard device.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct WgDevice {
    pub name: [u8; 16],
    pub ifindex: u32,
    pub flags: WgDeviceFlags,
    pub public_key: WgKey,
    pub private_key: WgKey,
    pub fwmark: u32,
    pub listen_port: u16,
    pub first_peer: *mut WgPeer,
    pub last_peer: *mut WgPeer,
}

//---------------------------------------------- Miscellaneous Structs ----------------------------------------------//

/// Represents a generic socket address.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Sockaddr {
    pub sa_family: u16,
    pub sa_data: [u8; 14],
}

/// Represents a socket address for IPv4.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct SockaddrIn {
    pub sin_family: u16,
    pub sin_port: u16,
    pub sin_addr: FfiIpv4Addr,
    pub sin_zero: [u8; 8],
}

/// Represents a socket address for IPv6.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct SockaddrIn6 {
    pub sin6_family: u16,
    pub sin6_port: u16,
    pub sin6_flowinfo: u32,
    pub sin6_addr: FfiIpv6Addr,
    pub sin6_scope_id: u32,
}

/// A union representing a WireGuard endpoint.
#[repr(C)]
#[derive(Copy, Clone)]
pub union WgEndpoint {
    pub addr: Sockaddr,
    pub addr4: SockaddrIn,
    pub addr6: SockaddrIn6,
}

impl fmt::Debug for WgEndpoint {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        unsafe {
            f.debug_struct("WgEndpoint")
                .field("addr", &self.addr)
                .field("addr4", &self.addr4)
                .field("addr6", &self.addr6)
                .finish()
        }
    }
}

impl From<&str> for WgEndpoint {
    fn from(input: &str) -> Self {
        if let Ok(std::net::SocketAddr::V4(addr)) = input.parse() {
            let ip = addr.ip().octets();
            let port = addr.port();

            let sockaddr_in = SockaddrIn {
                sin_family: AF_INET,
                sin_port: port.to_be(),
                sin_addr: FfiIpv4Addr { octets: ip },
                sin_zero: [0u8; 8],
            };

            WgEndpoint { addr4: sockaddr_in }
        } else if let Ok(std::net::SocketAddr::V6(addr)) = input.parse() {
            let ip = addr.ip().segments();
            let port = addr.port();

            let sockaddr_in6 = SockaddrIn6 {
                sin6_family: AF_INET6,
                sin6_port: port.to_be(),
                sin6_flowinfo: 0,
                sin6_addr: FfiIpv6Addr { segments: ip },
                sin6_scope_id: 0,
            };

            WgEndpoint { addr6: sockaddr_in6 }
        } else {
            panic!("Invalid socket address");
        }
    }
}


/// Represents an allowed IP address for a WireGuard peer.
#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct WgAllowedIp {
    pub family: u16,
    pub cidr: u8,
    pub ip: Ip,
    pub next_allowed_ip: *mut WgAllowedIp,
}

impl From<&str> for WgAllowedIp {
    fn from(input: &str) -> Self {
        // Split the input string into IP and CIDR parts
        let parts: Vec<&str> = input.split('/').collect();
        if parts.len() != 2 {
            panic!("Invalid IP/CIDR format");
        }

        let ip_str = parts[0];
        let cidr: u8 = parts[1].parse().expect("Invalid CIDR value");

        if let Ok(ipv4) = ip_str.parse::<Ipv4Addr>() {
            WgAllowedIp {
                family: AF_INET,
                cidr,
                ip: Ip {
                    ip4: FfiIpv4Addr {
                        octets: ipv4.octets(),
                    },
                },
                next_allowed_ip: std::ptr::null_mut(),
            }
        } else if let Ok(ipv6) = ip_str.parse::<Ipv6Addr>() {
            WgAllowedIp {
                family: AF_INET6,
                cidr,
                ip: Ip {
                    ip6: FfiIpv6Addr {
                        segments: ipv6.segments(),
                    },
                },
                next_allowed_ip: std::ptr::null_mut(),
            }
        } else {
            panic!("Invalid IP address format");
        }
    }
}

/// A union for IPv4 or IPv6 addresses.
#[repr(C)]
#[derive(Clone, Copy)]
pub union Ip {
    pub ip4: FfiIpv4Addr,
    pub ip6: FfiIpv6Addr,
}

impl fmt::Debug for Ip {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        unsafe {
            match self {
                Ip { ip4 } => f.debug_struct("Ip").field("ip4", &ip4).finish(),
            }
        }
    }
}

/// Represents a 64-bit timestamp.
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Timespec64 {
    pub tv_sec: i64,
    pub tv_nsec: i64,
}

//---------------------------------------------- FFI Functions ----------------------------------------------//

extern "C" {
    pub fn wg_generate_private_key(private_key: *mut WgKey);
    pub fn wg_key_to_base64(wg_key_b64_string: *mut WgKeyBase64String, wg_key_int: *const WgKey);
    pub fn wg_key_from_base64(wg_key_int: *mut WgKey, wg_key_b64_string: *const WgKeyBase64String);
    pub fn wg_generate_public_key(public_key: *mut WgKey, private_key: *const WgKey);
    pub fn wg_list_device_names() -> *const c_char;
    pub fn wg_get_device(dev: *mut *mut WgDevice, device_name: *const c_char) -> i32;
    pub fn wg_set_device(dev: *mut WgDevice) -> i32;
}
