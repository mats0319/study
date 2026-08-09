use std::error::Error;
use std::{env, fs, io, process};

pub fn grep() {
    let args: Vec<_> = env::args().collect();

    let config = Config::new(env::args()).unwrap_or_else(|err| {
        println!("Error: {}", err);
        process::exit(1);
    });

    let v = config.search();

    println!("{:#?}", v);
}

pub struct Config {
    text: String,
    content: String,
}

impl Config {
    fn new(mut args: impl Iterator<Item = String>) -> Result<Config, Box<dyn Error>> {
        args.next();

        let text = match args.next() {
            Some(arg) => arg,
            None => return Err(Box::from("No Query Str")),
        };

        let file_path = match args.next() {
            Some(arg) => arg,
            None => return Err(Box::from("No File Path")),
        };

        let content = fs::read_to_string(file_path)?;

        Ok(Config { text, content })
    }

    fn search(&self) -> Vec<&str> {
        self.content
            .lines()
            .filter(|line| line.contains(self.text.as_str()))
            .collect()
    }
}
