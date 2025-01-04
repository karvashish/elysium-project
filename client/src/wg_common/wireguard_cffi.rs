use std::{fmt, os::raw::c_char};
use bitflags::bitflags;

//---------------------------------------------- Constants ----------------------------------------------//

// Define any necessary constants here.

//---------------------------------------------- Enums and Flags ----------------------------------------------//

bitflags! {
    #[repr(transparent)]
    pub struct WgDeviceFlags: u32 {
        const REPLACE_PEERS = 1 << 0;
        const HAS_PRIVATE_KEY = 1 << 1;
        const HAS_PUBLIC_KEY = 1 << 2;
        const HAS_LISTEN_PORT = 1 << 3;
        const HAS_FWMARK = 1 << 4;
    }
}

bitflags! {
    #[repr(transparent)]
    pub struct WgPeerFlags: u32 {
        const REMOVE_ME = 1 << 0;
        const REPLACE_ALLOWEDIPS = 1 << 1;
        const HAS_PUBLIC_KEY = 1 << 2;
        const HAS_PRESHARED_KEY = 1 << 3;
        const HAS_PERSISTENT_KEEPALIVE_INTERVAL = 1 << 4;
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
        FfiIpv4Addr { octets: ip.octets() }
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
        FfiIpv6Addr { segments: ip.segments() }
    }
}

impl From<FfiIpv6Addr> for std::net::Ipv6Addr {
    fn from(ip: FfiIpv6Addr) -> Self {
        std::net::Ipv6Addr::from(ip.segments)
    }
}

/// A wrapper for a WireGuard key.
#[repr(C)]
#[derive(Debug, PartialEq, Eq)]
pub struct WgKey(pub [u8; 32]);

/// A base64-encoded WireGuard key.
#[repr(C)]
#[derive(PartialEq, Eq)]
pub struct WgKeyBase64String(pub [u8; 45]);

impl fmt::Debug for WgKeyBase64String {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let as_str: String = self.0.iter().map(|&c| c as char).collect();
        write!(f, "{}", as_str)
    }
}


//---------------------------------------------- Core Structs ----------------------------------------------//

/// Represents a WireGuard peer.
#[repr(C)]
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

/// Represents a WireGuard device.
#[repr(C)]
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
#[derive(Debug, Clone, Copy)]
pub struct SockaddrIn {
    pub sin_family: u16,
    pub sin_port: u16,
    pub sin_addr: FfiIpv4Addr,
    pub sin_zero: [u8; 8],
}

/// Represents a socket address for IPv6.
#[repr(C)]
#[derive(Debug, Clone, Copy)]
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

/// Represents an allowed IP address for a WireGuard peer.
#[repr(C)]
pub struct WgAllowedIp {
    pub family: u16,
    pub cidr: u8,
    pub ip: Ip,
    pub next_allowed_ip: *mut WgAllowedIp,
}

/// A union for IPv4 or IPv6 addresses.
#[repr(C)]
#[derive(Clone, Copy)]
pub union Ip {
    pub ip4: FfiIpv4Addr,
    pub ip6: FfiIpv6Addr,
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
