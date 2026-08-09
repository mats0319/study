/* 第8章练习题 */

// 第1题：给出一个数组，统计它的中位数和众数
// example: f1(&vec![3, 5, 2, 1, 4, 5, 3])

pub fn f1(vec: &[i32]) {
    if vec.is_empty() {
        println!("empty vector");
        return;
    }

    let mut vec_clone = vec.to_vec(); // 避免修改原始数组
    vec_clone.sort_unstable();

    let mid_1 = (vec_clone.len() - 1) / 2;
    let mid_2 = vec_clone.len() / 2;
    let mid = (vec_clone[mid_1] as f64 + vec_clone[mid_2] as f64) / 2.0;

    let mut current = vec_clone[0];
    let mut count = 0;
    let mut max_count = 0;
    let mut nums = Vec::new(); // 众数

    for &item in &vec_clone {
        if item == current {
            count += 1;
        } else {
            if count > max_count {
                max_count = count;
                nums.clear();
                nums.push(current);
            } else if count == max_count {
                nums.push(current);
            }

            current = item;
            count = 1;
        }
    }

    if count > max_count {
        max_count = count;
        nums.clear();
        nums.push(current);
    } else if count == max_count {
        nums.push(current);
    }

    println!("vector mid: {mid}, max count: {nums:?} ({max_count} times)");
}
