# Validation server for SpaceAPI endpoints

Written in [Rust](https://rust-lang.org).

[![Travis CI][travis-ci-badge]][travis-ci]

# How to build?

    git submodule init
    git submodule update

    cd service
    cargo build

# Test validator locally

Start the frontend:

    cd frontend
    python -m http.server 8080

Start the service:

    cd service
    cargo run

Then open the page with your favorite browser:

    http://localhost:8080/validator.html


<!-- Badges -->
[travis-ci]: https://travis-ci.org/spacedirectory/validator
[travis-ci-badge]: https://img.shields.io/travis/spacedirectory/validator.svg
