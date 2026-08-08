use crate::core::common::CHUNK_SIZE;
use crossbeam_channel::{Receiver, Sender, bounded};

pub struct BufferPool {
    tx: Sender<Vec<u8>>,
    rx: Receiver<Vec<u8>>,
}

impl BufferPool {
    pub fn new(count: usize) -> Self {
        let (tx, rx) = bounded(count);

        for _ in 0..count {
            // 检查：文档描述encrypt_in_place在buffer容量不足时会返回错误，实际对Vec会自动扩容，
            // 但我们还是增加密文tag的长度，避免它自动扩容
            tx.send(Vec::with_capacity(CHUNK_SIZE + 16)).unwrap();
        }

        Self { tx, rx }
    }

    pub fn get(&self) -> Vec<u8> {
        self.rx.recv().unwrap()
    }

    pub fn release(&self, mut buf: Vec<u8>) {
        buf.clear();
        self.tx.send(buf).unwrap();
    }
}
