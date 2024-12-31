use crate::wg_common::wireguard_cffi::{
    wg_generate_private_key, wg_generate_public_key, wg_key_from_base64, wg_key_to_base64, WgKey,
    WgKeyBase64String,
};

pub fn gen_private_key() -> WgKeyBase64String {
    let mut private_key_int: WgKey = WgKey::new([0; 32]);
    let mut private_key: WgKeyBase64String = WgKeyBase64String::new([b' ' as u8; 45]);

    unsafe {
        wg_generate_private_key(&mut private_key_int as *mut WgKey);
        wg_key_to_base64(
            &mut private_key as *mut WgKeyBase64String,
            &mut private_key_int as *const WgKey,
        );
    }
    return private_key;
}

pub fn gen_public_key(private_key: &WgKeyBase64String) -> WgKeyBase64String {
    let mut private_key_int: WgKey = WgKey::new([0; 32]);
    let mut public_key_int: WgKey = WgKey::new([0; 32]);
    let mut public_key: WgKeyBase64String = WgKeyBase64String::new([b' ' as u8; 45]);

    unsafe {
        wg_key_from_base64(&mut private_key_int as *mut WgKey, private_key);
        wg_generate_public_key(
            &mut public_key_int as *mut WgKey,
            &mut private_key_int as *const WgKey,
        );
        wg_key_to_base64(
            &mut public_key as *mut WgKeyBase64String,
            &mut public_key_int as *const WgKey,
        );
    }
    return public_key;
}
