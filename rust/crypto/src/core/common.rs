use std::fs::File;
use std::io::{BufReader, BufWriter, Read, Write};

pub const MAGIC: &[u8; 4] = b"VPF1";
pub const FORMAT_VERSION: u8 = 1; // FINAL V1: frames include is_last and use encrypt_last/decrypt_last

// alg_id: 1 = XChaCha20-Poly1305 (streaming)
pub const ALG_ID_XCHACHA20POLY1305_STREAM: u8 = 1;
// kdf_id: 1 = HKDF-SHA256
pub const KDF_ID_HKDF_SHA256: u8 = 1;

// StreamBE32 overhead is 5 bytes (4B counter + 1B last-flag).
// XChaCha20Poly1305 nonce is 24 bytes => base nonce = 24 - 5 = 19.
pub const STREAM_BASE_NONCE_LEN: usize = 19;

pub const SALT_LEN: usize = 16; // per-file random salt
pub const DEK_LEN: usize = 32; // derived data encryption key
pub const CHUNK_SIZE: usize = 1024 * 1024; // 1MB streaming chunk

// Header layout (binary):
// 0..4   MAGIC "VPF1"
// 4      format_version (FINAL V1 = 1)
// 5      alg_id
// 6      kdf_id
// 7      reserved
// 8      salt_len (u8)
// 9      nonce_len (u8)  (base nonce length for streaming)
// 10..   salt bytes
// ...    nonce bytes
pub const HEADER_FIXED_LEN: usize = 10;

pub fn write_header(
    w: &mut BufWriter<File>,
    file_salt: &[u8],
    base_nonce: &[u8],
) -> Result<(), String> {
    if file_salt.len() > 255 {
        return Err("salt too long".to_string());
    }
    if base_nonce.len() > 255 {
        return Err("nonce too long".to_string());
    }

    let mut fixed = [0u8; HEADER_FIXED_LEN];
    fixed[0..4].copy_from_slice(MAGIC);
    fixed[4] = FORMAT_VERSION;
    fixed[5] = ALG_ID_XCHACHA20POLY1305_STREAM;
    fixed[6] = KDF_ID_HKDF_SHA256;
    fixed[7] = 0;
    fixed[8] = file_salt.len() as u8;
    fixed[9] = base_nonce.len() as u8;

    w.write_all(&fixed)
        .map_err(|e| format!("write header fixed failed: {e}"))?;
    w.write_all(file_salt)
        .map_err(|e| format!("write salt failed: {e}"))?;
    w.write_all(base_nonce)
        .map_err(|e| format!("write nonce failed: {e}"))?;
    Ok(())
}

pub fn read_header(r: &mut BufReader<File>) -> Result<(Vec<u8>, [u8; 19]), String> {
    let mut fixed = [0u8; HEADER_FIXED_LEN];
    r.read_exact(&mut fixed)
        .map_err(|e| format!("read header fixed failed: {e}"))?;

    if &fixed[0..4] != MAGIC {
        return Err("bad magic".to_string());
    }

    let format_version = fixed[4];
    if format_version != FORMAT_VERSION {
        return Err(format!("unsupported format version: {format_version}"));
    }

    let alg_id = fixed[5];
    let kdf_id = fixed[6];

    if alg_id != ALG_ID_XCHACHA20POLY1305_STREAM {
        return Err(format!("unsupported alg_id: {alg_id}"));
    }
    if kdf_id != KDF_ID_HKDF_SHA256 {
        return Err(format!("unsupported kdf_id: {kdf_id}"));
    }

    let salt_len = fixed[8] as usize;
    let nonce_len = fixed[9] as usize;

    if nonce_len != STREAM_BASE_NONCE_LEN {
        return Err(format!("unexpected base nonce len: {nonce_len}"));
    }
    if salt_len == 0 || salt_len > 64 {
        return Err(format!("unexpected salt len: {salt_len}"));
    }

    let mut salt = vec![0u8; salt_len];
    r.read_exact(&mut salt)
        .map_err(|e| format!("read salt failed: {e}"))?;

    let mut base_nonce = [0u8; STREAM_BASE_NONCE_LEN];
    r.read_exact(&mut base_nonce)
        .map_err(|e| format!("read nonce failed: {e}"))?;

    Ok((salt, base_nonce))
}

pub fn get_thread_count() -> usize {
    // 减去给读取、写入线程留的cpu
    let cpu_count = num_cpus::get().saturating_sub(2);
    // 因为最终要在android上运行，满核运行可能导致降频
    std::cmp::min(4, cpu_count.max(1))
}

pub fn make_nonce(base_nonce: &[u8; 19], is_last_frame: bool, counter: u32) -> [u8; 24] {
    let mut nonce = [0u8; 24];

    nonce[..STREAM_BASE_NONCE_LEN].copy_from_slice(base_nonce);
    nonce[STREAM_BASE_NONCE_LEN] = if is_last_frame { 1 } else { 0 };
    nonce[STREAM_BASE_NONCE_LEN + 1..].copy_from_slice(&counter.to_le_bytes());

    nonce
}

pub struct Frame {
    pub index: u32,
    pub is_last_frame: bool,
    pub data: Vec<u8>,
}
