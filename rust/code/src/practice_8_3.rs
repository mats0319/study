/* 第8章练习题 */

// 第3题：允许通过文本接口向公司指定部门添加员工。允许根据部门查看员工、或所有员工的字典序
// example: f3()

use regex::Regex;
use std::collections::{HashMap, HashSet};
use std::io;

pub fn f3() {
    let mut department_has: HashMap<String, HashSet<String>> = HashMap::new(); // department - users 1:n
    let mut user_belong: HashMap<String, String> = HashMap::new(); // user - department 1:1

    let re_add = Regex::new(r"(?i)^add\s+(\w+)\s+to\s+(\w+)\s*$").unwrap();
    let re_list = Regex::new(r"(?i)^list(?:\s+(\w+))?\s*$").unwrap();

    loop {
        println!("\n> Input your command: ");

        let mut input = String::new();

        io::stdin().read_line(&mut input).expect("read error");

        let input = input.trim();

        if let Some(caps) = re_add.captures(input) {
            // add user

            let user_name = caps[1].to_string();
            let department = caps[2].to_string();

            if user_belong.contains_key(&user_name) {
                println!("> User already Exist.");
                continue;
            }

            user_belong.insert(user_name.clone(), department.clone());
            department_has
                .entry(department)
                .or_default()
                .insert(user_name);

            println!("> Add User Success.");
        } else if let Some(caps) = re_list.captures(input) {
            // list

            match caps.get(1) {
                Some(str) => match department_has.get(str.as_str()) {
                    Some(users) => {
                        let mut users_for_sort: Vec<_> = users.iter().collect();
                        users_for_sort.sort();

                        println!("> Department: {}, Users: {users_for_sort:?}", str.as_str());
                    }
                    None => println!("> Unknown Department: {}", str.as_str()),
                },
                None => {
                    // list all

                    let mut department_list: Vec<_> = department_has.keys().collect();
                    department_list.sort();

                    for d in department_list {
                        let mut users_for_sort: Vec<_> =
                            department_has.get(d).unwrap().iter().collect();
                        users_for_sort.sort();

                        println!("> Department: {d}, Users: {users_for_sort:?}");
                    }
                }
            }
        } else {
            println!("> Unknown Command: {input}");
        }
    }
}
