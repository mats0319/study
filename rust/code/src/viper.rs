// src/lib.rs
use std::cell::RefCell;
use std::ffi::{CStr, CString};
use std::fs::File;
use std::io::{BufReader, BufWriter, Read, Write};
use std::os::raw::{c_char, c_int};
use std::slice;

use aead::generic_array::GenericArray;
use aead::stream::{DecryptorBE32, EncryptorBE32};
use aead::{KeyInit, Payload};
use chacha20poly1305::XChaCha20Poly1305;

use hkdf::Hkdf;
use rand::RngCore;
use sha2::Sha256;

const MAGIC: &[u8; 4] = b"VPF1";
const FORMAT_VERSION: u8 = 1; // FINAL V1: frames include is_last and use encrypt_last/decrypt_last

// alg_id: 1 = XChaCha20-Poly1305 (streaming)
const ALG_ID_XCHACHA20POLY1305_STREAM: u8 = 1;
// kdf_id: 1 = HKDF-SHA256
const KDF_ID_HKDF_SHA256: u8 = 1;

// StreamBE32 overhead is 5 bytes (4B counter + 1B last-flag).
// XChaCha20Poly1305 nonce is 24 bytes => base nonce = 24 - 5 = 19.
const STREAM_BASE_NONCE_LEN: usize = 19;

const SALT_LEN: usize = 16; // per-file random salt
const DEK_LEN: usize = 32; // derived data encryption key
const CHUNK_SIZE: usize = 1024 * 1024; // 1MB streaming chunk

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
const HEADER_FIXED_LEN: usize = 10;

thread_local! {
    static LAST_ERROR: RefCell<CString> =
        RefCell::new(CString::new("OK").unwrap());
}

fn set_last_error(msg: impl AsRef<str>) {
    let s = msg.as_ref();
    let c = CString::new(s).unwrap_or_else(|_| CString::new("Invalid error string").unwrap());
    LAST_ERROR.with(|e| *e.borrow_mut() = c);
}

#[unsafe(no_mangle)]
pub extern "C" fn get_last_error_message() -> *const c_char {
    LAST_ERROR.with(|e| e.borrow().as_ptr())
}

// ============ FFI helpers ============

fn cstr_to_string(ptr: *const c_char, name: &str) -> Result<String, String> {
    if ptr.is_null() {
        return Err(format!("{name} is null"));
    }
    let s = unsafe { CStr::from_ptr(ptr) }
        .to_str()
        .map_err(|_| format!("{name} is not valid UTF-8"))?;
    if s.is_empty() {
        return Err(format!("{name} is empty"));
    }
    Ok(s.to_string())
}

fn bytes_from_ptr<'a>(ptr: *const u8, len: c_int, name: &str) -> Result<&'a [u8], String> {
    if len < 0 {
        return Err(format!("{name} len < 0"));
    }

    let len = len as usize;

    if len == 0 {
        // SAFETY: zero-length slice from dangling pointer is allowed
        return Ok(unsafe { slice::from_raw_parts(ptr, 0) });
    }

    if ptr.is_null() {
        return Err(format!("{name} ptr is null but len > 0"));
    }

    Ok(unsafe { slice::from_raw_parts(ptr, len) })
}

// ============ KDF ============

fn derive_dek(
    local_material: &[u8],
    cloud_material: &[u8],
    file_id: &[u8],
    file_salt: &[u8],
) -> Result<[u8; DEK_LEN], String> {
    // IKM = local || cloud
    let mut ikm = Vec::with_capacity(local_material.len() + cloud_material.len());
    ikm.extend_from_slice(local_material);
    ikm.extend_from_slice(cloud_material);

    // HKDF salt = file_salt
    let hk = Hkdf::<Sha256>::new(Some(file_salt), &ikm);

    // info = tag || 0x00 || file_id
    let mut info = Vec::with_capacity(32 + 1 + file_id.len());
    info.extend_from_slice(b"VIPER_FILE_KEY_V1");
    info.push(0);
    info.extend_from_slice(file_id);

    let mut out = [0u8; DEK_LEN];
    hk.expand(&info, &mut out)
        .map_err(|_| "HKDF expand failed".to_string())?;

    Ok(out)
}

// ============ Header ============

fn write_header(
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

struct ParsedHeader {
    salt: Vec<u8>,
    base_nonce: Vec<u8>,
}

fn read_header(r: &mut BufReader<File>) -> Result<ParsedHeader, String> {
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
    let mut base_nonce = vec![0u8; nonce_len];
    r.read_exact(&mut salt)
        .map_err(|e| format!("read salt failed: {e}"))?;
    r.read_exact(&mut base_nonce)
        .map_err(|e| format!("read nonce failed: {e}"))?;

    Ok(ParsedHeader { salt, base_nonce })
}

// ============ Public FFI APIs ============

/// encrypt_file (FINAL V1):
/// - Writes header (FORMAT_VERSION=1)
/// - Frames: [1-byte is_last][4-byte BE ct_len][ct bytes]
/// - Uses encrypt_last on final frame
#[unsafe(no_mangle)]
pub extern "C" fn encrypt_file(
    in_path: *const c_char,
    out_path: *const c_char,
    local_ptr: *const u8,
    local_len: c_int,
    cloud_ptr: *const u8,
    cloud_len: c_int,
    file_id_ptr: *const u8,
    file_id_len: c_int,
    aad_ptr: *const u8,
    aad_len: c_int,
    _flags: u32,
) -> c_int {
    match encrypt_file_impl(
        in_path,
        out_path,
        local_ptr,
        local_len,
        cloud_ptr,
        cloud_len,
        file_id_ptr,
        file_id_len,
        aad_ptr,
        aad_len,
    ) {
        Ok(()) => {
            set_last_error("OK");
            0
        }
        Err(msg) => {
            set_last_error(&msg);
            -1
        }
    }
}

fn encrypt_file_impl(
    in_path: *const c_char,
    out_path: *const c_char,
    local_ptr: *const u8,
    local_len: c_int,
    cloud_ptr: *const u8,
    cloud_len: c_int,
    file_id_ptr: *const u8,
    file_id_len: c_int,
    aad_ptr: *const u8,
    aad_len: c_int,
) -> Result<(), String> {
    let in_path = cstr_to_string(in_path, "in_path")?;
    let out_path = cstr_to_string(out_path, "out_path")?;

    let local_material = bytes_from_ptr(local_ptr, local_len, "local_material")?;
    let cloud_material = bytes_from_ptr(cloud_ptr, cloud_len, "cloud_material")?;
    let file_id = bytes_from_ptr(file_id_ptr, file_id_len, "file_id")?;
    let aad = bytes_from_ptr(aad_ptr, aad_len, "aad")?;

    if local_material.is_empty() {
        return Err("local_material is empty".to_string());
    }
    if cloud_material.is_empty() {
        return Err("cloud_material is empty".to_string());
    }
    if file_id.is_empty() {
        return Err("file_id is empty".to_string());
    }

    // Generate per-file salt
    let mut salt = [0u8; SALT_LEN];
    rand::thread_rng().fill_bytes(&mut salt);

    // Derive DEK
    let dek = derive_dek(local_material, cloud_material, file_id, &salt)?;

    // Create AEAD
    let aead =
        XChaCha20Poly1305::new_from_slice(&dek).map_err(|_| "invalid DEK length".to_string())?;

    // Streaming base nonce: StreamBE32 overhead = 5 bytes => base nonce = 24 - 5 = 19
    let mut base_nonce = [0u8; STREAM_BASE_NONCE_LEN];
    rand::thread_rng().fill_bytes(&mut base_nonce);

    // Open files
    let fin = File::open(&in_path).map_err(|e| format!("open in_path failed: {e}"))?;
    let fout = File::create(&out_path).map_err(|e| format!("create out_path failed: {e}"))?;

    let mut reader = BufReader::new(fin);
    let mut writer = BufWriter::new(fout);

    // Write header
    write_header(&mut writer, &salt, &base_nonce)?;

    // Encrypt stream
    let nonce = GenericArray::from_slice(&base_nonce);
    let mut enc = EncryptorBE32::from_aead(aead, nonce);

    // FINAL V1 framing:
    // Each frame: [1-byte is_last] + [4-byte big-endian ct_len] + [ct bytes]
    // IMPORTANT: In some aead versions, encrypt_last consumes `enc` (moves it),
    // so we must call encrypt_last only once at the end.
    let mut buf_cur = vec![0u8; CHUNK_SIZE];
    let mut buf_next = vec![0u8; CHUNK_SIZE];

    // Read first chunk into cur
    let mut cur_n = reader
        .read(&mut buf_cur)
        .map_err(|e| format!("read plaintext failed: {e}"))?;

    // Empty file => only header, no frames
    if cur_n == 0 {
        writer
            .flush()
            .map_err(|e| format!("flush out failed: {e}"))?;
        return Ok(());
    }

    // Read next chunk into next (lookahead)
    let mut next_n = reader
        .read(&mut buf_next)
        .map_err(|e| format!("read plaintext failed: {e}"))?;

    // While there is a "next" chunk, current is NOT last => encrypt_next
    while next_n != 0 {
        let chunk = &buf_cur[..cur_n];

        let ct = enc
            .encrypt_next(Payload { msg: chunk, aad })
            .map_err(|_| "encrypt_next failed".to_string())?;

        let ct_len = ct.len() as u32;
        writer
            .write_all(&[0u8]) // is_last = 0
            .map_err(|e| format!("write is_last failed: {e}"))?;
        writer
            .write_all(&ct_len.to_be_bytes())
            .map_err(|e| format!("write ct len failed: {e}"))?;
        writer
            .write_all(&ct)
            .map_err(|e| format!("write ciphertext failed: {e}"))?;

        // Rotate: next becomes cur; read a new next
        std::mem::swap(&mut buf_cur, &mut buf_next);
        cur_n = next_n;

        next_n = reader
            .read(&mut buf_next)
            .map_err(|e| format!("read plaintext failed: {e}"))?;
    }

    // Now next_n == 0 => cur_n is the LAST chunk => encrypt_last (may consume enc)
    let last_chunk = &buf_cur[..cur_n];
    let ct = enc
        .encrypt_last(Payload {
            msg: last_chunk,
            aad,
        })
        .map_err(|_| "encrypt_last failed".to_string())?;

    let ct_len = ct.len() as u32;
    writer
        .write_all(&[1u8]) // is_last = 1
        .map_err(|e| format!("write is_last failed: {e}"))?;
    writer
        .write_all(&ct_len.to_be_bytes())
        .map_err(|e| format!("write ct len failed: {e}"))?;
    writer
        .write_all(&ct)
        .map_err(|e| format!("write ciphertext failed: {e}"))?;

    writer
        .flush()
        .map_err(|e| format!("flush out failed: {e}"))?;
    Ok(())
}

/// decrypt_file (FINAL V1):
/// - Reads header (FORMAT_VERSION=1)
/// - Frames: [1-byte is_last][4-byte BE ct_len][ct bytes]
/// - Uses decrypt_last on final frame, and rejects missing-last / trailing-data
#[unsafe(no_mangle)]
pub extern "C" fn decrypt_file(
    in_path: *const c_char,
    out_path: *const c_char,
    local_ptr: *const u8,
    local_len: c_int,
    cloud_ptr: *const u8,
    cloud_len: c_int,
    file_id_ptr: *const u8,
    file_id_len: c_int,
    aad_ptr: *const u8,
    aad_len: c_int,
    _flags: u32,
) -> c_int {
    match decrypt_file_impl(
        in_path,
        out_path,
        local_ptr,
        local_len,
        cloud_ptr,
        cloud_len,
        file_id_ptr,
        file_id_len,
        aad_ptr,
        aad_len,
    ) {
        Ok(()) => {
            set_last_error("OK");
            0
        }
        Err(msg) => {
            set_last_error(&msg);
            -1
        }
    }
}

fn decrypt_file_impl(
    in_path: *const c_char,
    out_path: *const c_char,
    local_ptr: *const u8,
    local_len: c_int,
    cloud_ptr: *const u8,
    cloud_len: c_int,
    file_id_ptr: *const u8,
    file_id_len: c_int,
    aad_ptr: *const u8,
    aad_len: c_int,
) -> Result<(), String> {
    let in_path = cstr_to_string(in_path, "in_path")?;
    let out_path = cstr_to_string(out_path, "out_path")?;

    let local_material = bytes_from_ptr(local_ptr, local_len, "local_material")?;
    let cloud_material = bytes_from_ptr(cloud_ptr, cloud_len, "cloud_material")?;
    let file_id = bytes_from_ptr(file_id_ptr, file_id_len, "file_id")?;
    let aad = bytes_from_ptr(aad_ptr, aad_len, "aad")?;

    if local_material.is_empty() {
        return Err("local_material is empty".to_string());
    }
    if cloud_material.is_empty() {
        return Err("cloud_material is empty".to_string());
    }
    if file_id.is_empty() {
        return Err("file_id is empty".to_string());
    }

    // Open input and parse header (FINAL V1 requires FORMAT_VERSION=1)
    let fin = File::open(&in_path).map_err(|e| format!("open in_path failed: {e}"))?;
    let mut reader = BufReader::new(fin);
    let header = read_header(&mut reader)?;

    // Re-derive DEK using header salt
    let dek = derive_dek(local_material, cloud_material, file_id, &header.salt)?;
    let aead =
        XChaCha20Poly1305::new_from_slice(&dek).map_err(|_| "invalid DEK length".to_string())?;

    let nonce = GenericArray::from_slice(header.base_nonce.as_slice());
    let mut dec = DecryptorBE32::from_aead(aead, nonce);

    // Prepare output file
    let fout = File::create(&out_path).map_err(|e| format!("create out_path failed: {e}"))?;
    let mut writer = BufWriter::new(fout);

    // FINAL V1 framing:
    // Each frame: [1-byte is_last] + [4-byte BE ct_len] + [ct bytes]
    // Must end with is_last=1, and then EOF (no trailing bytes).
    loop {
        // Read is_last. EOF here means missing last frame => error.
        let mut flag = [0u8; 1];
        reader.read_exact(&mut flag).map_err(|e| {
            if e.kind() == std::io::ErrorKind::UnexpectedEof {
                "unexpected EOF while reading is_last (missing last frame)".to_string()
            } else {
                format!("read is_last failed: {e}")
            }
        })?;

        let is_last = match flag[0] {
            0 => false,
            1 => true,
            x => return Err(format!("invalid is_last flag: {x}")),
        };

        // Read ct_len
        let mut len_buf = [0u8; 4];
        reader
            .read_exact(&mut len_buf)
            .map_err(|e| format!("read frame length failed: {e}"))?;
        let ct_len = u32::from_be_bytes(len_buf) as usize;

        if ct_len == 0 || ct_len > (CHUNK_SIZE + 32) {
            return Err(format!("invalid ct_len: {ct_len}"));
        }

        // Read ciphertext
        let mut ct = vec![0u8; ct_len];
        reader
            .read_exact(&mut ct)
            .map_err(|e| format!("read ciphertext failed: {e}"))?;

        if !is_last {
            // Non-last frame => decrypt_next (does NOT move dec)
            let pt = dec
                .decrypt_next(Payload { msg: &ct, aad })
                .map_err(|_| "decrypt_next failed (AAD mismatch / corrupted data)".to_string())?;

            writer
                .write_all(&pt)
                .map_err(|e| format!("write plaintext failed: {e}"))?;

            continue;
        }

        // Last frame => decrypt_last (may MOVE dec in this crate version)
        let pt = dec
            .decrypt_last(Payload { msg: &ct, aad })
            .map_err(|_| "decrypt_last failed (AAD mismatch / corrupted data)".to_string())?;

        writer
            .write_all(&pt)
            .map_err(|e| format!("write plaintext failed: {e}"))?;

        // Strictly require EOF after last frame (no trailing garbage)
        let mut extra = [0u8; 1];
        match reader.read(&mut extra) {
            Ok(0) => {} // EOF OK
            Ok(_) => return Err("trailing data after last frame".to_string()),
            Err(e) => return Err(format!("read trailing data failed: {e}")),
        }

        break;
    }

    writer
        .flush()
        .map_err(|e| format!("flush out failed: {e}"))?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::ffi::CString;
    use std::fs;
    use std::os::raw::c_int;
    use tempfile::tempdir;

    fn last_err() -> String {
        unsafe {
            let p = get_last_error_message();
            if p.is_null() {
                return "<null>".to_string();
            }
            std::ffi::CStr::from_ptr(p).to_string_lossy().into_owned()
        }
    }

    #[test]
    fn roundtrip_small_v1_final() {
        let dir = tempdir().unwrap();

        let in_path = dir.path().join("plain.bin");
        let enc_path = dir.path().join("enc.vpf");
        let out_path = dir.path().join("out.bin");

        fs::write(&in_path, b"hello viper").unwrap();

        let local = b"local_key_material_32_bytes________";
        let cloud = b"cloud_key_material_32_bytes________";
        let file_id = b"file1";
        let aad = b"aad";

        let in_c = CString::new(in_path.to_str().unwrap()).unwrap();
        let enc_c = CString::new(enc_path.to_str().unwrap()).unwrap();
        let out_c = CString::new(out_path.to_str().unwrap()).unwrap();

        let rc = encrypt_file(
            in_c.as_ptr(),
            enc_c.as_ptr(),
            local.as_ptr(),
            local.len() as c_int,
            cloud.as_ptr(),
            cloud.len() as c_int,
            file_id.as_ptr(),
            file_id.len() as c_int,
            aad.as_ptr(),
            aad.len() as c_int,
            0,
        );
        assert_eq!(rc, 0, "encrypt failed: {}", last_err());

        let rc = decrypt_file(
            enc_c.as_ptr(),
            out_c.as_ptr(),
            local.as_ptr(),
            local.len() as c_int,
            cloud.as_ptr(),
            cloud.len() as c_int,
            file_id.as_ptr(),
            file_id.len() as c_int,
            aad.as_ptr(),
            aad.len() as c_int,
            0,
        );
        assert_eq!(rc, 0, "decrypt failed: {}", last_err());

        let out = fs::read(&out_path).unwrap();
        assert_eq!(out, b"hello viper");
    }

    #[test]
    fn roundtrip_cross_chunk_v1_final() {
        let dir = tempdir().unwrap();

        let in_path = dir.path().join("plain_big.bin");
        let enc_path = dir.path().join("enc_big.vpf");
        let out_path = dir.path().join("out_big.bin");

        let mut data = vec![0u8; CHUNK_SIZE * 2 + 123];
        rand::thread_rng().fill_bytes(&mut data);
        fs::write(&in_path, &data).unwrap();

        let local = b"local_key_material_32_bytes________";
        let cloud = b"cloud_key_material_32_bytes________";
        let file_id = b"file_big";
        let aad = b"aad";

        let in_c = CString::new(in_path.to_str().unwrap()).unwrap();
        let enc_c = CString::new(enc_path.to_str().unwrap()).unwrap();
        let out_c = CString::new(out_path.to_str().unwrap()).unwrap();

        let rc = encrypt_file(
            in_c.as_ptr(),
            enc_c.as_ptr(),
            local.as_ptr(),
            local.len() as c_int,
            cloud.as_ptr(),
            cloud.len() as c_int,
            file_id.as_ptr(),
            file_id.len() as c_int,
            aad.as_ptr(),
            aad.len() as c_int,
            0,
        );
        assert_eq!(rc, 0, "encrypt failed: {}", last_err());

        let rc = decrypt_file(
            enc_c.as_ptr(),
            out_c.as_ptr(),
            local.as_ptr(),
            local.len() as c_int,
            cloud.as_ptr(),
            cloud.len() as c_int,
            file_id.as_ptr(),
            file_id.len() as c_int,
            aad.as_ptr(),
            aad.len() as c_int,
            0,
        );
        assert_eq!(rc, 0, "decrypt failed: {}", last_err());

        let out = fs::read(&out_path).unwrap();
        assert_eq!(out, data);
    }

    #[test]
    fn decrypt_fails_on_truncated_missing_last_frame_v1_final() {
        let dir = tempdir().unwrap();

        let in_path = dir.path().join("plain.bin");
        let enc_path = dir.path().join("enc.vpf");
        let out_path = dir.path().join("out.bin");

        let mut data = vec![0u8; CHUNK_SIZE + 10];
        rand::thread_rng().fill_bytes(&mut data);
        fs::write(&in_path, &data).unwrap();

        let local = b"local_key_material_32_bytes________";
        let cloud = b"cloud_key_material_32_bytes________";
        let file_id = b"file1";
        let aad = b"aad";

        let in_c = CString::new(in_path.to_str().unwrap()).unwrap();
        let enc_c = CString::new(enc_path.to_str().unwrap()).unwrap();
        let out_c = CString::new(out_path.to_str().unwrap()).unwrap();

        let rc = encrypt_file(
            in_c.as_ptr(),
            enc_c.as_ptr(),
            local.as_ptr(),
            local.len() as c_int,
            cloud.as_ptr(),
            cloud.len() as c_int,
            file_id.as_ptr(),
            file_id.len() as c_int,
            aad.as_ptr(),
            aad.len() as c_int,
            0,
        );
        assert_eq!(rc, 0, "encrypt failed: {}", last_err());

        // truncate some bytes from the end (breaks last frame)
        let mut bytes = fs::read(&enc_path).unwrap();
        bytes.truncate(bytes.len().saturating_sub(10));
        fs::write(&enc_path, bytes).unwrap();

        let rc = decrypt_file(
            enc_c.as_ptr(),
            out_c.as_ptr(),
            local.as_ptr(),
            local.len() as c_int,
            cloud.as_ptr(),
            cloud.len() as c_int,
            file_id.as_ptr(),
            file_id.len() as c_int,
            aad.as_ptr(),
            aad.len() as c_int,
            0,
        );
        assert_eq!(rc, -1, "decrypt should fail on truncated file");
    }
}
