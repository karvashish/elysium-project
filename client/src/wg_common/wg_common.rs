use crate::wg_common::wireguard_cffi::{
    wg_generate_private_key, wg_generate_public_key, wg_key_from_base64, wg_key_to_base64, wg_list_device_names, WgKey,
    WgKeyBase64String,
};

use std::ffi::{CStr, CString};

use super::wireguard_cffi::{wg_get_device, WgDevice};

pub fn gen_private_key() -> WgKeyBase64String {
    let mut private_key_int= WgKey([0; 32]);
    let mut private_key = WgKeyBase64String([b' ' as u8; 45]);

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
    let mut private_key_int = WgKey([0; 32]);
    let mut public_key_int = WgKey([0; 32]);
    let mut public_key = WgKeyBase64String([b' ' as u8; 45]);

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

pub fn list_device_names() -> Vec<String> {
    let mut result = Vec::new();
    let mut offset = 0;

    unsafe {
        let ptr = wg_list_device_names();
        if ptr.is_null() {
            return Vec::new();
        }

        loop {
            let c_str = CStr::from_ptr(ptr.add(offset));
            let slice = c_str.to_bytes_with_nul();
            if slice.is_empty() || slice == [0] {
                break;
            }

            result.push(String::from_utf8_lossy(slice).trim_end_matches('\0').to_string());
            offset += slice.len();
        }
    }
    return result;
}

pub fn get_device(device_name: &str, device: &mut *mut WgDevice) -> Result<(), i32> {
    let c_device_name = CString::new(device_name)
        .map_err(|_| -1)?;

    unsafe {
        let result = wg_get_device(device, c_device_name.as_ptr());
        if result == 0 {
            Ok(())
        } else {
            Err(result)
        }
    }
}
