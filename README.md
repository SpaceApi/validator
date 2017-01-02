# Validation server for SpaceAPI endpoints

Written in [Rust](https://rust-lang.org).

Insert Build Badge here [ ...... ]

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

# Goals

  - Validator
    - Form
    - Route
  - Directory
    - In Memory
    - Persistent Directory
