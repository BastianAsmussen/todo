use clap::{Args, Parser, Subcommand};

/// A simple todo program.
#[derive(Parser)]
#[command(version, about, long_about = None)]
pub struct Cli {
    /// The action to perform.
    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Subcommand)]
pub enum Commands {
    /// Show all the tasks to perform.
    List {
        /// Show completed entries.
        #[arg(short, long)]
        all: bool,
    },
    /// Create a new task.
    Add(AddArgs),
    /// Delete a task.
    Remove {},
    /// Mark a task as completed.
    Complete {},
}

#[derive(Args)]
pub struct AddArgs {
    pub task: String,
}
