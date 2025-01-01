use std::fmt;
use std::os::raw::c_char;
use bitflags::bitflags;

//---------------------------------------------- Constants ----------------------------------------------//
// (Define any necessary constants here.)

//---------------------------------------------- Enums ----------------------------------------------//

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

//---------------------------------------------- Interfaces ----------------------------------------------//

#[repr(C)]
#[derive(Debug, PartialEq)]
pub struct WgKey([u8; 32]);

impl WgKey {
    pub fn new(input: [u8; 32]) -> Self {
        WgKey(input)
    }
}

#[repr(C)]
pub struct WgKeyBase64String([u8; 45]);

impl WgKeyBase64String {
    pub fn new(input: [u8; 45]) -> Self {
        WgKeyBase64String(input)
    }
}

impl fmt::Debug for WgKeyBase64String {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let as_str: String = self.0.iter().map(|&c| c as char).collect();
        write!(f, "{}", as_str)
    }
}

#[repr(C)]
#[derive(Debug, Copy, Clone, PartialEq, Eq)]
pub struct Timespec64 {
    pub tv_sec: i64,
    pub tv_nsec: i64,
}

impl Timespec64 {
    pub fn new(tv_sec: i64, tv_nsec: i64) -> Self {
        Timespec64 { tv_sec, tv_nsec }
    }
}

#[repr(C)]
#[derive(Copy, Clone)]
pub union WgEndpoint {
    pub addr: Sockaddr,
    pub addr4: SockaddrIn,
    pub addr6: SockaddrIn6,
}

impl WgEndpoint {
    pub fn new_addr(addr: Sockaddr) -> Self {
        WgEndpoint { addr }
    }

    pub fn new_addr4(addr4: SockaddrIn) -> Self {
        WgEndpoint { addr4 }
    }

    pub fn new_addr6(addr6: SockaddrIn6) -> Self {
        WgEndpoint { addr6 }
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct Sockaddr {
    pub sa_family: u16,
    pub sa_data: [u8; 14],
}

impl Sockaddr {
    pub fn new(sa_family: u16, sa_data: [u8; 14]) -> Self {
        Sockaddr { sa_family, sa_data }
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct SockaddrIn {
    pub sin_family: u16,
    pub sin_port: u16,
    pub sin_addr: std::net::Ipv4Addr,
    pub sin_zero: [u8; 8],
}

impl SockaddrIn {
    pub fn new(
        sin_family: u16,
        sin_port: u16,
        sin_addr: std::net::Ipv4Addr,
        sin_zero: [u8; 8],
    ) -> Self {
        SockaddrIn {
            sin_family,
            sin_port,
            sin_addr,
            sin_zero,
        }
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct SockaddrIn6 {
    pub sin6_family: u16,
    pub sin6_port: u16,
    pub sin6_flowinfo: u32,
    pub sin6_addr: std::net::Ipv6Addr,
    pub sin6_scope_id: u32,
}

impl SockaddrIn6 {
    pub fn new(
        sin6_family: u16,
        sin6_port: u16,
        sin6_flowinfo: u32,
        sin6_addr: std::net::Ipv6Addr,
        sin6_scope_id: u32,
    ) -> Self {
        SockaddrIn6 {
            sin6_family,
            sin6_port,
            sin6_flowinfo,
            sin6_addr,
            sin6_scope_id,
        }
    }
}

#[repr(C)]
#[derive(Clone, Copy)]
pub struct WgAllowedIp {
    pub family: u16,
    pub cidr: u8,
    pub ip: Ip,
    pub next_allowedip: *mut WgAllowedIp,
}

impl WgAllowedIp {
    pub fn new(family: u16, cidr: u8, ip: Ip, next_allowedip: *mut WgAllowedIp) -> Self {
        WgAllowedIp {
            family,
            cidr,
            ip,
            next_allowedip,
        }
    }
}

#[repr(C)]
#[derive(Clone, Copy)]
pub union Ip {
    pub ip4: InAddr,
    pub ip6: In6Addr,
}

impl Ip {
    pub fn new_ip4(ip4: InAddr) -> Self {
        Ip { ip4 }
    }

    pub fn new_ip6(ip6: In6Addr) -> Self {
        Ip { ip6 }
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct InAddr {
    pub s_addr: u32,
}

impl InAddr {
    pub fn new(s_addr: u32) -> Self {
        InAddr { s_addr }
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct In6Addr {
    pub s6_addr: [u8; 16],
}

impl In6Addr {
    pub fn new(s6_addr: [u8; 16]) -> Self {
        In6Addr { s6_addr }
    }
}

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
    pub first_allowedip: *mut WgAllowedIp,
    pub last_allowedip: *mut WgAllowedIp,
    pub next_peer: *mut WgPeer,
}

impl WgPeer {
    pub fn new(
        flags: WgPeerFlags,
        public_key: WgKey,
        preshared_key: WgKey,
        endpoint: WgEndpoint,
        last_handshake_time: Timespec64,
        rx_bytes: u64,
        tx_bytes: u64,
        persistent_keepalive_interval: u16,
        first_allowedip: *mut WgAllowedIp,
        last_allowedip: *mut WgAllowedIp,
        next_peer: *mut WgPeer,
    ) -> Self {
        WgPeer {
            flags,
            public_key,
            preshared_key,
            endpoint,
            last_handshake_time,
            rx_bytes,
            tx_bytes,
            persistent_keepalive_interval,
            first_allowedip,
            last_allowedip,
            next_peer,
        }
    }
}

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

impl WgDevice {
    pub fn new(
        name: [u8; 16],
        ifindex: u32,
        flags: WgDeviceFlags,
        public_key: WgKey,
        private_key: WgKey,
        fwmark: u32,
        listen_port: u16,
        first_peer: *mut WgPeer,
        last_peer: *mut WgPeer,
    ) -> Self {
        WgDevice {
            name,
            ifindex,
            flags,
            public_key,
            private_key,
            fwmark,
            listen_port,
            first_peer,
            last_peer,
        }
    }
}

//---------------------------------------------- Functions ----------------------------------------------//

extern "C" {
    pub fn wg_generate_private_key(private_key: *mut WgKey);
    pub fn wg_key_to_base64(
        wg_key_b64_string: *mut WgKeyBase64String,
        wg_key_int: *const WgKey,
    );
    pub fn wg_key_from_base64(
        wg_key_int: *mut WgKey,
        wg_key_b64_string: *const WgKeyBase64String,
    );
    pub fn wg_generate_public_key(public_key: *mut WgKey, private_key: *const WgKey);
    pub fn wg_list_device_names() -> *const c_char;
}
