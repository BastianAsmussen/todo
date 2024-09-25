use std::{fs::File, path::Path};

use clap::Parser;
use cli::Cli;
use fs::TaskHandler;

mod cli;
mod fs;

fn main() {
    let _cli = Cli::parse();

    let path = Path::new("test.csv");
    if !path.exists() {
        println!("{} doesn't exist, creating it...", path.display());

        if let Err(why) = File::create(path) {
            eprintln!("Failed to create file: {why}");
            std::process::exit(1);
        }
    }

    let handler = TaskHandler::new(path);
    println!("{:#?}", handler.read().expect("Failed to read!"));
}
