use hkdf::Hkdf;
use sha2::Sha256;

use crate::core::common::DEK_LEN;

pub fn derive_dek(
    local_material: &[u8],
    cloud_material: &[u8],
    file_id: &[u8],
    file_salt: &[u8],
) -> Result<[u8; DEK_LEN], String> {
    // HKDF salt = file_salt
    let hk = Hkdf::<Sha256>::new(Some(file_salt), &[local_material, cloud_material].concat());

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
