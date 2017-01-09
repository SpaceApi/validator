extern crate valico;
extern crate serde_json;

use std::error::Error;
pub use valico::common::error::ValicoError;

use serde_json::{Value};
use valico::json_schema;

pub fn get_schema(api_version: &str) -> Result<Value, Box<Error>> {
    let schema = match api_version {
        "0.13" => include_str!("../schema/13.json"),
        "0.12" => include_str!("../schema/12.json"),
        "0.11" => include_str!("../schema/11.json"),
        "0.9" => include_str!("../schema/9.json"),
        "0.8" => include_str!("../schema/8.json"),
        _ => return Err(format!("Invalid api version: {}", api_version).into())
    }.parse()?;
    Ok(schema)
}

pub fn validate_spaceapi_json(json: &str) -> Result<json_schema::ValidationState, Box<Error>> {
    let json_value: Value = json.parse()?;

    let json_obj = json_value.as_object().ok_or("Invalid JSON")?;

    let version = json_obj.get("api").ok_or("api missing")?
        .as_str().ok_or("api key not a string")?;

    let json_schema: Value = get_schema(version)?;
    let mut scope = json_schema::Scope::new();
    let schema = scope.compile_and_return(json_schema, false).unwrap();
    let validation_result = schema.validate(&json_value);

    Ok(validation_result)
}


#[cfg(test)]
mod test {
    use ::validate_spaceapi_json;
    #[test]
    fn validate_coredump() {
        let obj = include_str!("../coredump.json");
        let validated = validate_spaceapi_json(obj);
        println!("{:?}", validated);
    }
}
