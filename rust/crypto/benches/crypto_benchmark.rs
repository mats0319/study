use criterion::{Criterion, criterion_group, criterion_main};
use std::time::Duration;

use crypto::core::Encrypt;
use crypto::core::encrypt_v1::EncryptV1;
use crypto::core::encrypt_v2::EncryptV2;

fn benchmark_encrypt(c: &mut Criterion) {
    let mut group = c.benchmark_group("encrypt");

    group.sample_size(30);
    group.measurement_time(Duration::from_secs(30));

    let local_material = "test local key".as_bytes();
    let cloud_material = "test cloud key".as_bytes();
    let file_id = "test file id".as_bytes();
    let aad = "test aad".as_bytes();

    // // enc 10m
    // group.bench_function("encrypt v1 10m", |b| {
    //     b.iter(|| {
    //         EncryptV1::encrypt(
    //             "./benches/testdata/test_10m.bin",
    //             "./benches/testdata/test_10m.enc",
    //             local_material,
    //             cloud_material,
    //             file_id,
    //             aad,
    //         )
    //         .unwrap();
    //     })
    // });
    // group.bench_function("encrypt v2 10m", |b| {
    //     b.iter(|| {
    //         EncryptV2::encrypt(
    //             "./benches/testdata/test_10m.bin",
    //             "./benches/testdata/test_10m.enc",
    //             local_material,
    //             cloud_material,
    //             file_id,
    //             aad,
    //         )
    //         .unwrap();
    //     })
    // });
    //
    // // enc 100m
    // group.bench_function("encrypt v1 100m", |b| {
    //     b.iter(|| {
    //         EncryptV1::encrypt(
    //             "./benches/testdata/test_100m.bin",
    //             "./benches/testdata/test_100m.enc",
    //             local_material,
    //             cloud_material,
    //             file_id,
    //             aad,
    //         )
    //         .unwrap();
    //     })
    // });
    // group.bench_function("encrypt v2 100m", |b| {
    //     b.iter(|| {
    //         EncryptV2::encrypt(
    //             "./benches/testdata/test_100m.bin",
    //             "./benches/testdata/test_100m.enc",
    //             local_material,
    //             cloud_material,
    //             file_id,
    //             aad,
    //         )
    //         .unwrap();
    //     })
    // });

    // enc 1g
    group.bench_function("encrypt v1 1g", |b| {
        b.iter(|| {
            EncryptV1::encrypt(
                "./benches/testdata/test_1g.bin",
                "./benches/testdata/test_1g.enc",
                local_material,
                cloud_material,
                file_id,
                aad,
            )
            .unwrap();
        })
    });
    group.bench_function("encrypt v2 1g", |b| {
        b.iter(|| {
            EncryptV2::encrypt(
                "./benches/testdata/test_1g.bin",
                "./benches/testdata/test_1g.enc",
                local_material,
                cloud_material,
                file_id,
                aad,
            )
            .unwrap();
        })
    });

    group.finish();
}

criterion_group!(benches, benchmark_encrypt);

criterion_main!(benches);
