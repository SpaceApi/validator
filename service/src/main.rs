#[macro_use] extern crate nickel;
extern crate rustc_serialize;
extern crate spaceapi_validator;

use std::collections::BTreeMap;
use nickel::status::StatusCode;
use nickel::{Nickel, JsonBody, HttpRouter, MediaType};
use rustc_serialize::json::{Json, ToJson, encode};

#[derive(RustcDecodable, RustcEncodable)]
struct ValidationRequest {
    version: String,
    schema:  String,
}

impl ToJson for ValidationRequest {
    fn to_json(&self) -> Json {
        let mut map = BTreeMap::new();
        map.insert("version".to_string(), self.version.to_json());
        map.insert("schema" .to_string(), self.schema .to_json());
        Json::Object(map)
    }
}

#[derive(RustcDecodable, RustcEncodable)]
struct Result {
    version: String,
    status: String,
    message: Option<String>,
}

fn main() {
    let mut server = Nickel::new();

    // try it with curl
    // curl 'http://localhost:6767/' -H 'Content-Type: application/json;charset=UTF-8'  --data-binary $'{"version":"0.13", "schema":"{...}"}'
    server.post("/", middleware! { |request, mut response|
        response.headers_mut().set_raw("Access-Control-Allow-Origin", vec![b"*".to_vec()]);
        response.headers_mut().set_raw("Access-Control-Allow-Methods", vec![b"POST,OPTIONS".to_vec()]);

        let vr = try_with!(response, {
            request.json_as::<ValidationRequest>().map_err(|e| (StatusCode::BadRequest, e))
        });


        // TODO: insert validation here

        //format!("Hello {} {}", person.version, person.schema)
        encode(&Result {
            version: "0.13".into(),
            status: "OK".into(),
            message: None,
        }).unwrap()
    });

    server.get("/", middleware! { |request, mut response|
        response.headers_mut().set_raw("Access-Control-Allow-Origin", vec![b"*".to_vec()]);
        response.headers_mut().set_raw("Access-Control-Allow-Methods", vec![b"POST,OPTIONS".to_vec()]);

        "OK"
    });

    server.listen("127.0.0.1:6767").unwrap();
}
