mod common;
mod core_test;
mod decrypt;
pub mod encrypt_v1;
pub mod encrypt_v2;
mod kdf;
mod buffer_pool;

pub use decrypt::decrypt_file_impl;
pub use encrypt_v1::encrypt_file_impl;

// for benchmark test
pub trait Encrypt {
    fn encrypt(
        in_path: &str,
        out_path: &str,
        local_material: &[u8],
        cloud_material: &[u8],
        file_id: &[u8],
        aad: &[u8],
    ) -> Result<(), String>;
}
