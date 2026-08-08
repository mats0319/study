use crate::core::{decrypt_file_impl, encrypt_file_impl};
use jni::JNIEnv;
use jni::objects::{JByteArray, JClass, JString};
use jni::sys::jint;

/// encrypt_file (FINAL V1):
/// - Writes header (FORMAT_VERSION=1)
/// - Frames: [1-byte is_last][4-byte BE ct_len][ct bytes]
/// - Uses encrypt_last on final frame
#[unsafe(no_mangle)]
pub extern "system" fn Java_com_example_demo_Crypto_encryptFile(
    mut env: JNIEnv,
    _class: JClass,

    in_path: JString,
    out_path: JString,

    local_key: JByteArray,
    cloud_key: JByteArray,

    file_id: JByteArray,
    aad: JByteArray,
) -> jint {
    let in_path: String = env.get_string(&in_path).unwrap().into();
    let out_path: String = env.get_string(&out_path).unwrap().into();
    let local_key: Vec<u8> = env.convert_byte_array(local_key).unwrap();
    let cloud_key: Vec<u8> = env.convert_byte_array(cloud_key).unwrap();
    let file_id: Vec<u8> = env.convert_byte_array(file_id).unwrap();
    let aad: Vec<u8> = env.convert_byte_array(aad).unwrap();

    match encrypt_file_impl(&in_path, &out_path, &local_key, &cloud_key, &file_id, &aad) {
        Ok(()) => 0,
        Err(err) => {
            let _ = env.throw_new("java/lang/RuntimeException", err);
            -1
        }
    }
}

/// decrypt_file (FINAL V1):
/// - Reads header (FORMAT_VERSION=1)
/// - Frames: [1-byte is_last][4-byte BE ct_len][ct bytes]
/// - Uses decrypt_last on final frame, and rejects missing-last / trailing-data
#[unsafe(no_mangle)]
pub extern "system" fn Java_com_example_demo_Crypto_decryptFile(
    mut env: JNIEnv,
    _class: JClass,

    in_path: JString,
    out_path: JString,

    local_key: JByteArray,
    cloud_key: JByteArray,

    file_id: JByteArray,
    aad: JByteArray,
) -> jint {
    let in_path: String = env.get_string(&in_path).unwrap().into();
    let out_path: String = env.get_string(&out_path).unwrap().into();
    let local_key: Vec<u8> = env.convert_byte_array(local_key).unwrap();
    let cloud_key: Vec<u8> = env.convert_byte_array(cloud_key).unwrap();
    let file_id: Vec<u8> = env.convert_byte_array(file_id).unwrap();
    let aad: Vec<u8> = env.convert_byte_array(aad).unwrap();

    match decrypt_file_impl(&in_path, &out_path, &local_key, &cloud_key, &file_id, &aad) {
        Ok(()) => 0,
        Err(err) => {
            let _ = env.throw_new("java/lang/RuntimeException", err);
            -1
        }
    }
}
