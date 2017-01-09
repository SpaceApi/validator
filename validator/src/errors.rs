error_chain! {
    errors {
        WrongApiVersion(api: String) {
            description("invalid api version")
            display("invalid api version: '{}'", api)
        }
    }

    foreign_links {
        Parse(::serde_json::error::Error);
    }
}
