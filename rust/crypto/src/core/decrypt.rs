use crate::core::buffer_pool::BufferPool;
use crate::core::common::{
    CHUNK_SIZE, Frame, STREAM_BASE_NONCE_LEN, get_thread_count, make_nonce, read_header,
};
use crate::core::kdf::derive_dek;
use aead::AeadInPlace;
use chacha20poly1305::{KeyInit, XChaCha20Poly1305};
use crossbeam_channel::{Receiver, Sender, bounded, unbounded};
use std::collections::HashMap;
use std::fs::File;
use std::io::{BufReader, BufWriter, Read, Write};
use std::sync::Arc;
use std::thread;

struct Decryptor {
    dec: XChaCha20Poly1305,
    base_nonce: [u8; STREAM_BASE_NONCE_LEN],
    aad: Arc<[u8]>,
}

impl Decryptor {
    fn decrypt_frame(&self, frame: Frame) -> Result<Frame, String> {
        let mut buf = frame.data;

        let nonce = make_nonce(&self.base_nonce, frame.is_last_frame, frame.index);

        self.dec
            .decrypt_in_place(&nonce.into(), &self.aad, &mut buf)
            .map_err(|e| format!("decrypt_in_place failed: {e}"))?;

        Ok(Frame {
            index: frame.index,
            is_last_frame: frame.is_last_frame,
            data: buf,
        })
    }
}

pub fn decrypt_file_impl(
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

    // channels
    let (cipher_tx, cipher_rx): (Sender<Frame>, Receiver<Frame>) = bounded(thread_count * 2);
    let (plain_tx, plain_rx): (Sender<Frame>, Receiver<Frame>) = unbounded();

    // file handler
    let fin = File::open(&in_path).map_err(|e| format!("open in_path failed: {e}"))?;
    let fout = File::create(&out_path).map_err(|e| format!("create out_path failed: {e}"))?;

    // parse 'salt' and 'base nonce'
    let mut reader = BufReader::new(fin);

    let (salt, base_nonce) = read_header(&mut reader)?;

    // derive key
    let dek = derive_dek(local_material, cloud_material, file_id, &salt)?;

    // 1. read thread
    let read_thread = thread::spawn(move || -> Result<(), String> {
        let mut index: u32 = 0;

        loop {
            let mut flag = [0u8; 1];
            reader.read_exact(&mut flag).map_err(|e| {
                if e.kind() == std::io::ErrorKind::UnexpectedEof {
                    "unexpected EOF while reading is_last (missing last frame)".to_string()
                } else {
                    format!("read is_last failed: {e}")
                }
            })?;

            let is_last_frame = match flag[0] {
                0 => false,
                1 => true,
                x => return Err(format!("invalid is_last flag: {x}")),
            };

            let mut ct_len = [0u8; 4];
            reader
                .read_exact(&mut ct_len)
                .map_err(|e| format!("read ct len failed: {e}"))?;

            let ct_len = u32::from_be_bytes(ct_len) as usize;

            if ct_len == 0 || ct_len > (CHUNK_SIZE + 32) {
                return Err(format!("invalid ct_len: {ct_len}"));
            }

            let mut buffer = read_pool.get();
            buffer.clear();
            buffer.resize(ct_len, 0);

            reader
                .read_exact(&mut buffer)
                .map_err(|e| format!("read ciphertext failed: {e}"))?;

            cipher_tx
                .send(Frame {
                    index,
                    is_last_frame,
                    data: buffer,
                })
                .map_err(|e| format!("send cipher frame failed: {e}"))?;

            if is_last_frame {
                break;
            }

            index += 1;
        }

        Ok(())
    });

    // 2. decrypt thread
    let dec = Arc::new(Decryptor {
        dec: XChaCha20Poly1305::new_from_slice(&dek)
            .map_err(|_| "invalid DEK length".to_string())?,
        base_nonce,
        aad: aad.clone(),
    });
    let mut dec_thread_list = Vec::new();
    for _ in 0..thread_count {
        let rx = cipher_rx.clone();
        let tx = plain_tx.clone();
        let dec = dec.clone();

        let thread = thread::spawn(move || -> Result<(), String> {
            for frame in rx {
                let decrypted_frame = dec.decrypt_frame(frame)?;
                tx.send(decrypted_frame)
                    .map_err(|e| format!("send plain frame failed: {e}"))?;
            }

            Ok(())
        });

        dec_thread_list.push(thread);
    }

    // 3. write thread
    let write_thread = thread::spawn(move || -> Result<(), String> {
        let mut writer = BufWriter::new(fout);

        let mut pending_frames: HashMap<u32, Frame> = HashMap::new();
        let mut last_index: u32 = 0;
        let mut write_index: u32 = 0;

        for frame in plain_rx {
            if frame.is_last_frame {
                last_index = frame.index;
            }

            if frame.index == write_index {
                writer
                    .write_all(&frame.data)
                    .map_err(|e| format!("write plaintext failed: {e}"))?;
                write_pool.release(frame.data);
                write_index += 1;

                while let Some(next_frame) = pending_frames.remove(&write_index) {
                    writer
                        .write_all(&next_frame.data)
                        .map_err(|e| format!("write plaintext failed: {e}"))?;
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
            .map_err(|e| format!("flush out failed: {e}"))?;

        Ok(())
    });

    let _ = read_thread
        .join()
        .map_err(|e| format!("read thread failed: {e:?}"))?;
    for thread in dec_thread_list {
        let _ = thread
            .join()
            .map_err(|e| format!("decrypt thread failed: {e:?}"))?;
    }
    let _ = write_thread
        .join()
        .map_err(|e| format!("write thread failed: {e:?}"))?;

    Ok(())
}
