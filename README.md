# Validation server for SpaceAPI endpoints

Written in [Rust](https://rust-lang.org).

[![Travis CI][travis-ci-badge]][travis-ci]

# How to build?

    git submodule init
    git submodule update

    cd service
    cargo build

# Test validator locally

Start the service:

    cd service
    cargo run

You can now query the API using a HTTP client:

    http://localhost:6767/

# License

Licensed under either of

 * Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or
   http://www.apache.org/licenses/LICENSE-2.0)
 * MIT license ([LICENSE-MIT](LICENSE-MIT) or
   http://opensource.org/licenses/MIT)

at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall
be dual licensed as above, without any additional terms or conditions.

<!-- Badges -->
[travis-ci]: https://travis-ci.org/spacedirectory/validator
[travis-ci-badge]: https://img.shields.io/travis/spacedirectory/validator/master.svg
