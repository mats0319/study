use std::cmp::Ordering;
use std::io;

pub fn guess_number() {
    println!("Guess the Number! (1~100)");

    let secret_number = rand::random_range(1..=100);
    let mut small = 1;
    let mut big = 100;

    loop {
        println!("Input your guess: (valid range [{small},{big}])");

        let mut guess = String::new();

        io::stdin()
            .read_line(&mut guess)
            .expect("Failed to read line");

        let guess: u32 = match guess.trim().parse() {
            Ok(num) => num,
            Err(_) => continue,
        };

        if !(small <= guess && guess <= big) {
            println!("Invalid Number: {guess}, Range is [{small},{big}]");
            continue;
        }

        print!("You guessed: {guess}, ");

        match guess.cmp(&secret_number) {
            Ordering::Less => {
                println!("Too small!");
                small = guess;
            }
            Ordering::Greater => {
                println!("Too big!");
                big = guess;
            }
            Ordering::Equal => {
                println!("You win!");
                break;
            }
        }
    }
}
