use crate::core::Encrypt;
use crate::core::buffer_pool::BufferPool;
use crate::core::common::{CHUNK_SIZE, Frame, SALT_LEN, STREAM_BASE_NONCE_LEN, get_thread_count, make_nonce, write_header};
use crate::core::kdf::derive_dek;
use aead::AeadInPlace;
use aead::rand_core::RngCore;
use chacha20poly1305::{KeyInit, XChaCha20Poly1305};
use crossbeam_channel::{Receiver, Sender, bounded, unbounded};
use std::collections::HashMap;
use std::fs::File;
use std::io::{BufReader, BufWriter, Read, Write};
use std::sync::Arc;
use std::thread;

pub struct EncryptV2;

impl Encrypt for EncryptV2 {
    fn encrypt(
        in_path: &str,
        out_path: &str,
        local_material: &[u8],
        cloud_material: &[u8],
        file_id: &[u8],
        aad: &[u8],
    ) -> Result<(), String> {
        encrypt_file_impl_2(
            in_path,
            out_path,
            local_material,
            cloud_material,
            file_id,
            aad,
        )
    }
}

struct Encryptor {
    enc: XChaCha20Poly1305,
    base_nonce: [u8; STREAM_BASE_NONCE_LEN],
    aad: Arc<[u8]>,
}

impl Encryptor {
    fn encrypt_frame(&self, frame: Frame) -> Result<Frame, String> {
        let mut buf = frame.data;

        let nonce = make_nonce(&self.base_nonce, frame.is_last_frame, frame.index);

        self.enc
            .encrypt_in_place(&nonce.into(), &self.aad, &mut buf)
            .map_err(|e| format!("encrypt_in_place failed: {e}"))?;

        Ok(Frame {
            index: frame.index,
            is_last_frame: frame.is_last_frame,
            data: buf,
        })
    }
}

pub fn encrypt_file_impl_2(
    in_path: &str,
    out_path: &str,
    local_material: &[u8],
    cloud_material: &[u8],
    file_id: &[u8],
    aad: &[u8],
) -> Result<(), String> {
    if local_material.is_empty() || cloud_material.is_empty() || file_id.is_empty() {
        return Err(format!(
            "local_material length: {}, cloud_material length: {}, file_id length: {}",
            local_material.len(),
            cloud_material.len(),
            file_id.len(),
        ));
    }

    // params
    let thread_count = get_thread_count();
    let pool = Arc::new(BufferPool::new(thread_count * 4));
    let read_pool = pool.clone();
    let write_pool = pool.clone();
    let aad = Arc::<[u8]>::from(aad);

    // generate 'salt' for derive key and 'base nonce' for encrypt
    let mut random = rand::thread_rng();

    let mut salt = [0u8; SALT_LEN];
    random.fill_bytes(&mut salt);

    let mut base_nonce = [0u8; STREAM_BASE_NONCE_LEN];
    random.fill_bytes(&mut base_nonce);

    // derive key
    let dek = derive_dek(local_material, cloud_material, file_id, &salt)?;

    // channels and copy data for multi-threads
    let (plain_tx, plain_rx): (Sender<Frame>, Receiver<Frame>) = bounded(thread_count * 2);
    let (cipher_tx, cipher_rx): (Sender<Frame>, Receiver<Frame>) = unbounded();

    // get file handler
    let fin = File::open(&in_path).map_err(|e| format!("open in_path failed: {e}"))?;
    let fout = File::create(&out_path).map_err(|e| format!("create out_path failed: {e}"))?;

    // 1. read thread
    let read_thread = thread::spawn(move || -> Result<(), String> {
        let mut reader = BufReader::new(fin);
        let mut index: u32 = 0;

        loop {
            let mut buffer = read_pool.get();
            buffer.clear();
            buffer.resize(CHUNK_SIZE,0);

            let n = reader
                .read(&mut buffer)
                .map_err(|e| format!("read plaintext failed: {e}"))?;

            buffer.truncate(n);

            plain_tx
                .send(Frame {
                    index,
                    is_last_frame: n == 0,
                    data: buffer,
                })
                .map_err(|e| format!("send plain frame failed: {e}"))?;

            if n == 0 {
                break;
            }

            index += 1;
        }

        Ok(())
    });

    // 2. encrypt thread(s)
    let enc = Arc::new(Encryptor {
        enc: XChaCha20Poly1305::new_from_slice(&dek)
            .map_err(|_| "invalid DEK length".to_string())?,
        base_nonce,
        aad: aad.clone(),
    });
    let mut enc_thread_list = Vec::new();
    for _ in 0..thread_count {
        let rx = plain_rx.clone();
        let tx = cipher_tx.clone();
        let enc = enc.clone();

        let thread = thread::spawn(move || -> Result<(), String> {
            for frame in rx {
                let encrypted_frame = enc.encrypt_frame(frame)?;
                tx.send(encrypted_frame)
                    .map_err(|e| format!("send cipher frame failed: {e}"))?;
            }

            Ok(())
        });

        enc_thread_list.push(thread);
    }

    // 3. write thread
    let write_thread = thread::spawn(move || -> Result<(), String> {
        let mut writer = BufWriter::new(fout);

        write_header(&mut writer, &salt, &base_nonce)?;

        let mut pending_frames: HashMap<u32, Frame> = HashMap::new();
        let mut last_index: u32 = 0;
        let mut write_index: u32 = 0;

        for frame in cipher_rx {
            if frame.is_last_frame {
                last_index = frame.index;
            }

            if frame.index == write_index {
                let _ = write_frame(&mut writer, &frame)?;
                write_pool.release(frame.data);
                write_index += 1;

                while let Some(next_frame) = pending_frames.remove(&write_index) {
                    let _ = write_frame(&mut writer, &next_frame)?;
                    write_pool.release(next_frame.data);
                    write_index += 1;
                }
            } else {
                pending_frames.insert(frame.index, frame);
            }

            // 检查：因为在前面+1，所以这里要写大于，或者'rc=li+1'；
            // 因为我们会发送额外的空白帧，所以last_index总是>0的
            if last_index > 0 && write_index > last_index {
                break;
            }
        }

        writer
            .flush()
            .map_err(|e| format!("flush out_path failed: {e}"))?;

        Ok(())
    });

    // 等待所有线程执行完成
    let _ = read_thread
        .join()
        .map_err(|e| format!("read_thread failed: {e:?}"))?;
    for thread in enc_thread_list {
        let _ = thread
            .join()
            .map_err(|e| format!("enc_thread failed: {e:?}"))?;
    }
    let _ = write_thread
        .join()
        .map_err(|e| format!("write_thread failed: {e:?}"))?;

    Ok(())
}

fn write_frame(writer: &mut BufWriter<File>, frame: &Frame) -> Result<(), String> {
    writer
        .write_all(&[if frame.is_last_frame { 1 } else { 0 }])
        .map_err(|e| format!("write_frame failed: {e}"))?;
    writer
        .write_all(&(frame.data.len() as u32).to_be_bytes())
        .map_err(|e| format!("write fixed failed: {e}"))?;
    writer
        .write_all(&frame.data)
        .map_err(|e| format!("write ciphertext failed: {e}"))?;

    Ok(())
}
