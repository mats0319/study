#[cfg(test)]
mod tests {
    use crate::core::decrypt_file_impl;
    use crate::core::encrypt_v2::encrypt_file_impl_2;
    use sha2::{Digest, Sha256};
    use std::fs;
    use std::time::Instant;

    const WORK_DIR: &str = "./testdata/";
    const ORIGIN_FILE: &str = "./testdata/test_100m.bin";
    const ENC_FILE: &str = "./testdata/test_100m.enc";
    const DEC_FILE: &str = "./testdata/test_100m.dec";
    const LOCAL_MATERIAL: &[u8] = b"test local material";
    const CLOUD_MATERIAL: &[u8] = b"test cloud material";
    const FILE_ID: &[u8] = b"test file id";
    const AAD: &[u8] = b"test aad";

    #[test]
    fn test_encdec_function() {
        if !fs::exists(ORIGIN_FILE).unwrap() {
            let _ = fs::create_dir_all(WORK_DIR);
            fs::write(ORIGIN_FILE, b"This is a test file.").unwrap();
        }

        let start = Instant::now();
        encrypt_file_impl_2(
            ORIGIN_FILE,
            ENC_FILE,
            LOCAL_MATERIAL,
            CLOUD_MATERIAL,
            FILE_ID,
            AAD,
        )
        .unwrap();
        println!("Encrypt completed in {:?}", start.elapsed());

        let start = Instant::now();
        decrypt_file_impl(
            ENC_FILE,
            DEC_FILE,
            LOCAL_MATERIAL,
            CLOUD_MATERIAL,
            FILE_ID,
            AAD,
        )
        .unwrap();
        println!("Decrypt completed in {:?}", start.elapsed());

        assert_eq!(file_hash(ORIGIN_FILE), file_hash(DEC_FILE));
    }

    fn file_hash(path: &str) -> String {
        let data = fs::read(path).unwrap();
        let hash = Sha256::digest(&data);
        hash.iter().map(|b| format!("{:02x}", b)).collect()
    }
}
