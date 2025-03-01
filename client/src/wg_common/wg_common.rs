use super::wireguard_cffi::{wg_get_device, wg_set_device, WgDevice};
use crate::wg_common::wireguard_cffi::{
    wg_generate_private_key, wg_generate_public_key, wg_key_from_base64, wg_key_to_base64,
    wg_list_device_names, WgAllowedIp, WgDeviceFlags, WgEndpoint, WgKey, WgKeyBase64String, WgPeer,
};
use std::ffi::{CStr, CString};

pub fn gen_private_key() -> WgKeyBase64String {
    let mut private_key_int = WgKey([0; 32]);
    let mut private_key = WgKeyBase64String([b' '; 45]);
    unsafe {
        wg_generate_private_key(&mut private_key_int);
        wg_key_to_base64(&mut private_key, &private_key_int);
    }
    private_key
}

pub fn gen_public_key(private_key: &WgKeyBase64String) -> WgKeyBase64String {
    let mut private_key_int = WgKey([0; 32]);
    let mut public_key_int = WgKey([0; 32]);
    let mut public_key = WgKeyBase64String([b' '; 45]);
    unsafe {
        wg_key_from_base64(&mut private_key_int, private_key);
        wg_generate_public_key(&mut public_key_int, &private_key_int);
        wg_key_to_base64(&mut public_key, &public_key_int);
    }
    public_key
}
pub fn wg_key_from_str(input: &str) -> WgKey {
    let key = WgKeyBase64String::from(input);
    let mut key_int = WgKey([0; 32]);

    unsafe {
        wg_key_from_base64(&mut key_int, &key);
    }

    key_int
}

pub fn list_device_names() -> Vec<String> {
    let mut result = Vec::new();
    unsafe {
        let ptr = wg_list_device_names();
        if ptr.is_null() {
            return result;
        }
        let mut offset = 0;
        loop {
            let c_str = CStr::from_ptr(ptr.add(offset));
            let slice = c_str.to_bytes_with_nul();
            if slice.is_empty() || slice == [0] {
                break;
            }
            result.push(
                String::from_utf8_lossy(slice)
                    .trim_end_matches('\0')
                    .to_string(),
            );
            offset += slice.len();
        }
    }
    result
}

pub fn update_device(
    priv_key: &WgKeyBase64String,
    listen_port: u16,
    device_name: &str,
    server_pub: &str,
    s_endpoint: &str,
    s_ip: &str,
) -> Result<(), i32> {
    let c_device_name = CString::new(device_name).map_err(|_| -1)?;
    let mut device: *mut WgDevice = std::ptr::null_mut();
    let mut private_key_int = WgKey([0; 32]);

    let server_pub_key_int = wg_key_from_str(server_pub);
    let server_endpoint = WgEndpoint::from(s_endpoint);
    let mut server_allowed_ip = WgAllowedIp::from(s_ip);
    let mut server_peer = WgPeer::init(server_pub_key_int, server_endpoint, &mut server_allowed_ip);

    unsafe {
        if wg_get_device(&mut device, c_device_name.as_ptr()) != 0 {
            let err = wg_get_device(&mut device, c_device_name.as_ptr());
            println!("Error {:?} getting device by name {}", err, device_name);
            return Err(err);
        }

        let mut temp_dev = device.read();
        wg_key_from_base64(&mut private_key_int, priv_key);
        temp_dev.private_key = private_key_int;
        temp_dev.listen_port = listen_port;
        temp_dev
            .flags
            .insert(WgDeviceFlags::HAS_PRIVATE_KEY | WgDeviceFlags::HAS_LISTEN_PORT);
        temp_dev.first_peer = &mut server_peer;

        if wg_set_device(&mut temp_dev) != 0 {
            let err = wg_set_device(&mut temp_dev);
            println!("Failed to set device: {}", err);
            return Err(err);
        }
    }

    println!("Device updated successfully.");
    Ok(())
}
