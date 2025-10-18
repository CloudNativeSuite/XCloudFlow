use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Play {
    pub name: String,
    pub tasks: Vec<Task>,
}

#[derive(Debug, Deserialize)]
pub struct Task {
    pub name: String,
    pub shell: Option<String>,
    pub script: Option<String>,
    pub template: Option<Template>,
}

#[derive(Debug, Deserialize)]
pub struct Template {
    pub src: String,
    pub dest: String,
}
