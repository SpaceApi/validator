#[macro_use] extern crate nickel;
extern crate rustc_serialize;
extern crate url;
extern crate spaceapi_validator;

use std::collections::BTreeMap;
use nickel::status::StatusCode;
use nickel::{Nickel, JsonBody, HttpRouter, MediaType};
use rustc_serialize::json::{encode};
use spaceapi_validator::{validate_spaceapi_json, ValicoError};


#[derive(RustcDecodable, RustcEncodable)]
struct ValidationRequest {
    schema:  String,
}

#[derive(RustcDecodable, RustcEncodable)]
struct Result {
    version: String,
    status: String,
    message: Option<String>,
    errors: Option<String>,
    missing: Option<String>,
}

fn main() {
    let mut server = Nickel::new();

    // try it with curl
    // curl 'http://localhost:6767/' -H 'Content-Type: application/json;charset=UTF-8'  --data-binary $'{"version":"0.13", "schema":"{...}"}'
    server.post("/", middleware! { |request, mut response|
        set_headers(&mut response);

        let vr = try_with!(response, {
            request.json_as::<ValidationRequest>().map_err(|e| (StatusCode::BadRequest, e))
        });

        let result = validate_spaceapi_json(&*vr.schema);
        let result = match result {
          Ok(state) => {
            if state.is_valid() {
              Result {
                version: "TODO".into(),
                status: "OK".into(),
                message: None,
                errors: None,
                missing: Some(encode_missing(&state.missing)),
              }
            } else {
              Result {
                version: "TODO".into(),
                status: "ERROR".into(),
                message: Some(format!("{:?}", state)),
                errors: Some(encode_errors(&state.errors)),
                missing: Some(encode_missing(&state.missing)),
              }
            }
          },
          Err(boxed_error) => {
            let err_msg: &str = (*boxed_error).description();
            Result {
              version: "undefined".into(), // TODO return something more useful?
              status: "ERROR".into(),
              message: Some(err_msg.into()),
              errors: None,
              missing: None,
            }
          },
        };

        response.set(MediaType::Json);
        encode(&result).unwrap()
    });


    // Supporting CORS
    server.get("/", middleware! { |_, mut response|
        set_headers(&mut response);

        "OK"
    });

    server.options("/", middleware! { |_, mut response|
        set_headers(&mut response);

        "OK"
    });

    // TODO hard coded
    server.listen("127.0.0.1:6767").unwrap();
}

/// set CORS header and disable Cache-Control
fn set_headers(response: &mut nickel::Response) {
    response.headers_mut().set_raw("Access-Control-Allow-Origin", vec![b"*".to_vec()]);
    response.headers_mut().set_raw("Access-Control-Allow-Methods", vec![b"GET, POST, OPTIONS".to_vec()]);
    response.headers_mut().set_raw("Access-Control-Allow-Headers", vec![b"Origin, Content-Type, X-Auth-Token".to_vec()]);

    response.headers_mut().set_raw("Cache-Control", vec![b"no-cache".to_vec()]);
}

fn encode_missing(missing: &Vec<url::Url>) -> String {
  let arr = missing.iter()
                       .fold("".to_string(), |base, &ref element| {
                        format!("{}{}...",
                            if base.len() > 0 { "," } else { "" },
                            &*element.path()
                        )
                       });
  format!("[{}]", arr)
}

fn encode_errors(errors: &Vec<Box<ValicoError>>) -> String {
    let v: Vec<BTreeMap<&str, &str>> = errors.iter()
                        .map(|&ref element| {
        let mut map: BTreeMap<&str, &str> = BTreeMap::new();
        map.insert("code" , element.get_code() );
        map.insert("title", element.get_title());
        map.insert("path" , element.get_path() );
        if let Some(d) = element.get_detail() {
            map.insert("detail", d);
        }
        map
    }).collect();
    encode(&v).unwrap_or("[]".into())
}
