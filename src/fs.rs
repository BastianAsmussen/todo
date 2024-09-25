use std::{error::Error, path::Path};

use chrono::{DateTime, Utc};
use csv::{Reader, Writer};
use serde::{Deserialize, Serialize};

pub struct TaskHandler<'a> {
    path: &'a Path,
}

impl<'a> TaskHandler<'a> {
    pub const fn new(path: &'a Path) -> Self {
        Self { path }
    }

    pub fn read(&self) -> Result<Vec<Task>, Box<dyn Error>> {
        let mut tasks = Vec::new();

        let mut reader = Reader::from_path(self.path)?;
        for result in reader.deserialize() {
            let record: Task = result?;

            tasks.push(record);
        }

        Ok(tasks)
    }

    pub fn write(&self, mut task: Task) -> Result<(), Box<dyn Error>> {
        let mut tasks = self.read()?;
        if tasks.iter().any(|t| t.id == task.id) {
            task.id = u32::try_from(tasks.len())?;
        }

        tasks.push(task);

        let mut writer = Writer::from_path(self.path)?;
        for task in tasks {
            writer.serialize(task)?;
        }

        writer.flush()?;

        Ok(())
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Task {
    id: u32,
    name: String,
    due: DateTime<Utc>,
    completed: bool,
}

impl Task {
    pub const fn new(id: u32, name: String, due: DateTime<Utc>, completed: bool) -> Self {
        Self {
            id,
            name,
            due,
            completed,
        }
    }
}
