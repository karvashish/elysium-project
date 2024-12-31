use std::fmt;

#[repr(C)]
#[derive(Debug, PartialEq)]
pub struct WgKey([u8; 32]);

#[repr(C)]
pub struct WgKeyBase64String([u8; 45]);

impl WgKey {
    pub fn new(input: [u8; 32]) -> WgKey {
        WgKey(input)
    }
}

impl WgKeyBase64String {
    pub fn new(input: [u8; 45]) -> WgKeyBase64String {
        WgKeyBase64String(input)
    }
}

impl fmt::Debug for WgKeyBase64String {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let as_str: String = self.0.iter().map(|&c| c as u8 as char).collect();
        if f.alternate() {
            write!(f, "WgKeyBase64String(\"{}\")", as_str)
        } else {
            write!(f, "{}", as_str)
        }
    }
}

unsafe extern "C" {
    pub unsafe fn wg_generate_private_key(private_key: *mut WgKey);
    pub unsafe fn wg_key_to_base64(
        wg_key_b64_string: *mut WgKeyBase64String,
        wg_key_int: *const WgKey,
    );
    pub unsafe fn wg_key_from_base64(
        wg_key_int: *mut WgKey,
        wg_key_b64_string: *const WgKeyBase64String,
    );
    pub unsafe fn wg_generate_public_key(public_key: *mut WgKey, private_key: *const WgKey);
}
